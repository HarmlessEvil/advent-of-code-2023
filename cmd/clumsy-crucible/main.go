package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
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

type CityMap [][]int

func (m CityMap) InBounds(p Point2D) bool {
	return p.Y >= 0 && p.Y < len(m) && p.X >= 0 && p.X < len(m[0])
}

func runMain() error {
	cityMap, err := parseCityMap()
	if err != nil {
		return fmt.Errorf("parse city map: %w", err)
	}

	state := dijkstra(cityMap, DijkstraPoint{})

	minState := DijkstraState{Distance: math.MaxInt}
	for _, direction := range []Direction{DirectionRight, DirectionDown} {
		for count := 4; count <= 10; count++ {
			if item, ok := state[DijkstraPoint{
				Count:     count,
				Direction: direction,
				Coordinate: Point2D{
					X: len(cityMap[0]) - 1,
					Y: len(cityMap) - 1,
				},
			}]; ok {
				if item.Distance < minState.Distance {
					minState = item
				}
			}
		}
	}

	solution := make([][]string, len(cityMap))
	for i := range solution {
		solution[i] = make([]string, len(cityMap[i]))

		for j := range cityMap[i] {
			solution[i][j] = strconv.Itoa(cityMap[i][j])
		}
	}

	for point := minState; point.Previous != nil; point = *point.Previous {
		var direction string
		switch point.Point.Direction {
		case DirectionRight:
			direction = ">"
		case DirectionDown:
			direction = "v"
		case DirectionLeft:
			direction = "<"
		case DirectionUp:
			direction = "^"
		}

		solution[point.Point.Coordinate.Y][point.Point.Coordinate.X] = direction
	}

	for _, row := range solution {
		fmt.Println(row)
	}

	fmt.Println(minState.Distance)

	return nil
}

func parseCityMap() (CityMap, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var cityMap CityMap

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		row := make([]int, len(line))
		for i, heatLoss := range line {
			row[i] = int(heatLoss - '0')
		}

		cityMap = append(cityMap, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return cityMap, nil
}

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

func (d Direction) Reverse() Direction {
	return map[Direction]Direction{
		DirectionRight: DirectionLeft,
		DirectionDown:  DirectionUp,
		DirectionLeft:  DirectionRight,
		DirectionUp:    DirectionDown,
	}[d]
}

type DijkstraPoint struct {
	Count      int
	Direction  Direction
	Coordinate Point2D
}

func (p DijkstraPoint) Move(direction Direction) DijkstraPoint {
	count := 1
	if p.Direction == direction {
		count = p.Count + 1
	}

	return DijkstraPoint{
		Count:      count,
		Direction:  direction,
		Coordinate: p.Coordinate.Add(direction.ToVector()),
	}
}

type DijkstraState struct {
	Distance int
	Point    DijkstraPoint
	Previous *DijkstraState
}

type DijkstraStateMap map[DijkstraPoint]DijkstraState

type DijkstraQueue []DijkstraState

func (d DijkstraQueue) Len() int {
	return len(d)
}

func (d DijkstraQueue) Less(i, j int) bool {
	return d[i].Distance < d[j].Distance
}

func (d DijkstraQueue) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d *DijkstraQueue) Push(x any) {
	*d = append(*d, x.(DijkstraState))
}

func (d *DijkstraQueue) Pop() any {
	n := len(*d)

	item := (*d)[n-1]
	*d = (*d)[0 : n-1]

	return item
}

func dijkstra(cityMap CityMap, start DijkstraPoint) DijkstraStateMap {
	queue := make(DijkstraQueue, 1, len(cityMap)*len(cityMap[0])*4*3)
	queue[0] = DijkstraState{Point: start}
	heap.Init(&queue)

	state := make(DijkstraStateMap, cap(queue))

	for queue.Len() > 0 {
		point := heap.Pop(&queue).(DijkstraState)

		var nextPoints []DijkstraPoint
		if point.Point.Count > 0 && point.Point.Count < 4 {
			nextPoints = []DijkstraPoint{point.Point.Move(point.Point.Direction)}
		} else {
			nextPoints = []DijkstraPoint{
				point.Point.Move(DirectionRight),
				point.Point.Move(DirectionDown),
				point.Point.Move(DirectionLeft),
				point.Point.Move(DirectionUp),
			}
		}

		for _, nextPoint := range nextPoints {
			if !cityMap.InBounds(nextPoint.Coordinate) ||
				nextPoint.Count > 10 ||
				nextPoint.Direction == point.Point.Direction.Reverse() {
				continue
			}

			nextDistance := point.Distance + cityMap[nextPoint.Coordinate.Y][nextPoint.Coordinate.X]
			if nextState, ok := state[nextPoint]; !ok || nextDistance < nextState.Distance {
				nextState := DijkstraState{
					Distance: nextDistance,
					Point:    nextPoint,
					Previous: &point,
				}

				state[nextPoint] = nextState
				heap.Push(&queue, nextState)
			}
		}
	}

	return state
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
