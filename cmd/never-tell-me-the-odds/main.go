package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/mitchellh/go-z3"
)

type Point3D struct {
	X float64
	Y float64
	Z float64
}

func (p Point3D) String() string {
	return fmt.Sprintf("%.15g, %.15g, %.15g", p.X, p.Y, p.Z)
}

func (p Point3D) Add(other Point3D) Point3D {
	return Point3D{
		X: p.X + other.X,
		Y: p.Y + other.Y,
		Z: p.Z + other.Z,
	}
}

func (p Point3D) Sub(other Point3D) Point3D {
	return Point3D{
		X: p.X - other.X,
		Y: p.Y - other.Y,
		Z: p.Z - other.Z,
	}
}

func (p Point3D) Dot(other Point3D) float64 {
	return p.X*other.X + p.Y*other.Y + p.Z*other.Z
}

func (p Point3D) Mul(scalar float64) Point3D {
	return Point3D{
		X: p.X * scalar,
		Y: p.Y * scalar,
		Z: p.Z * scalar,
	}
}

type Line struct {
	Start     Point3D
	Direction Point3D
}

func (h Line) String() string {
	return fmt.Sprintf("%s @ %s", h.Start, h.Direction)
}

func runMain() error {
	hailstones, err := parseHailstones("input.txt")
	if err != nil {
		return fmt.Errorf("parse hailstones: %w", err)
	}

	config := z3.NewConfig()
	z3Ctx := z3.NewContext(config)
	config.Close()
	defer z3Ctx.Close()

	solver := z3Ctx.NewSolver()
	defer solver.Close()

	x := z3Ctx.Const(z3Ctx.Symbol("x"), z3Ctx.IntSort())
	y := z3Ctx.Const(z3Ctx.Symbol("y"), z3Ctx.IntSort())
	z := z3Ctx.Const(z3Ctx.Symbol("z"), z3Ctx.IntSort())

	dx := z3Ctx.Const(z3Ctx.Symbol("dx"), z3Ctx.IntSort())
	dy := z3Ctx.Const(z3Ctx.Symbol("dy"), z3Ctx.IntSort())
	dz := z3Ctx.Const(z3Ctx.Symbol("dz"), z3Ctx.IntSort())

	t := [3]*z3.AST{
		z3Ctx.Const(z3Ctx.Symbol("t1"), z3Ctx.IntSort()),
		z3Ctx.Const(z3Ctx.Symbol("t2"), z3Ctx.IntSort()),
		z3Ctx.Const(z3Ctx.Symbol("t3"), z3Ctx.IntSort()),
	}

	for i, hailstone := range hailstones[:3] {
		px := z3Ctx.Int64(int64(hailstone.Start.X), z3Ctx.IntSort())
		pdx := z3Ctx.Int64(int64(hailstone.Direction.X), z3Ctx.IntSort())
		solver.Assert(x.Add(dx.Mul(t[i])).Eq(px.Add(pdx.Mul(t[i]))))

		py := z3Ctx.Int64(int64(hailstone.Start.Y), z3Ctx.IntSort())
		pdy := z3Ctx.Int64(int64(hailstone.Direction.Y), z3Ctx.IntSort())
		solver.Assert(y.Add(dy.Mul(t[i])).Eq(py.Add(pdy.Mul(t[i]))))

		pz := z3Ctx.Int64(int64(hailstone.Start.Z), z3Ctx.IntSort())
		pdz := z3Ctx.Int64(int64(hailstone.Direction.Z), z3Ctx.IntSort())
		solver.Assert(z.Add(dz.Mul(t[i])).Eq(pz.Add(pdz.Mul(t[i]))))
	}

	if v := solver.Check(); v != z3.True {
		return fmt.Errorf("no solutions")
	}

	model := solver.Model()
	assignments := model.Assignments()
	_ = model.Close()

	fmt.Println(assignments["x"].Int64() + +assignments["y"].Int64() + +assignments["z"].Int64())

	return nil
}

func parseHailstones(filename string) ([]Line, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}
	defer f.Close()

	var hailstones []Line

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		hailstones = append(hailstones, parseHailstone(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return hailstones, nil
}

func parseHailstone(line string) Line {
	var hailstone Line
	_, _ = fmt.Sscanf(
		line,
		"%f, %f, %f @ %f, %f, %f",
		&hailstone.Start.X,
		&hailstone.Start.Y,
		&hailstone.Start.Z,
		&hailstone.Direction.X,
		&hailstone.Direction.Y,
		&hailstone.Direction.Z,
	)

	return hailstone
}

func dot(m, n, o, p Point3D) float64 {
	return m.Sub(n).Dot(o.Sub(p))
}

type Intersection struct {
	Point Point3D
	Value float64
}

// Source: https://stackoverflow.com/a/2316934/7149107
//
// Source: https://paulbourke.net/geometry/pointlineplane/
func intersect(a, b Line) (Intersection, Intersection) {
	p1 := a.Start
	p2 := p1.Add(a.Direction)

	p3 := b.Start
	p4 := p3.Add(b.Direction)

	mua := (dot(p1, p3, p4, p3)*dot(p4, p3, p2, p1) - dot(p1, p3, p2, p1)*dot(p4, p3, p4, p3)) /
		(dot(p2, p1, p2, p1)*dot(p4, p3, p4, p3) - dot(p4, p3, p2, p1)*dot(p4, p3, p2, p1))

	mub := (dot(p1, p3, p4, p3) + mua*dot(p4, p3, p2, p1)) / dot(p4, p3, p4, p3)

	return Intersection{
			Point: p1.Add(p2.Sub(p1).Mul(mua)),
			Value: mua,
		}, Intersection{
			Point: p3.Add(p4.Sub(p3).Mul(mub)),
			Value: mub,
		}
}

func main() {
	if err := runMain(); err != nil {
		slog.Error("program aborted", slog.Any("error", err))
		os.Exit(1)
	}
}
