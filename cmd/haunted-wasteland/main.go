package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	instructions := parseInstructions(scanner.Text())
	_ = instructions

	scanner.Scan() // \n

	var startNodes []string

	network := map[string][2]string{}
	for scanner.Scan() {
		line := scanner.Text()

		node := line[0:3]
		network[node] = [2]string{line[7:10], line[12:15]}

		if node[2] == 'A' {
			startNodes = append(startNodes, node)
		}
	}

	stepCount := make([]int, len(startNodes))
	for i, node := range startNodes {
		stepCount[i] = countSteps(network, instructions, node)
	}

	res := 1
	for _, count := range stepCount {
		res = lcm(res, count)
	}

	fmt.Println(res)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	return nil
}

func parseInstructions(text string) []int {
	instructions := make([]int, len(text))
	for i, direction := range text {
		switch direction {
		case 'L':
			instructions[i] = 0
		case 'R':
			instructions[i] = 1
		}
	}

	return instructions
}

func countSteps(network map[string][2]string, instructions []int, node string) int {
	stepCount := 0

	for i := 0; node[2] != 'Z'; i = (i + 1) % len(instructions) {
		node = network[node][instructions[i]]
		stepCount++
	}

	return stepCount
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
