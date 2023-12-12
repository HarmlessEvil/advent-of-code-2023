package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type CacheKey struct {
	Springs      string
	DamagedCount [30]int // abuse input property: it's guaranteed to have at most 6 * 5 numbers
}

type Cache map[CacheKey]int

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}
	defer f.Close()

	sum := 0
	cache := Cache{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, " ", 2)
		springs := parts[0]

		var damagedCount []int
		for _, item := range strings.Split(parts[1], ",") {
			count, err := strconv.Atoi(item)
			if err != nil {
				return fmt.Errorf("parse damaged count %q: %w", item, err)
			}

			damagedCount = append(damagedCount, count)
		}

		springs, damagedCount = unfold(springs, damagedCount)

		arrangements := countArrangements(cache, springs, toArray(damagedCount), len(damagedCount))
		fmt.Println(parts[0], damagedCount, arrangements)

		sum += arrangements
	}

	fmt.Println(sum)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func unfold(springs string, damagedCount []int) (string, []int) {
	unfoldedDamagedCount := make([]int, 0, len(damagedCount)*5)
	springsParts := make([]string, 5)

	for i := 0; i < 5; i++ {
		springsParts[i] = springs
		unfoldedDamagedCount = append(unfoldedDamagedCount, damagedCount...)
	}

	return strings.Join(springsParts, "?"), unfoldedDamagedCount
}

func toArray(s []int) [30]int {
	var res [30]int
	copy(res[:], s)

	return res
}

func countArrangements(cache Cache, springs string, damagedCount [30]int, groupCount int) int {
	if groupCount == 0 {
		if couldBeAllOperational(springs) {
			return 1
		}

		return 0
	}

	target := damagedCount[0]
	arrangements := 0

	var match bool
	for i := 0; i < len(springs)-target+1; i++ {
		match = true
		for _, spring := range springs[i : i+target] {
			if spring == '.' {
				match = false
				break
			}
		}

		isValid := i+target < len(springs) && springs[i+target] != '#'
		if match {
			logMatch(springs, i, target)

			if isValid {
				cacheKey := CacheKey{
					Springs:      springs[i+target+1:],
					DamagedCount: toArray(damagedCount[1:]),
				}

				arrangementsInSuffix, ok := cache[cacheKey]
				if !ok {
					arrangementsInSuffix = countArrangements(cache, cacheKey.Springs, cacheKey.DamagedCount, groupCount-1)
					cache[cacheKey] = arrangementsInSuffix
					slog.Debug("cache miss")
				} else {
					slog.Debug("cache hit")
				}

				if arrangementsInSuffix == 0 {
					isValid = false
				}

				arrangements += arrangementsInSuffix
			} else if i+target >= len(springs) && groupCount == 1 {
				arrangements++
			}
		}

		if springs[i] == '#' {
			if match && groupCount == 1 && arrangements == 0 && isValid {
				return 1
			}

			return arrangements
		}
	}

	if match && groupCount == 1 && arrangements == 0 {
		return 1
	}

	return arrangements
}

func logMatch(springs string, i int, target int) {
	slog.Debug(
		"matched",
		slog.String("springs", springs),
		slog.String("pattern", springs[i:i+target]),
		slog.String("match", springs[:i]+strings.Repeat("#", target)+springs[i+target:]),
		slog.Int("target", target),
	)
}

func couldBeAllOperational(springs string) bool {
	for _, spring := range springs {
		if spring == '#' {
			return false
		}
	}

	return true
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
