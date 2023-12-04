package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const filename = "input.txt"

func runMain() error {
	cardsAmount, err := countFileLines(filename)
	if err != nil {
		return fmt.Errorf("count cards: %w", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	cardsTotal := 0

	cardCopies := make([]int, cardsAmount)
	for i := range cardCopies {
		cardCopies[i] = 1
	}

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		lineParts := strings.SplitN(line, ": ", 2)

		scratchcardParts := strings.SplitN(lineParts[1], " | ", 2)

		winningNumbers := make(map[int]struct{})
		for _, item := range strings.Fields(scratchcardParts[0]) {
			number, err := strconv.Atoi(item)
			if err != nil {
				return fmt.Errorf("parse winning number %q in line %q: %w", item, line, err)
			}

			winningNumbers[number] = struct{}{}
		}

		numbersMatchedAmount := 0
		for _, item := range strings.Fields(scratchcardParts[1]) {
			number, err := strconv.Atoi(item)
			if err != nil {
				return fmt.Errorf("parse picked number %q in line %q: %w", item, line, err)
			}

			if _, ok := winningNumbers[number]; ok {
				numbersMatchedAmount++
			}
		}

		for j := 0; j < numbersMatchedAmount; j++ {
			cardCopies[i+j+1] += cardCopies[i]

			slog.Debug(
				"has winning number",
				slog.String("id", lineParts[0]),
				slog.Int("copies", cardCopies[i]),
				slog.Int("matches", numbersMatchedAmount),
			)
		}

		cardsTotal += cardCopies[i]
	}

	fmt.Println(cardsTotal)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func countFileLines(filename string) (int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSeparator := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSeparator)

		switch {
		case errors.Is(err, io.EOF):
			return count, nil

		case err != nil:
			return 0, err
		}
	}
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
