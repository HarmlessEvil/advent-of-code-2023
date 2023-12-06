package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strings"
	"unicode"
)

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	raceTime := parseNumber(strings.TrimPrefix(scanner.Text(), "Time:"))

	scanner.Scan()
	distance := parseNumber(strings.TrimPrefix(scanner.Text(), "Distance:"))

	d := raceTime*raceTime - 4*distance
	minT := math.Ceil((float64(raceTime) - math.Sqrt(float64(d))) / 2)
	maxT := math.Ceil((float64(raceTime) + math.Sqrt(float64(d))) / 2)

	fmt.Println(int(maxT - minT))

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func parseNumber(text string) int {
	res := 0
	for _, char := range text {
		if !unicode.IsSpace(char) {
			res = res*10 + int(char-'0')
		}
	}

	return res
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
