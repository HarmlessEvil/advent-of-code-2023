package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"unicode"
)

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var prevLine []rune
	var currentLine []rune

	scanner.Scan()
	nextLine := []rune(scanner.Text())

	sum := 0

	for scanner.Scan() {
		prevLine = currentLine
		currentLine = nextLine
		nextLine = []rune(scanner.Text())

		sum += sumOfGearRatiosInLine(prevLine, currentLine, nextLine)
	}

	sum += sumOfGearRatiosInLine(currentLine, nextLine, nil)

	fmt.Println(sum)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

type Range struct {
	Start int
	End   int
}

func sumOfGearRatiosInLine(prevLine []rune, currentLine []rune, nextLine []rune) int {
	sum := 0

	for i, char := range currentLine {
		if char != '*' {
			continue
		}

		adjacentDigitsAmount := 0
		adjacentDigitRanges := [3][]Range{
			make([]Range, 0, 2),
			make([]Range, 0, 2),
			make([]Range, 0, 2),
		}

		if len(prevLine) > 0 {
			if i > 0 && unicode.IsDigit(prevLine[i-1]) {
				adjacentDigitRanges[0] = append(adjacentDigitRanges[0], Range{
					Start: i - 1,
					End:   i,
				})
			}

			if unicode.IsDigit(prevLine[i]) {
				if len(adjacentDigitRanges[0]) == 0 {
					adjacentDigitRanges[0] = append(adjacentDigitRanges[0], Range{
						Start: i,
						End:   i + 1,
					})
				} else {
					adjacentDigitRanges[0][0].End++
				}
			}

			if i+1 < len(prevLine) && unicode.IsDigit(prevLine[i+1]) {
				if len(adjacentDigitRanges[0]) == 0 || adjacentDigitRanges[0][0].End != i+1 {
					adjacentDigitRanges[0] = append(adjacentDigitRanges[0], Range{
						Start: i + 1,
						End:   i + 2,
					})
				} else {
					adjacentDigitRanges[0][0].End++
				}
			}

			adjacentDigitsAmount += len(adjacentDigitRanges[0])
		}

		if i > 0 && unicode.IsDigit(currentLine[i-1]) {
			adjacentDigitRanges[1] = append(adjacentDigitRanges[1], Range{
				Start: i - 1,
				End:   i,
			})

			adjacentDigitsAmount++
		}

		if i+1 < len(currentLine) && unicode.IsDigit(currentLine[i+1]) {
			adjacentDigitRanges[1] = append(adjacentDigitRanges[1], Range{
				Start: i + 1,
				End:   i + 2,
			})

			adjacentDigitsAmount++
		}

		if len(nextLine) > 0 {
			if i > 0 && unicode.IsDigit(nextLine[i-1]) {
				adjacentDigitRanges[2] = append(adjacentDigitRanges[2], Range{
					Start: i - 1,
					End:   i,
				})
			}

			if unicode.IsDigit(nextLine[i]) {
				if len(adjacentDigitRanges[2]) == 0 {
					adjacentDigitRanges[2] = append(adjacentDigitRanges[2], Range{
						Start: i,
						End:   i + 1,
					})
				} else {
					adjacentDigitRanges[2][0].End++
				}
			}

			if i+1 < len(nextLine) && unicode.IsDigit(nextLine[i+1]) {
				if len(adjacentDigitRanges[2]) == 0 || adjacentDigitRanges[2][0].End != i+1 {
					adjacentDigitRanges[2] = append(adjacentDigitRanges[2], Range{
						Start: i + 1,
						End:   i + 2,
					})
				} else {
					adjacentDigitRanges[2][0].End++
				}
			}

			adjacentDigitsAmount += len(adjacentDigitRanges[2])
		}

		if adjacentDigitsAmount != 2 {
			continue
		}

		gearRatio := 1
		for j, numberRanges := range adjacentDigitRanges {
			line := [3][]rune{prevLine, currentLine, nextLine}[j]

			for _, numberRange := range numberRanges {
				gearRatio *= parseNumber(line, numberRange.Start, numberRange.End)
			}
		}

		sum += gearRatio
	}

	return sum
}

func parseNumber(line []rune, start int, end int) int {
	for ; start > 0 && unicode.IsDigit(line[start-1]); start-- {
	}
	for ; end < len(line) && unicode.IsDigit(line[end]); end++ {
	}

	number := 0
	for i := start; i < end; i++ {
		number = number*10 + int(line[i]-'0')
	}

	slog.Debug(
		"parse number",
		slog.Int("start", start),
		slog.Int("end", end),
		slog.String("line", string(line[start:end])),
		slog.Int("number", number),
	)

	return number
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
