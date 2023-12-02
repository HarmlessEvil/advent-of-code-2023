package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	sum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		lineParts := strings.SplitN(line, ": ", 2)

		var gameID int
		if _, err := fmt.Sscanf(lineParts[0], "Game %d", &gameID); err != nil {
			return fmt.Errorf("parse game id in %q: %w", line, err)
		}

		cubes := make(map[string]int, 3)

		rounds := strings.Split(lineParts[1], "; ")
		for _, round := range rounds {
			sets := strings.Split(round, ", ")
			for _, set := range sets {
				var amount int
				var color string
				if _, err := fmt.Sscanf(set, "%d %s", &amount, &color); err != nil {
					return fmt.Errorf("parse set %q: %w", set, err)
				}

				cubes[color] = max(cubes[color], amount)
			}
		}

		power := 1
		for _, amount := range cubes {
			power *= amount
		}

		sum += power
	}

	fmt.Println(sum)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
