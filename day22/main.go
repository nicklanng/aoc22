package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"strings"
)

const (
	Void  = ' '
	Floor = '.'
	Wall  = '#'
)

type Direction byte

const (
	East Direction = iota
	South
	West
	North
)

func (d Direction) Left() Direction {
	return (d - 1 + 4) % 4
}

func (d Direction) Right() Direction {
	return (d + 1) % 4
}

func (d Direction) Vector() lib.Vector2i {
	return [4]lib.Vector2i{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}}[d]
}

type CmdForward struct {
	Distance int
}
type CmdTurnLeft struct{}
type CmdTurnRight struct{}

//go:embed input
var input []byte

func main() {
	width, height, tiles, commands := parseInput()

	position := lib.Vector2i{X: bytes.IndexByte(tiles, Void)}
	direction := East

	for _, cmd := range commands {
		switch v := cmd.(type) {
		case CmdTurnLeft:
			direction = direction.Left()
		case CmdTurnRight:
			direction = direction.Right()
		case CmdForward:
			for m := 0; m < v.Distance; m++ {
				nextPosition := moveWithinBoard(position, direction, width, height)
				targetTile := tiles[nextPosition.ToIndex(width)]

				for targetTile == Void {
					nextPosition = moveWithinBoard(nextPosition, direction, width, height)
					targetTile = tiles[nextPosition.ToIndex(width)]
				}

				if targetTile == Wall {
					break
				}

				position = nextPosition
			}
		}
	}

	partOneScore := 1000*(position.Y+1) + 4*(position.X+1) + int(direction)
	fmt.Println(partOneScore)
}

func moveWithinBoard(position lib.Vector2i, direction Direction, width int, height int) lib.Vector2i {
	nextPosition := position.Add(direction.Vector())
	if nextPosition.X <= 0 {
		nextPosition.X += width
	}
	if nextPosition.X >= width {
		nextPosition.X %= width
	}
	if nextPosition.Y <= 0 {
		nextPosition.Y += height
	}
	if nextPosition.Y >= height {
		nextPosition.Y %= height
	}
	return nextPosition
}

func parseInput() (int, int, []byte, []any) {
	var width int
	var tiles []byte
	var commands []any

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" {
			break
		}
		if len(s) > width {
			width = len(s)
		}
	}

	scanner = bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()
		switch s {
		case "":
			scanner.Scan()
			fields := splitOnAndAfter(scanner.Text(), "LR")
			for _, f := range fields {
				switch f {
				case "":
					continue
				case "L":
					commands = append(commands, CmdTurnLeft{})
				case "R":
					commands = append(commands, CmdTurnRight{})
				default:
					commands = append(commands, CmdForward{Distance: lib.MustParseInt(f)})
				}
			}
		default:
			if len(s) < width {
				s += strings.Repeat(" ", width-len(s))
			}
			tiles = append(tiles, s...)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return width, len(tiles) / width, tiles, commands
}

func splitOnAndAfter(s string, seps string) []string {
	var fields []string
	last := 0
	for i, r := range s {
		for _, sep := range seps {
			if r == sep {
				fields = append(fields, s[last:i])
				last = i
				fields = append(fields, s[last:i+1])
				last = i + 1
				break
			}
		}
	}
	fields = append(fields, s[last:len(s)])
	return fields
}
