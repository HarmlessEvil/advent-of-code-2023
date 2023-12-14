package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Tile int

const (
	TileEmptySpace Tile = iota
	TileRoundedRock
	TileCubeShapedRock
)

type Row []Tile

type Platform []Row

func runMain() error {
	platform, err := parsePlatform("input.txt")
	if err != nil {
		return fmt.Errorf("parse platform: %w", err)
	}

	for i := 0; i < 1_000; i++ {
		platform.tiltNorth()
		platform.tiltWest()
		platform.tiltSouth()
		platform.tiltEast()
	}

	fmt.Println(totalLoad(platform))

	return nil
}

func parsePlatform(filename string) (Platform, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("read input file: %w", err)
	}
	defer f.Close()

	var platform Platform

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		row := make(Row, len(line))
		for j, symbol := range line {
			var tile Tile
			switch symbol {
			case 'O':
				tile = TileRoundedRock
			case '#':
				tile = TileCubeShapedRock
			}

			row[j] = tile
		}

		platform = append(platform, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return platform, nil
}

func (p *Platform) tiltNorth() {
	cursors := make([]int, len((*p)[0]))

	for i, row := range *p {
		for j, tile := range row {
			if tile == TileCubeShapedRock {
				cursors[j] = i + 1
			} else if tile == TileRoundedRock {
				if cursors[j] != i {
					(*p)[cursors[j]][j] = TileRoundedRock
					(*p)[i][j] = TileEmptySpace
				}

				cursors[j]++
			}
		}
	}
}

func (p *Platform) tiltWest() {
	for i, row := range *p {
		cursor := 0

		for j, tile := range row {
			if tile == TileCubeShapedRock {
				cursor = j + 1
			} else if tile == TileRoundedRock {
				if cursor != j {
					(*p)[i][cursor] = TileRoundedRock
					(*p)[i][j] = TileEmptySpace
				}

				cursor++
			}
		}
	}
}

func (p *Platform) tiltSouth() {
	lastRowIndex := len(*p) - 1

	cursors := make([]int, len((*p)[0]))
	for i := range cursors {
		cursors[i] = lastRowIndex
	}

	for i := range *p {
		i = lastRowIndex - i

		for j, tile := range (*p)[i] {
			if tile == TileCubeShapedRock {
				cursors[j] = i - 1
			} else if tile == TileRoundedRock {
				if cursors[j] != i {
					(*p)[cursors[j]][j] = TileRoundedRock
					(*p)[i][j] = TileEmptySpace
				}

				cursors[j]--
			}
		}
	}
}

func (p *Platform) tiltEast() {
	for i, row := range *p {
		cursor := len(row) - 1

		for j := range row {
			j = len(row) - 1 - j
			tile := (*p)[i][j]

			if tile == TileCubeShapedRock {
				cursor = j - 1
			} else if tile == TileRoundedRock {
				if cursor != j {
					(*p)[i][cursor] = TileRoundedRock
					(*p)[i][j] = TileEmptySpace
				}

				cursor--
			}
		}
	}
}

func (r Row) String() string {
	var b strings.Builder
	for _, tile := range r {
		b.WriteRune(map[Tile]rune{
			TileEmptySpace:     '.',
			TileCubeShapedRock: '#',
			TileRoundedRock:    'O',
		}[tile])
	}

	return b.String()
}

func (p Platform) String() string {
	var b strings.Builder
	for i, row := range p {
		b.WriteString(row.String())

		if i != len(p)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func totalLoad(platform Platform) int {
	sum := 0

	for i, row := range platform {
		for _, tile := range row {
			if tile == TileRoundedRock {
				sum += len(platform) - i
			}
		}
	}

	return sum
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
