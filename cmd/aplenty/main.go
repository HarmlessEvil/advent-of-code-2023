package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
)

type WorkflowContext map[string]Workflow

type Workflow struct {
	Rules []Rule
}

type Range struct {
	Min int
	Max int
}

func (r Range) Intersect(other Range) Range {
	return Range{
		Min: max(r.Min, other.Min),
		Max: min(r.Max, other.Max),
	}
}

func (r Range) Subtract(other Range) Range {
	if r.Min == other.Min {
		return Range{
			Min: other.Max,
			Max: r.Max,
		}
	}

	if r.Max == other.Max {
		return Range{
			Min: r.Min,
			Max: other.Min,
		}
	}

	panic(fmt.Errorf("not implemented in general case"))
}

type Condition struct {
	Category string
	Target   Range
}

type Rule struct {
	Condition *Condition
	Target    string
}

type PartPattern struct {
	X Range
	M Range
	A Range
	S Range
}

var MaxRange = Range{
	Min: 1,
	Max: 4001,
}

func runMain() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	workflowContext, err := parseWorkflows(scanner)
	if err != nil {
		return fmt.Errorf("parse workflow context: %w", err)
	}

	fmt.Println(countAcceptedParts(workflowContext, "in", PartPattern{
		X: MaxRange,
		M: MaxRange,
		A: MaxRange,
		S: MaxRange,
	}))

	return nil
}

func parseWorkflows(scanner *bufio.Scanner) (WorkflowContext, error) {
	res := WorkflowContext{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			break
		}

		i := bytes.IndexRune(line, '{')
		name := string(line[:i])

		var rules []Rule
		for _, rule := range bytes.Split(line[i+1:len(line)-1], []byte(",")) {
			parts := bytes.SplitN(rule, []byte(":"), 2)

			if len(parts) == 1 {
				rules = append(rules, Rule{
					Target: string(parts[0]),
				})
			} else {
				criteria := parts[0]

				targetValue := parseInt(criteria[2:])

				target := MaxRange
				switch criteria[1] {
				case '<':
					target.Max = targetValue
				case '>':
					target.Min = targetValue + 1
				}

				rules = append(rules, Rule{
					Target: string(parts[1]),
					Condition: &Condition{
						Category: string(criteria[0]),
						Target:   target,
					},
				})
			}
		}

		res[name] = Workflow{Rules: rules}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return res, nil
}

func parseInt(b []byte) int {
	res := 0
	for _, symbol := range b {
		res = res*10 + int(symbol-'0')
	}

	return res
}

func countAcceptedParts(c WorkflowContext, workflowName string, p PartPattern) int {
	if workflowName == "R" {
		return 0
	}

	if workflowName == "A" {
		return (p.X.Max - p.X.Min) * (p.M.Max - p.M.Min) * (p.A.Max - p.A.Min) * (p.S.Max - p.S.Min)
	}

	res := 0

	for _, rule := range c[workflowName].Rules {
		next := p

		if rule.Condition != nil {
			switch rule.Condition.Category {
			case "x":
				next.X = p.X.Intersect(rule.Condition.Target)
				if next.X.Max <= next.X.Min {
					continue
				}
			case "m":
				next.M = p.M.Intersect(rule.Condition.Target)
				if next.M.Max <= next.M.Min {
					continue
				}
			case "a":
				next.A = p.A.Intersect(rule.Condition.Target)
				if next.A.Max <= next.A.Min {
					continue
				}
			case "s":
				next.S = p.S.Intersect(rule.Condition.Target)
				if next.S.Max <= next.S.Min {
					continue
				}
			}
		}

		res += countAcceptedParts(c, rule.Target, next)

		if rule.Condition != nil {
			switch rule.Condition.Category {
			case "x":
				p.X = p.X.Intersect(MaxRange.Subtract(rule.Condition.Target))
			case "m":
				p.M = p.M.Intersect(MaxRange.Subtract(rule.Condition.Target))
			case "a":
				p.A = p.A.Intersect(MaxRange.Subtract(rule.Condition.Target))
			case "s":
				p.S = p.S.Intersect(MaxRange.Subtract(rule.Condition.Target))
			}
		}
	}

	return res
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
