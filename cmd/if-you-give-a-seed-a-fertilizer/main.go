package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
)

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	seeds, err := parseSeeds(strings.TrimPrefix(scanner.Text(), "seeds: "))
	if err != nil {
		return err
	}

	scanner.Scan() // \n
	scanner.Scan() // seed-to-soil map:
	seedToSoil, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse seed-to-soil map: %w", err)
	}

	soilRanges := lookupAllRanges(seedToSoil, seeds)

	scanner.Scan() // soil-to-fertilizer map:
	soilToFertilizer, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse soil-to-fertilizer map: %w", err)
	}

	fertilizerRanges := lookupAllRanges(soilToFertilizer, soilRanges)

	scanner.Scan() // fertilizer-to-water map:
	fertilizerToWater, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse fertilizer-to-water map: %w", err)
	}

	waterRanges := lookupAllRanges(fertilizerToWater, fertilizerRanges)

	scanner.Scan() // water-to-light map:
	waterToLight, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse water-to-light map: %w", err)
	}

	lightRanges := lookupAllRanges(waterToLight, waterRanges)

	scanner.Scan() // light-to-temperature map:
	lightToTemperature, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse light-to-temperature map: %w", err)
	}

	temperatureRanges := lookupAllRanges(lightToTemperature, lightRanges)

	scanner.Scan() // temperature-to-humidity map:
	temperatureToHumidity, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse temperature-to-humidity map: %w", err)
	}

	humidityRanges := lookupAllRanges(temperatureToHumidity, temperatureRanges)

	scanner.Scan() // humidity-to-location map:
	humidityToLocation, err := parseMap(scanner)
	if err != nil {
		return fmt.Errorf("parse humidity-to-location map: %w", err)
	}

	locationRanges := lookupAllRanges(humidityToLocation, humidityRanges)

	fmt.Println(locationRanges[0].Start)

	return nil
}

func lookupAllRanges(haystack []Range, needles []Range) []Range {
	var res []Range
	for _, needle := range needles {
		res = append(res, lookupRange(haystack, needle)...)
	}

	slices.SortFunc(res, func(a, b Range) int {
		return cmp.Compare(a.Start, b.Start)
	})

	return res
}

func parseSeeds(text string) ([]Range, error) {
	var seeds []Range
	parts := strings.Split(text, " ")
	for i := 0; i < len(parts); i += 2 {
		start, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("parse seed start %q: %w", parts[i], err)
		}

		length, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse seed length %q: %w", parts[i+1], err)
		}

		seeds = append(seeds, Range{
			Start: start,
			End:   start + length,
		})
	}

	slices.SortFunc(seeds, func(a, b Range) int {
		return cmp.Compare(a.Start, b.Start)
	})

	return seeds, nil
}

type Range struct {
	Start int
	End   int
	Value int
}

func parseMap(scanner *bufio.Scanner) ([]Range, error) {
	var res []Range

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		parts := strings.SplitN(line, " ", 3)

		destination, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("parse destination %q: %w", parts[0], err)
		}

		start, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("parse source %q: %w", parts[1], err)
		}

		length, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("parse length %q: %w", parts[2], err)
		}

		res = append(res, Range{
			Start: start,
			End:   start + length,
			Value: destination - start,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	slices.SortFunc(res, func(a, b Range) int {
		return cmp.Compare(a.Start, b.Start)
	})

	return res, nil
}

func lookupRange(haystack []Range, needle Range) []Range {
	var res []Range

	start, foundStart := slices.BinarySearchFunc(haystack, needle.Start, func(r Range, target int) int {
		if r.Start <= target && r.End > target {
			return 0
		}

		return cmp.Compare(r.Start, target)
	})

	end, foundEnd := slices.BinarySearchFunc(haystack, needle.End, func(r Range, target int) int {
		if r.Start <= target && r.End > target {
			return 0
		}

		return cmp.Compare(r.Start, target)
	})

	if !foundStart {
		gapEnd := needle.Start
		if start < len(haystack) {
			gapEnd = min(haystack[start].Start, gapEnd)
		}

		res = append(res, Range{
			Start: needle.Start,
			End:   gapEnd,
		})
	}

	for i := start; i < end; i++ {
		res = append(res, Range{
			Start: max(haystack[i].Start, needle.Start) + haystack[i].Value,
			End:   min(haystack[i].End, needle.End) + haystack[i].Value,
		})

		gapEnd := needle.End
		if i+1 < len(haystack) {
			gapEnd = min(haystack[i+1].Start, gapEnd)
		}

		if gapEnd > haystack[i].End {
			res = append(res, Range{
				Start: haystack[i].End,
				End:   gapEnd,
			})
		}
	}

	if foundEnd {
		res = append(res, Range{
			Start: max(haystack[end].Start, needle.Start) + haystack[end].Value,
			End:   needle.End + haystack[end].Value,
		})
	}

	slices.SortFunc(res, func(a, b Range) int {
		return cmp.Compare(a.Start, b.Start)
	})

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
