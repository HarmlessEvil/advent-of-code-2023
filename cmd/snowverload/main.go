package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Wiring map[string]map[string]struct{}

func runMain() error {
	wiring, err := parseWiring("input.txt")
	if err != nil {
		return fmt.Errorf("parse wiring: %w", err)
	}

	var cut []Edge
	for len(cut) != 3 {
		cut = contract(wiring)
	}

	fmt.Println(cut)

	for _, edge := range cut {
		delete(wiring[edge.From], edge.To)
		delete(wiring[edge.To], edge.From)
	}

	edge := cut[0]

	count1 := dfs(wiring, map[string]struct{}{edge.From: {}}, edge.From)
	count2 := dfs(wiring, map[string]struct{}{edge.To: {}}, edge.To)

	fmt.Println(count1, count2, count1*count2)

	return nil
}

func parseWiring(filename string) (Wiring, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	wiring := Wiring{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		from := line[:3]

		for _, to := range strings.Split(line[5:], " ") {
			if _, ok := wiring[from]; !ok {
				wiring[from] = map[string]struct{}{}
			}
			wiring[from][to] = struct{}{}

			if _, ok := wiring[to]; !ok {
				wiring[to] = map[string]struct{}{}
			}
			wiring[to][from] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return wiring, nil
}

type Edge struct {
	From string
	To   string
}

func contract(wiring Wiring) []Edge {
	graph := make(map[string]map[string][]Edge, len(wiring))
	for src, dst := range wiring {
		graph[src] = make(map[string][]Edge, len(dst))

		for to := range dst {
			graph[src][to] = append(graph[src][to], Edge{
				From: src,
				To:   to,
			})
		}
	}

	for len(graph) > 2 {
		from, to := getRandomEdge(graph)

		next := from + to

		for dest := range graph[from] {
			if dest == to {
				continue
			}

			edges := graph[dest][from]

			delete(graph[dest], from)
			graph[dest][next] = append(graph[dest][next], edges...)

			if _, ok := graph[next]; !ok {
				graph[next] = map[string][]Edge{}
			}
			graph[next][dest] = append(graph[next][dest], edges...)
		}
		delete(graph, from)

		for dest := range graph[to] {
			if dest == from {
				continue
			}

			edges := graph[dest][to]

			delete(graph[dest], to)
			graph[dest][next] = append(graph[dest][next], edges...)

			if _, ok := graph[next]; !ok {
				graph[next] = map[string][]Edge{}
			}
			graph[next][dest] = append(graph[next][dest], edges...)
		}
		delete(graph, to)
	}

	for _, vertices := range graph {
		for _, edges := range vertices {
			return edges
		}
	}

	return nil
}

func getRandomEdge(graph map[string]map[string][]Edge) (string, string) {
	for from, vertices := range graph {
		for to := range vertices {
			return from, to
		}
	}

	return "", ""
}

func dfs(wiring Wiring, visited map[string]struct{}, vertex string) int {
	res := 1

	for to := range wiring[vertex] {
		if _, ok := visited[to]; ok {
			continue
		}

		visited[to] = struct{}{}
		res += dfs(wiring, visited, to)
	}

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
