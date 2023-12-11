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

type Location struct {
	GalaxyID int
}

type Image struct {
	ExpandedRows    map[int]struct{}
	ExpandedColumns map[int]struct{}
	Galaxies        []Point2D
	Locations       [][]Location
}

func (i Image) InBounds(point Point2D) bool {
	return point.Y >= 0 && point.Y < len(i.Locations) && point.X >= 0 && point.X < len(i.Locations[0])
}

func runMain() error {
	image, err := parseImage("input.txt")
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	sum := 0

	for i := 0; i < len(image.Galaxies); i++ {
		for _, dist := range bfs(image, i+1) {
			sum += dist
		}
	}

	fmt.Println(sum)

	return nil
}

func parseImage(filename string) (Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Image{}, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var image Image
	var galaxyCount int

	scanner := bufio.NewScanner(f)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		row := make([]Location, len(line))
		for j, symbol := range line {
			if symbol == '#' {
				galaxyCount++
				row[j].GalaxyID = galaxyCount

				image.Galaxies = append(image.Galaxies, Point2D{X: j, Y: i})
			}
		}

		image.Locations = append(image.Locations, row)
	}

	if err := scanner.Err(); err != nil {
		return Image{}, fmt.Errorf("scan: %w", err)
	}

	image.ExpandedRows = expandSpaceRows(image.Locations)
	image.ExpandedColumns = expandSpaceCols(image.Locations)

	return image, nil
}

func expandSpaceRows(locations [][]Location) map[int]struct{} {
	expendedRows := map[int]struct{}{}

searchEmptyRows:
	for i, row := range locations {
		for _, loc := range row {
			if loc.GalaxyID != 0 {
				continue searchEmptyRows
			}
		}

		expendedRows[i] = struct{}{}
	}

	return expendedRows
}

func expandSpaceCols(locations [][]Location) map[int]struct{} {
	expendedCols := map[int]struct{}{}

searchEmptyCols:
	for j := range locations[0] {
		for i := range locations {
			if locations[i][j].GalaxyID != 0 {
				continue searchEmptyCols
			}
		}

		expendedCols[j] = struct{}{}
	}

	return expendedCols
}

type VisitStatus int

const (
	VisitStatusNotVisited VisitStatus = iota
	VisitStatusPendingVisit
	VisitStatusVisited
)

type BFSPoint struct {
	Point    Point2D
	Distance int
}

const expansionRate = 1000000

func bfs(image Image, sourceGalaxyID int) []int {
	res := make([]int, len(image.Galaxies)-sourceGalaxyID)

	target := make(map[int]struct{}, len(res))
	for i := 0; i < len(res); i++ {
		target[sourceGalaxyID+i+1] = struct{}{}
	}

	sourceGalaxy := image.Galaxies[sourceGalaxyID-1]
	visitStatus := map[Point2D]VisitStatus{}

	for points := []BFSPoint{{Point: sourceGalaxy}}; len(points) > 0; points = points[1:] {
		point := points[0]

		visitStatus[point.Point] = VisitStatusVisited

		for _, delta := range [4]Point2D{{X: 0, Y: -1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 1, Y: 0}} {
			nextPoint := Point2D{
				X: point.Point.X + delta.X,
				Y: point.Point.Y + delta.Y,
			}

			if !image.InBounds(nextPoint) || visitStatus[nextPoint] != VisitStatusNotVisited {
				continue
			}

			distance := point.Distance + 1
			if _, ok := image.ExpandedRows[nextPoint.Y]; ok && delta.Y != 0 {
				distance += expansionRate - 1
			}
			if _, ok := image.ExpandedColumns[nextPoint.X]; ok && delta.X != 0 {
				distance += expansionRate - 1
			}

			nextLocation := image.Locations[nextPoint.Y][nextPoint.X]
			if nextLocation.GalaxyID > sourceGalaxyID {
				res[nextLocation.GalaxyID-sourceGalaxyID-1] = distance

				delete(target, nextLocation.GalaxyID)
				if len(target) == 0 {
					return res
				}
			}

			visitStatus[nextPoint] = VisitStatusPendingVisit
			points = append(points, BFSPoint{
				Point:    nextPoint,
				Distance: distance,
			})
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
