package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"
)

type Point2D struct {
	X int
	Y int
}

type Tile struct {
	IsOnMainLoop bool
	Neighbors    []Point2D
	VisitStatus  VisitStatus
}

type Sketch [][]Tile

func (s Sketch) InBounds(point Point2D) bool {
	return point.Y >= 0 && point.Y < len(s) && point.X >= 0 && point.X < len(s[0])
}

type VisitStatus int

const (
	VisitStatusNotVisited VisitStatus = iota
	VisitStatusPendingVisit
	VisitStatusVisited
)

func runMain() error {
	sketch, start, err := parseSketch("input.txt")
	if err != nil {
		return fmt.Errorf("parse sketch: %w", err)
	}

	totalTiles := len(sketch) * len(sketch[0])
	mainLoopTiles := findMainLoop(sketch, start)
	outsideTiles := countOutsideTiles(sketch)

	fmt.Println(totalTiles - mainLoopTiles - outsideTiles)

	return nil
}

var symbolToDelta = map[rune][2]Point2D{
	'|': {{X: 0, Y: -1}, {X: 0, Y: 1}},
	'-': {{X: -1, Y: 0}, {X: 1, Y: 0}},
	'L': {{X: 0, Y: -1}, {X: 1, Y: 0}},
	'J': {{X: 0, Y: -1}, {X: -1, Y: 0}},
	'7': {{X: -1, Y: 0}, {X: 0, Y: 1}},
	'F': {{X: 1, Y: 0}, {X: 0, Y: 1}},
}

