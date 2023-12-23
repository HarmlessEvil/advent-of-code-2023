package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

type Point2D struct {
	X int
	Y int
}

type HikingTrailMap struct {
	Size Point2D
	Map  map[Point2D][]Point2D
}

func runMain() error {
	trailMap, err := parseHikingTrailMap("input.txt")
	if err != nil {
		return fmt.Errorf("parse hiking trail map: %w", err)
	}

	startPoint := Point2D{X: 1}
	nodeToColor := colorGraph(trailMap.Map, startPoint)
	graph := condenseGraph(nodeToColor)

	if len(nodeToColor[startPoint]) == 0 {
		return fmt.Errorf("start point's id not found")
	}

	var start int
	for scc := range nodeToColor[startPoint] {
		start = scc
		break
	}

	endPoint := Point2D{
		X: trailMap.Size.X - 2,
		Y: trailMap.Size.Y - 1,
	}

	if len(nodeToColor[endPoint]) == 0 {
		return fmt.Errorf("end point's id not found")
	}

	var end int
	for scc := range nodeToColor[endPoint] {
		end = scc
		break
	}

	path := findLongestPath(graph, Crossroad{
		ID:    start,
		Point: startPoint,
	}, end)

	fmt.Println(path.Length - 1)

	return nil
}

func parseHikingTrailMap(filename string) (HikingTrailMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return HikingTrailMap{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	trailMap := HikingTrailMap{Map: map[Point2D][]Point2D{}}

	scanner := bufio.NewScanner(f)

	var previousLine string

	scanner.Scan()
	line := scanner.Text()

	i := 0
	for ; scanner.Scan(); i++ {
		nextLine := scanner.Text()

		parseLine(trailMap, i, previousLine, line, nextLine)

		previousLine, line = line, nextLine
	}

	parseLine(trailMap, i, previousLine, line, "")

	trailMap.Size = Point2D{X: len(line), Y: i + 1}

	if err := scanner.Err(); err != nil {
		return HikingTrailMap{}, fmt.Errorf("scan: %w", err)
	}

	return trailMap, nil
}

func parseLine(trailMap HikingTrailMap, i int, previousLine string, line string, nextLine string) {
	for j, symbol := range line {
		point := Point2D{X: j, Y: i}

		if symbol == '#' {
			continue
		}

		if len(previousLine) > 0 && previousLine[j] != '#' {
			trailMap.Map[point] = append(trailMap.Map[point], Point2D{
				X: point.X,
				Y: point.Y - 1,
			})
		}

		if j < len(line) && line[j+1] != '#' {
			trailMap.Map[point] = append(trailMap.Map[point], Point2D{
				X: point.X + 1,
				Y: point.Y,
			})
		}

		if len(nextLine) > 0 && nextLine[j] != '#' {
			trailMap.Map[point] = append(trailMap.Map[point], Point2D{
				X: point.X,
				Y: point.Y + 1,
			})
		}

		if j > 0 && line[j-1] != '#' {
			trailMap.Map[point] = append(trailMap.Map[point], Point2D{
				X: point.X - 1,
				Y: point.Y,
			})
		}
	}
}

type Node struct {
	ID     int
	Edges  map[int]Point2D
	Weight int
}

func colorGraph(graph map[Point2D][]Point2D, start Point2D) map[Point2D]map[int]struct{} {
	nodeToColor := make(map[Point2D]map[int]struct{}, len(graph))

	nextFreeColor := 0

	var dfs func(point Point2D, index int)
	dfs = func(point Point2D, index int) {
		neighbors := graph[point]

		for ; len(neighbors) <= 2; neighbors = graph[point] {
			nodeToColor[point] = map[int]struct{}{index: {}}

			var next *Point2D
			for _, neighbor := range neighbors {
				if _, ok := nodeToColor[neighbor]; !ok {
					n := neighbor
					next = &n
					break
				}
			}

			if next == nil {
				return
			}

			point = *next
		}

		if _, ok := nodeToColor[point]; !ok {
			nodeToColor[point] = map[int]struct{}{}
		}
		nodeToColor[point][index] = struct{}{}

		for _, next := range neighbors {
			if _, ok := nodeToColor[next]; !ok {
				nextFreeColor++
				dfs(next, nextFreeColor)
			}

			for id := range nodeToColor[next] {
				nodeToColor[point][id] = struct{}{}
			}
		}
	}

	dfs(start, nextFreeColor)

	return nodeToColor
}

func setToSlice[E comparable](set map[E]struct{}) []E {
	res := make([]E, 0, len(set))
	for elem := range set {
		res = append(res, elem)
	}

	return res
}

func condenseGraph(coloring map[Point2D]map[int]struct{}) map[int]Node {
	res := map[int]Node{}

	for point, c := range coloring {
		colors := setToSlice(c)
		if len(colors) > 1 {
			for i := range colors {
				for j := i + 1; j < len(colors); j++ {
					if _, ok := res[colors[i]]; !ok {
						res[colors[i]] = Node{
							ID:    colors[i],
							Edges: map[int]Point2D{},
						}
					}

					node := res[colors[i]]
					node.Edges[colors[j]] = point
					res[colors[i]] = node

					if _, ok := res[colors[j]]; !ok {
						res[colors[j]] = Node{
							ID:    colors[j],
							Edges: map[int]Point2D{},
						}
					}

					node = res[colors[j]]
					node.Edges[colors[i]] = point
					res[colors[j]] = node
				}
			}
		} else {
			if _, ok := res[colors[0]]; !ok {
				res[colors[0]] = Node{
					ID:    colors[0],
					Edges: map[int]Point2D{},
				}
			}

			node := res[colors[0]]
			node.Weight++
			res[colors[0]] = node
		}
	}

	return res
}

type Path struct {
	Nodes  []Node
	Length int
}

func findLongestPath(graph map[int]Node, start Crossroad, end int) Path {
	return dfs(graph, end, map[Point2D]*Crossroad{}, start)
}

type Crossroad struct {
	ID    int
	Point Point2D
}

func dfs(graph map[int]Node, target int, prev map[Point2D]*Crossroad, current Crossroad) Path {
	if current.ID == target {
		res := Path{
			Nodes:  make([]Node, 0, len(prev)),
			Length: 0,
		}

		for cur := &current; cur != nil; cur = prev[cur.Point] {
			node := graph[cur.ID]
			node.ID = cur.ID

			res.Nodes = append(res.Nodes, node)
			res.Length += node.Weight + 1 // crossroad at the end
		}

		res.Length-- // extra crossroad at the end
		return res
	}

	var res Path
	for next, nextPoint := range graph[current.ID].Edges {
		if _, ok := prev[nextPoint]; ok {
			continue
		}

		prev[nextPoint] = &Crossroad{
			ID:    current.ID,
			Point: current.Point,
		}

		path := dfs(graph, target, prev, Crossroad{
			ID:    next,
			Point: nextPoint,
		})

		if path.Length > res.Length {
			res = path
		}

		delete(prev, nextPoint)
	}

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
