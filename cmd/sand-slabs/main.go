package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type Point3D struct {
	X int
	Y int
	Z int
}

func (p Point3D) Add(other Point3D) Point3D {
	return Point3D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

type Stack struct {
	Bricks map[Point3D]int
	Slabs  []Slab
}

type Slab struct {
	ID        int
	Start     Point3D
	Direction Point3D
	Length    int
}

func isSubset[T comparable](set map[T]struct{}, subset map[T]struct{}) bool {
	for elem := range subset {
		if _, ok := set[elem]; !ok {
			return false
		}
	}

	return true
}

func union[T comparable](a, b map[T]struct{}) map[T]struct{} {
	res := make(map[T]struct{}, len(a))

	for elem := range a {
		res[elem] = struct{}{}
	}

	for elem := range b {
		res[elem] = struct{}{}
	}

	return res
}

func runMain() error {
	stack, err := parseStack("input.txt")
	if err != nil {
		return fmt.Errorf("parse stack: %w", err)
	}

	stack.fallSlabs()
	supportGraph := stack.getSupportGraph()

	slices.Reverse(stack.Slabs)

	sum := 0

	for _, slab := range stack.Slabs {
		fallen := getFallen(supportGraph, slab.ID, map[int]struct{}{slab.ID: {}})
		sum += len(fallen) - 1
	}

	fmt.Println(sum)

	return nil
}

func parseStack(filename string) (Stack, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Stack{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	stack := Stack{Bricks: map[Point3D]int{}}

	scanner := bufio.NewScanner(f)

	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()

		points := strings.Split(line, "~")
		from := parsePoint3D(points[0])
		to := parsePoint3D(points[1])

		slab := Slab{
			ID:    i,
			Start: from,
		}

		switch {
		case from.X != to.X:
			slab.Direction = Point3D{X: 1}
			slab.Length = to.X + 1 - from.X
		case from.Y != to.Y:
			slab.Direction = Point3D{Y: 1}
			slab.Length = to.Y + 1 - from.Y
		case from.Z != to.Z:
			slab.Direction = Point3D{Z: 1}
			slab.Length = to.Z + 1 - from.Z
		default:
			slab.Length = 1
		}

		for j, direction := 0, (Point3D{}); j < slab.Length; j, direction = j+1, direction.Add(slab.Direction) {
			stack.Bricks[from.Add(direction)] = i
		}

		stack.Slabs = append(stack.Slabs, slab)
	}

	if err := scanner.Err(); err != nil {
		return Stack{}, fmt.Errorf("scan: %w", err)
	}

	return stack, nil
}

func parsePoint3D(line string) Point3D {
	var res Point3D
	_, _ = fmt.Sscanf(line, "%d,%d,%d", &res.X, &res.Y, &res.Z)
	return res
}

func (s *Stack) fallSlabs() {
	slices.SortFunc(s.Slabs, func(a, b Slab) int {
		return cmp.Compare(a.Start.Z, b.Start.Z)
	})

	for i, slab := range s.Slabs {
		original := slab.Start

		for slab.Start.Z > 1 {
			if !s.canSlabFall(slab) {
				break
			}

			slab.Start = slab.Start.Add(Point3D{Z: -1})
		}

		if slab.Start != original {
			brick := slab.Start
			originalBrick := original

			for j := 0; j < slab.Length; j++ {
				delete(s.Bricks, originalBrick)
				s.Bricks[brick] = slab.ID

				brick = brick.Add(slab.Direction)
				originalBrick = originalBrick.Add(slab.Direction)
			}

			s.Slabs[i] = slab
		}
	}
}

func (s *Stack) canSlabFall(slab Slab) bool {
	for i, brick := 0, slab.Start; i < slab.Length; i, brick = i+1, brick.Add(slab.Direction) {
		if other, ok := s.Bricks[brick.Add(Point3D{Z: -1})]; ok && other != slab.ID {
			return false
		}
	}

	return true
}

type SupportGraph struct {
	Support     map[int]map[int]struct{}
	SupportedBy map[int]map[int]struct{}
}

func (s *Stack) getSupportGraph() SupportGraph {
	res := SupportGraph{
		Support:     make(map[int]map[int]struct{}, len(s.Slabs)),
		SupportedBy: make(map[int]map[int]struct{}, len(s.Slabs)),
	}

	for _, slab := range s.Slabs {
		for i, brick := 0, slab.Start; i < slab.Length; i, brick = i+1, brick.Add(slab.Direction) {
			if other, ok := s.Bricks[brick.Add(Point3D{Z: 1})]; ok && other != slab.ID {
				support, ok := res.Support[slab.ID]
				if !ok {
					support = make(map[int]struct{})
					res.Support[slab.ID] = support
				}

				support[other] = struct{}{}
			}

			if other, ok := s.Bricks[brick.Add(Point3D{Z: -1})]; ok && other != slab.ID {
				supportedBy, ok := res.SupportedBy[slab.ID]
				if !ok {
					supportedBy = make(map[int]struct{})
					res.SupportedBy[slab.ID] = supportedBy
				}

				supportedBy[other] = struct{}{}
			}
		}
	}

	return res
}

func getFallen(supportGraph SupportGraph, slab int, fallen map[int]struct{}) map[int]struct{} {
	for target := range supportGraph.Support[slab] {
		if isSubset(fallen, supportGraph.SupportedBy[target]) {
			fallen = union(fallen, getFallen(supportGraph, target, union(fallen, map[int]struct{}{target: {}})))
		}
	}

	return fallen
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