func parseSketch(filename string) (Sketch, Point2D, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, Point2D{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var sketch Sketch
	var start Point2D

	scanner := bufio.NewScanner(f)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		row := make([]Tile, len(line))
		for j, tile := range line {
			switch tile {
			case 'S':
				start = Point2D{X: j, Y: i}
				continue
			case '.':
				continue
			}

			delta := symbolToDelta[tile]
			row[j].Neighbors = []Point2D{
				{X: delta[0].X + j, Y: delta[0].Y + i},
				{X: delta[1].X + j, Y: delta[1].Y + i},
			}
		}

		sketch = append(sketch, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, Point2D{}, fmt.Errorf("scan: %w", err)
	}

	fixupStartShape(sketch, start)

	return sketch, start, nil
}

func fixupStartShape(sketch Sketch, start Point2D) {
	for _, delta := range [][2]Point2D{
		symbolToDelta['|'],
		symbolToDelta['-'],
		symbolToDelta['L'],
		symbolToDelta['J'],
		symbolToDelta['7'],
		symbolToDelta['F'],
	} {
		from := Point2D{
			X: start.X + delta[0].X,
			Y: start.Y + delta[0].Y,
		}

		to := Point2D{
			X: start.X + delta[1].X,
			Y: start.Y + delta[1].Y,
		}

		if !sketch.InBounds(from) || !sketch.InBounds(to) {
			continue
		}

		fromTile := sketch[from.Y][from.X]
		toTile := sketch[to.Y][to.X]

		if !slices.Contains(fromTile.Neighbors, start) || !slices.Contains(toTile.Neighbors, start) {
			continue
		}

		sketch[start.Y][start.X].Neighbors = []Point2D{from, to}
		return
	}

	panic(fmt.Errorf("fixup start shape %+v", start))
}

func findMainLoop(sketch Sketch, start Point2D) int {
	prevTile := start
	currentTile := sketch[start.Y][start.X].Neighbors[0]
	mainLoopTiles := 1

	for currentTile != start {
		sketch[currentTile.Y][currentTile.X].IsOnMainLoop = true
		neighbors := sketch[currentTile.Y][currentTile.X].Neighbors

		if neighbors[0] != prevTile {
			prevTile = currentTile
			currentTile = neighbors[0]
		} else {
			prevTile = currentTile
			currentTile = neighbors[1]
		}

		mainLoopTiles++
	}

	sketch[currentTile.Y][currentTile.X].IsOnMainLoop = true

	return mainLoopTiles
}

func countOutsideTiles(sketch Sketch) int {
	res := 0

	for i := range sketch {
		for _, j := range []int{0, len(sketch[0]) - 1} {
			if sketch[i][j].VisitStatus == VisitStatusNotVisited {
				res += countOutsideTilesFrom(sketch, i, j)
			}
		}
	}

	for j := range sketch[0] {
		for _, i := range []int{0, len(sketch) - 1} {
			if sketch[i][j].VisitStatus == VisitStatusNotVisited {
				res += countOutsideTilesFrom(sketch, i, j)
			}
		}
	}

	return res
}

type FloodFillPoint struct {
	Point Point2D
	Side  Point2D
}

func countOutsideTilesFrom(sketch Sketch, i int, j int) int {
	res := 0

	for points := []FloodFillPoint{{Point: Point2D{X: j, Y: i}}}; len(points) > 0; points = points[1:] {
		point := points[0]

		sketch[point.Point.Y][point.Point.X].VisitStatus = VisitStatusVisited

		tile := sketch[point.Point.Y][point.Point.X]
		if !tile.IsOnMainLoop {
			res++
		}

		for _, delta := range [8]Point2D{
			{X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1},
			{X: -1, Y: 0}, {X: 1, Y: 0},
			{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 1, Y: 1},
		} {
			nextPoint := Point2D{
				X: point.Point.X + delta.X,
				Y: point.Point.Y + delta.Y,
			}

			if !sketch.InBounds(nextPoint) {
				continue
			}

			nextTile := sketch[nextPoint.Y][nextPoint.X]
			if nextTile.VisitStatus != VisitStatusNotVisited {
				continue
			}

			if tile.IsOnMainLoop && areVectorsIntersect(Point2D{
				X: point.Point.X - point.Side.X,
				Y: point.Point.Y - point.Side.Y,
			}, nextPoint, tile.Neighbors[0], tile.Neighbors[1]) {
				continue
			}

			sketch[nextPoint.Y][nextPoint.X].VisitStatus = VisitStatusPendingVisit
			points = append(points, FloodFillPoint{
				Point: nextPoint,
				Side: Point2D{
					X: nextPoint.X - point.Point.X,
					Y: nextPoint.Y - point.Point.Y,
				},
			})
		}
	}

	return res
}

// orientation determines the orientation of three given points (a, b, c).
//
// Assume the segment from a to b as a vector.
// We can find the order by checking the direction of the vector from b to c relative to the vector from a to b.
//
//   - If the vector from b to c deviates to the left of the vector from a to b,
//     then the points are in counter-clockwise order.
//   - If the vector from b to c deviates to the right of the vector from a to b,
//     then the points are in clockwise order.
//   - If the vector from b to c is in the same direction as that from a to b or in the exact opposite direction,
//     then the points are collinear.
//
// The function returns 1 if the points are in counter-clockwise order,
// -1 if they are in clockwise order, and 0 if they are collinear.
func orientation(a, b, c Point2D) int {
	switch val := (c.Y-a.Y)*(b.X-a.X) - (b.Y-a.Y)*(c.X-a.X); {
	case val > 0:
		return 1
	case val < 0:
		return -1
	default:
		return 0
	}
}

// areVectorsIntersect determines if the line segments formed by points a and b, and points c and d intersect.
func areVectorsIntersect(a, b, c, d Point2D) bool {
	o1 := orientation(a, b, c)
	o2 := orientation(a, b, d)
	o3 := orientation(c, d, a)
	o4 := orientation(c, d, b)

	if o1 != o2 && o3 != o4 {
		return true
	}

	if o1 == 0 && isPointOnSegment(a, c, b) {
		return true
	}

	if o2 == 0 && isPointOnSegment(a, d, b) {
		return true
	}

	if o3 == 0 && isPointOnSegment(c, a, d) {
		return true
	}

	if o4 == 0 && isPointOnSegment(c, b, d) {
		return true
	}

	return false
}

// isPointOnSegment determines if point b lies on the line segment formed by points a and c.
func isPointOnSegment(a, b, c Point2D) bool {
	return b.X <= max(a.X, c.X) && b.X >= min(a.X, c.X) &&
		b.Y <= max(a.Y, c.Y) && b.Y >= min(a.Y, c.Y)
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
