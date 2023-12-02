package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Match struct {
	Index int
	Value int
}

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

		firstMatch := Match{Index: len(line)}
		lastMatch := Match{Index: -1}

		if i := strings.IndexAny(line, "123456789"); i != -1 && i < firstMatch.Index {
			firstMatch = Match{
				Index: i,
				Value: int(line[i] - '0'),
			}
		}

		if i := strings.LastIndexAny(line, "123456789"); i != -1 && i > lastMatch.Index {
			lastMatch = Match{
				Index: i,
				Value: int(line[i] - '0'),
			}
		}

		for i, digit := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"} {
			if j := strings.Index(line, digit); j != -1 && j < firstMatch.Index {
				firstMatch = Match{
					Index: j,
					Value: i + 1,
				}
			}

			if j := strings.LastIndex(line, digit); j != -1 && j > lastMatch.Index {
				lastMatch = Match{
					Index: j,
					Value: i + 1,
				}
			}
		}

		sum += firstMatch.Value*10 + lastMatch.Value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	fmt.Println(sum)

	return nil
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
