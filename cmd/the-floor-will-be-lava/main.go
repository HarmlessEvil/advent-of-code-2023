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

func (p Point2D) Add(other Point2D) Point2D {
	return Point2D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

type Tile int

const (
	TileEmptySpace Tile = iota
	TileMirrorForward
	TileMirrorBack
	TileSplitterVertical
	TileSplitterHorizontal
)

type Contraption [][]Tile

func (c Contraption) InBounds(point Point2D) bool {
	return point.Y >= 0 && point.Y < len(c) && point.X >= 0 && point.X < len(c[0])
}

func runMain() error {
	contraption, err := parseContraption()
	if err != nil {
		return fmt.Errorf("parse contraption: %w", err)
	}

	var maxEnergised int

	for i := range contraption {
		visitStatus := simulateLightBeam(contraption, BFSPoint{
			Point: Point2D{
				X: 0,
				Y: i,
			},
			Direction: DirectionRight,
		})

		maxEnergised = max(maxEnergised, countEnergized(contraption, visitStatus))
	}

	for i := range contraption {
		visitStatus := simulateLightBeam(contraption, BFSPoint{
			Point: Point2D{
				X: len(contraption[0]) - 1,
				Y: i,
			},
			Direction: DirectionLeft,
		})

		maxEnergised = max(maxEnergised, countEnergized(contraption, visitStatus))
	}

	for j := range contraption[0] {
		visitStatus := simulateLightBeam(contraption, BFSPoint{
			Point: Point2D{
				X: j,
				Y: 0,
			},
			Direction: DirectionDown,
		})

		maxEnergised = max(maxEnergised, countEnergized(contraption, visitStatus))
	}

	for j := range contraption[0] {
		visitStatus := simulateLightBeam(contraption, BFSPoint{
			Point: Point2D{
				X: j,
				Y: len(contraption) - 1,
			},
			Direction: DirectionUp,
		})

		maxEnergised = max(maxEnergised, countEnergized(contraption, visitStatus))
	}

	fmt.Println(maxEnergised)

	return nil
}

func parseContraption() (Contraption, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var contraption Contraption

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		row := make([]Tile, len(line))

		for i, symbol := range line {
			var tile Tile

			switch symbol {
			case '/':
				tile = TileMirrorForward
			case '\\':
				tile = TileMirrorBack
			case '|':
				tile = TileSplitterVertical
			case '-':
				tile = TileSplitterHorizontal
			}

			row[i] = tile
		}

		contraption = append(contraption, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return contraption, nil
}

type VisitStatus int

const (
	VisitStatusNotVisited VisitStatus = iota
	VisitStatusPendingVisit
	VisitStatusVisited
)

type Direction int

const (
	DirectionRight Direction = iota
	DirectionDown
	DirectionLeft
	DirectionUp
)

func (d Direction) ToVector() Point2D {
	switch d {
	case DirectionRight:
		return Point2D{X: 1, Y: 0}
	case DirectionDown:
		return Point2D{X: 0, Y: 1}
	case DirectionLeft:
		return Point2D{X: -1, Y: 0}
	case DirectionUp:
		return Point2D{X: 0, Y: -1}
	}

	return Point2D{}
}

func (d Direction) MirrorForward() Direction {
	return map[Direction]Direction{
		DirectionRight: DirectionUp,
		DirectionDown:  DirectionLeft,
		DirectionLeft:  DirectionDown,
		DirectionUp:    DirectionRight,
	}[d]
}

func (d Direction) MirrorBack() Direction {
	return map[Direction]Direction{
		DirectionRight: DirectionDown,
		DirectionDown:  DirectionRight,
		DirectionLeft:  DirectionUp,
		DirectionUp:    DirectionLeft,
	}[d]
}

type BeamVisitStatus struct {
	Dirs [4]VisitStatus
}

func (s BeamVisitStatus) IsEnergized() bool {
	return s.Dirs[DirectionRight] == VisitStatusVisited ||
		s.Dirs[DirectionDown] == VisitStatusVisited ||
		s.Dirs[DirectionLeft] == VisitStatusVisited ||
		s.Dirs[DirectionUp] == VisitStatusVisited
}

type BFSPoint struct {
	Point     Point2D
	Direction Direction
}

func (p BFSPoint) Move(direction Direction) BFSPoint {
	return BFSPoint{
		Point:     p.Point.Add(direction.ToVector()),
		Direction: direction,
	}
}

func simulateLightBeam(contraption Contraption, start BFSPoint) map[Point2D]BeamVisitStatus {
	visitStatus := make(map[Point2D]BeamVisitStatus, len(contraption)*len(contraption[0]))

	for points := []BFSPoint{start}; len(points) > 0; points = points[1:] {
		point := points[0]

		pointVisitStatus := visitStatus[point.Point]
		pointVisitStatus.Dirs[point.Direction] = VisitStatusVisited
		visitStatus[point.Point] = pointVisitStatus

		nextPoints := make([]BFSPoint, 0, 2)
		switch contraption[point.Point.Y][point.Point.X] {
		case TileEmptySpace:
			nextPoints = append(nextPoints, point.Move(point.Direction))
		case TileSplitterVertical:
			nextPoints = append(nextPoints, point.Move(DirectionUp), point.Move(DirectionDown))
		case TileSplitterHorizontal:
			nextPoints = append(nextPoints, point.Move(DirectionLeft), point.Move(DirectionRight))
		case TileMirrorForward:
			nextPoints = append(nextPoints, point.Move(point.Direction.MirrorForward()))
		case TileMirrorBack:
			nextPoints = append(nextPoints, point.Move(point.Direction.MirrorBack()))
		}

		for _, nextPoint := range nextPoints {
			nextPointVisitStatus := visitStatus[nextPoint.Point]
			if contraption.InBounds(nextPoint.Point) && nextPointVisitStatus.Dirs[nextPoint.Direction] == VisitStatusNotVisited {
				nextPointVisitStatus.Dirs[nextPoint.Direction] = VisitStatusPendingVisit
				visitStatus[nextPoint.Point] = nextPointVisitStatus

				points = append(points, nextPoint)
			}
		}
	}

	return visitStatus
}

func countEnergized(contraption Contraption, visitStatus map[Point2D]BeamVisitStatus) int {
	count := 0

	for i, row := range contraption {
		for j := range row {
			if visitStatus[Point2D{
				X: j,
				Y: i,
			}].IsEnergized() {
				count++
			}
		}
	}

	return count
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
