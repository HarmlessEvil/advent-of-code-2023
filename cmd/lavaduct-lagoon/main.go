package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Point2D struct {
	X int64
	Y int64
}

func (p Point2D) Add(other Point2D) Point2D {
	return Point2D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Point2D) Mul(scalar int64) Point2D {
	return Point2D{
		X: p.X * scalar,
		Y: p.Y * scalar,
	}
}

type Direction int

const (
	DirectionRight Direction = iota
	DirectionDown
	DirectionLeft
	DirectionUp
)

func (d Direction) ToVector() Point2D {
	return map[Direction]Point2D{
		DirectionRight: {X: 1, Y: 0},
		DirectionDown:  {X: 0, Y: 1},
		DirectionLeft:  {X: -1, Y: 0},
		DirectionUp:    {X: 0, Y: -1},
	}[d]
}

type Polygon struct {
	Perimeter int64
	Points    []Point2D
}

type Instruction struct {
	Direction Direction
	Distance  int64
}

func runMain() error {
	polygon, err := parsePolygon("input.txt")
	if err != nil {
		return fmt.Errorf("parse polygon: %w", err)
	}

	fmt.Println(polygon.Area() + polygon.Perimeter/2 + 1)

	return nil
}

func parsePolygon(filename string) (Polygon, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Polygon{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var polygon Polygon
	var current Point2D

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		instruction := parseInstruction(scanner.Bytes())

		polygon.Points = append(polygon.Points, current)
		polygon.Perimeter += instruction.Distance

		current = current.Add(instruction.Direction.ToVector().Mul(instruction.Distance))
	}

	if err := scanner.Err(); err != nil {
		return Polygon{}, fmt.Errorf("scan: %w", err)
	}

	return polygon, nil
}

func parseInstruction(line []byte) Instruction {
	var res Instruction

	n := len(line)
	res.Distance, _ = strconv.ParseInt(string(line[n-7:n-2]), 16, 64)
	res.Direction = Direction(line[n-2] - '0')

	return res
}

func (p Polygon) Area() int64 {
	area := int64(0)

	for i := 1; i < len(p.Points); i++ {
		area += p.Points[i-1].X*p.Points[i].Y - p.Points[i-1].Y*p.Points[i].X
	}

	return area / 2
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
