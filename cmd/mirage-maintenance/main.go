package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
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
		var numbers []int
		for _, item := range strings.Split(scanner.Text(), " ") {
			number, err := strconv.Atoi(item)
			if err != nil {
				return fmt.Errorf("parse number %q: %w", item, err)
			}

			numbers = append(numbers, number)
		}

		sum += extrapolateOnePointBackwards(numbers)
	}

	fmt.Println(sum)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func extrapolateOnePointBackwards(numbers []int) int {
	differences := make([]int, len(numbers)-1)

	allZeroes := true
	for i := 1; i < len(numbers); i++ {
		difference := numbers[i] - numbers[i-1]
		if difference != 0 {
			allZeroes = false
		}

		differences[i-1] = difference
	}

	if allZeroes {
		return numbers[0]
	}

	next := extrapolateOnePointBackwards(differences)
	return numbers[0] - next
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
