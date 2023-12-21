package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

type Tile int

const (
	TileGardenPlot Tile = iota
	TileRock
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

type GardenMap struct {
	Tiles [][]Tile
	Start Point2D
}

func (m GardenMap) At(p Point2D) Tile {
	return m.Tiles[mod(p.Y, len(m.Tiles))][mod(p.X, len(m.Tiles[0]))]
}

func mod(a, b int) int {
	return (a%b + b) % b
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

func f(x int, a0, a1, a2 int) int {
	b0 := a0
	b1 := a1 - a0
	b2 := a2 - a1

	return b0 + b1*x + (x*(x-1)/2)*(b2-b1)
}

const Steps = 26_501_365

func runMain() error {
	gardenMap, err := parseGardenMap("input.txt")
	if err != nil {
		return fmt.Errorf("parse garden map: %w", err)
	}

	coefficients := bfs(gardenMap, []int{65, 196, 327})
	fmt.Println(f(Steps/len(gardenMap.Tiles), coefficients[0], coefficients[1], coefficients[2]))

	return nil
}

func parseGardenMap(filename string) (GardenMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return GardenMap{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var gardenMap GardenMap

	scanner := bufio.NewScanner(f)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		row := make([]Tile, len(line))

		for j, symbol := range line {
			var tile Tile
			switch symbol {
			case '#':
				tile = TileRock
			case 'S':
				gardenMap.Start = Point2D{
					X: j,
					Y: i,
				}
			}

			row[j] = tile
		}

		gardenMap.Tiles = append(gardenMap.Tiles, row)
	}

	if err := scanner.Err(); err != nil {
		return GardenMap{}, fmt.Errorf("scan: %w", err)
	}

	return gardenMap, nil
}

func bfs(gardenMap GardenMap, steps []int) []int {
	front := []Point2D{gardenMap.Start}
	var visited map[Point2D]struct{}

	res := make([]int, len(steps))

	for i, j := 0, 0; j < len(steps); i++ {
		var nextFront []Point2D
		visited = map[Point2D]struct{}{}

		for _, point := range front {
			for _, direction := range []Direction{DirectionRight, DirectionDown, DirectionLeft, DirectionUp} {
				nextPoint := point.Add(direction.ToVector())

				if _, ok := visited[nextPoint]; ok {
					continue
				}

				if gardenMap.At(nextPoint) == TileRock {
					continue
				}

				visited[nextPoint] = struct{}{}
				nextFront = append(nextFront, nextPoint)
			}
		}

		front = nextFront

		if i+1 == steps[j] {
			res[j] = len(visited)
			j++
		}
	}

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
