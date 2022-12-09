package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math"
	"strings"
)

//go:embed input
var input []byte

type Vector2i struct {X, Y int}
func (v Vector2i) Add(v2 Vector2i) Vector2i     {return Vector2i{X: v.X+v2.X, Y: v.Y+v2.Y}}
func (v Vector2i) Sub(v2 Vector2i) Vector2i     {return Vector2i{X: v.X-v2.X, Y: v.Y-v2.Y}}
func (v Vector2i) Magnitude() float64           {return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))}
func (v Vector2i) Distance(v2 Vector2i) float64 {return v.Sub(v2).Magnitude()}
func (v Vector2i) UnitNormalize() Vector2i {
	if v.X != 0 {
		v.X = v.X / lib.IntAbs(v.X)
	}
	if v.Y != 0 {
		v.Y = v.Y / lib.IntAbs(v.Y)
	}
	return v
}

type MoveCommand struct {
	Direction Vector2i
	Distance  int
}

type Rope struct {
	Knots []Vector2i
}

func NewRope(length int) Rope {
	return Rope{Knots: make([]Vector2i, length)}
}

func (r *Rope) Move(direction Vector2i) {
	r.Knots[0] = r.Knots[0].Add(direction)
	for i := 0; i < len(r.Knots) -1 ; i++ {
		if r.Knots[i].Distance(r.Knots[i+1]) > 1.5 {
			r.Knots[i+1] = r.Knots[i+1].Add(r.Knots[i].Sub(r.Knots[i+1]).UnitNormalize())
		}
	}
}

func (r *Rope) Tail() Vector2i {
	return r.Knots[len(r.Knots)-1]
}

func stringToDirection(s string) Vector2i {
	switch s {
	case "U": return Vector2i{Y: 1}
	case "D": return Vector2i{Y: -1}
	case "L": return Vector2i{X: -1}
	case "R": return Vector2i{X: 1}
	default: panic("unknown direction")
	}
}

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	// scan input
	var commands []MoveCommand
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		direction := fields[0]
		distance := lib.MustParseInt(fields[1])
		commands = append(commands, MoveCommand{Direction: stringToDirection(direction), Distance: distance})
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	
	// part 1
	tailPositions := map[Vector2i]struct{}{}
	rope := NewRope(2)
	for _, cmd := range commands {
		for i := 0; i < cmd.Distance; i++ {
			rope.Move(cmd.Direction)
			tailPositions[rope.Tail()] = struct{}{}
		}
	}
	partOnePositionCount := len(tailPositions)
	tailPositions = map[Vector2i]struct{}{}

	// part 2
	partTwoTailPositions := map[Vector2i]struct{}{}
	rope = NewRope(10)
	for _, cmd := range commands {
		for i := 0; i < cmd.Distance; i++ {
			rope.Move(cmd.Direction)
			partTwoTailPositions[rope.Tail()] = struct{}{}
		}
	}
	partTwoPositionCount := len(partTwoTailPositions)

	fmt.Printf("Part 1 positions: %d\n", partOnePositionCount)
	fmt.Printf("Part 2 positions: %d\n", partTwoPositionCount)
}