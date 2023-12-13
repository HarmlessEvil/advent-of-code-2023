package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

type Pattern int

const (
	PatternAsh Pattern = iota
	PatternRocks
)

type Note [][]Pattern

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	note, err := parseNote(scanner)
	if err != nil {
		return fmt.Errorf("parse note: %w", err)
	}

	sum := 0

	for len(note) > 0 {
		row := findReflectionRow(note)
		column := findReflectionColumn(note)

		sum += 100*row + column

		slog.Debug("reflection", slog.Int("row", row), slog.Int("column", column))

		note, err = parseNote(scanner)
		if err != nil {
			return fmt.Errorf("parse note: %w", err)
		}
	}

	fmt.Println(sum)

	return nil
}

func parseNote(scanner *bufio.Scanner) (Note, error) {
	var note Note

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		row := make([]Pattern, len(line))
		for i, pattern := range line {
			if pattern == '#' {
				row[i] = PatternRocks
			}
		}

		note = append(note, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return note, nil
}

func findReflectionColumn(note Note) int {
searchLine:
	for line := 1; line < len(note[0]); line++ {
		smudgeCount := 0

		for j := 0; j < min(line, len(note[0])-line); j++ {
			for i := range note {
				if note[i][line-j-1] != note[i][line+j] {
					if smudgeCount == 0 {
						smudgeCount++
					} else {
						continue searchLine
					}
				}
			}
		}

		if smudgeCount == 1 {
			return line
		}
	}

	return 0
}

func findReflectionRow(note Note) int {
searchLine:
	for line := 1; line < len(note); line++ {
		smudgeCount := 0

		for i := 0; i < min(line, len(note)-line); i++ {
			for j := range note[0] {
				if note[line-i-1][j] != note[line+i][j] {
					if smudgeCount == 0 {
						smudgeCount++
					} else {
						continue searchLine
					}
				}
			}
		}

		if smudgeCount == 1 {
			return line
		}
	}

	return 0
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
