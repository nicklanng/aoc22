package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"strings"
)

//go:embed input
var input []byte

type MoveCommand struct {
	Direction lib.Vector2i
	Distance  int
}

type Rope struct {
	Knots []lib.Vector2i
}

func NewRope(length int) Rope {
	return Rope{Knots: make([]lib.Vector2i, length)}
}

func (r *Rope) Move(direction lib.Vector2i) {
	r.Knots[0] = r.Knots[0].Add(direction)
	for i := 0; i < len(r.Knots) -1 ; i++ {
		if r.Knots[i].Distance(r.Knots[i+1]) > 1.5 {
			r.Knots[i+1] = r.Knots[i+1].Add(r.Knots[i].Sub(r.Knots[i+1]).UnitNormalize())
		}
	}
}

func (r *Rope) Tail() lib.Vector2i {
	return r.Knots[len(r.Knots)-1]
}

func stringToDirection(s string) lib.Vector2i {
	switch s {
	case "U": return lib.Vector2i{Y: 1}
	case "D": return lib.Vector2i{Y: -1}
	case "L": return lib.Vector2i{X: -1}
	case "R": return lib.Vector2i{X: 1}
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
	tailPositions := map[lib.Vector2i]struct{}{}
	rope := NewRope(2)
	for _, cmd := range commands {
		for i := 0; i < cmd.Distance; i++ {
			rope.Move(cmd.Direction)
			tailPositions[rope.Tail()] = struct{}{}
		}
	}
	partOnePositionCount := len(tailPositions)
	tailPositions = map[lib.Vector2i]struct{}{}

	// part 2
	rope = NewRope(10)
	for _, cmd := range commands {
		for i := 0; i < cmd.Distance; i++ {
			rope.Move(cmd.Direction)
			tailPositions[rope.Tail()] = struct{}{}
		}
	}
	partTwoPositionCount := len(tailPositions)

	fmt.Printf("Part 1 positions: %d\n", partOnePositionCount)
	fmt.Printf("Part 2 positions: %d\n", partTwoPositionCount)
}