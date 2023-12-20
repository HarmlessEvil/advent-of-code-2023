package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Pulse int

const (
	PulseLow Pulse = iota
	PulseHigh
)

func (p Pulse) String() string {
	return map[Pulse]string{
		PulseLow:  "-low-",
		PulseHigh: "-high-",
	}[p]
}

type Module interface {
	Receive(from string, pulse Pulse) *Pulse
}

type FlipFlop struct {
	IsActive bool
}

func (f *FlipFlop) Receive(_ string, pulse Pulse) *Pulse {
	if pulse == PulseHigh {
		return nil
	}

	f.IsActive = !f.IsActive

	var res Pulse
	if f.IsActive {
		res = PulseHigh
	}

	return &res
}

type Conjunction struct {
	Inputs map[string]Pulse
}

func (c *Conjunction) Receive(from string, pulse Pulse) *Pulse {
	c.Inputs[from] = pulse

	var res Pulse
	for _, p := range c.Inputs {
		if p != PulseHigh {
			res = PulseHigh
			return &res
		}
	}

	return &res
}

type Broadcaster struct{}

func (b Broadcaster) Receive(_ string, pulse Pulse) *Pulse {
	return &pulse
}

type WiredModule struct {
	Module  Module
	Outputs []string
}

type ModuleConfiguration map[string]WiredModule

type Transmission struct {
	Pulse       Pulse
	Source      string
	Destination string
}

func (t Transmission) String() string {
	return fmt.Sprintf("%s %s> %s", t.Source, t.Pulse, t.Destination)
}

func (c ModuleConfiguration) PushButton(targets map[Transmission]struct{}) *Transmission {
	var target *Transmission

	for queue := []Transmission{{Source: "button", Destination: "broadcaster"}}; len(queue) > 0; queue = queue[1:] {
		transmission := queue[0]

		if target == nil {
			if _, ok := targets[transmission]; ok {
				target = &transmission
			}
		}

		slog.Debug("sent", slog.Any("transmission", transmission))

		destination := c[transmission.Destination]
		if destination.Module == nil {
			continue
		}

		pulse := destination.Module.Receive(transmission.Source, transmission.Pulse)
		if pulse == nil {
			continue
		}

		for _, output := range destination.Outputs {
			queue = append(queue, Transmission{
				Pulse:       *pulse,
				Source:      transmission.Destination,
				Destination: output,
			})
		}
	}

	return target
}

// ButtonPresses is a maximum amount of 12-bit numbers. It will overflow all counters at least once.
const ButtonPresses = 1 << 12

func runMain() error {
	config, target, sources, err := parseModuleConfiguration("input.txt")
	if err != nil {
		return fmt.Errorf("parse module configuration: %w", err)
	}

	logger := slog.Default()

	transmissions := make(map[Transmission]struct{}, len(sources))
	for _, source := range sources {
		transmissions[Transmission{
			Pulse:       PulseHigh,
			Source:      source,
			Destination: target,
		}] = struct{}{}
	}

	results := make([]int, 0, len(sources))
	for i := 0; i < ButtonPresses; i++ {
		slog.SetDefault(logger.With(slog.Int("round", i)))

		if found := config.PushButton(transmissions); found != nil {
			results = append(results, i+1)

			if len(results) == cap(results) {
				break
			}
		}
	}

	res := 1
	for _, result := range results {
		res = lcm(res, result)
	}

	fmt.Println(res)

	return nil
}

func parseModuleConfiguration(filename string) (ModuleConfiguration, string, []string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", nil, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	config := ModuleConfiguration{}
	inputs := map[string][]string{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, " -> ", 2)

		name, module := parseModule(parts[0])
		outputs := strings.Split(parts[1], ", ")

		config[name] = WiredModule{
			Module:  module,
			Outputs: outputs,
		}

		for _, output := range outputs {
			inputs[output] = append(inputs[output], name)
		}
	}

	for name, modules := range inputs {
		if c, ok := config[name].Module.(*Conjunction); ok {
			c.Inputs = make(map[string]Pulse, len(modules))

			for _, module := range modules {
				c.Inputs[module] = PulseLow
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", nil, fmt.Errorf("scan: %w", err)
	}

	target := inputs["rx"][0] // rx has always 1 input
	return config, target, inputs[target], nil
}

func parseModule(s string) (string, Module) {
	if s[0] == '%' {
		return s[1:], &FlipFlop{}
	}

	if s[0] == '&' {
		return s[1:], &Conjunction{}
	}

	return s, Broadcaster{}
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
