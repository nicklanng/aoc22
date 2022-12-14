package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"strings"
)

const width, height = 800, 200

const (
	Air byte = iota
	Rock
	Sand
)

//go:embed input
var input []byte

func main() {
	var (
		tiles     [width * height]byte
		sandCount int
	)

	// part 1
	sandCount = 0
	tiles = parseInput(false)
	for simulate(tiles[:]) {
		sandCount++
	}
	fmt.Printf("Part 1 output: %d\n", sandCount)

	// part 2
	sandCount = 0
	tiles = parseInput(true)
	for simulate(tiles[:]) {
		sandCount++
	}
	fmt.Printf("Part 2 output: %d\n", sandCount+1)
}

func simulate(tiles []byte) bool {
	point := 500

	for {
		if point+width >= len(tiles) {
			return false
		}

		// try to move down
		if tiles[point+width] == Air {
			point += width
			continue
		}

		// FYI these will wrap around if on the edges of map, fine for this puzzle though

		// try to move down left
		if tiles[point+width-1] == Air {
			point += width - 1
			continue
		}
		// try to move down right
		if tiles[point+width+1] == Air {
			point += width + 1
			continue
		}

		tiles[point] = Sand
		return point != 500
	}
}

func parseInput(drawFloor bool) [width * height]byte {
	var tiles [width * height]byte
	var highestY int

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		// read line
		var points []lib.Vector2i
		pointsStr := strings.Split(scanner.Text(), " -> ")
		for _, str := range pointsStr {
			var point lib.Vector2i
			if err := point.ParseString(str); err != nil {
				panic(err)
			}
			points = append(points, point)
		}

		// draw rocks
		for i := 0; i < len(points)-1; i++ {
			start, end := points[i], points[i+1]

			if start.Y > highestY {
				highestY = start.Y
			}

			if start.Y == end.Y {
				// horizontal
				if end.X < start.X {
					end, start = start, end
				}
				for ; start.X <= end.X; start.X++ {
					tiles[start.ToIndex(width)] = Rock
				}
			} else {
				// vertical
				if end.Y < start.Y {
					end, start = start, end
				}
				for ; start.Y <= end.Y; start.Y++ {
					tiles[start.ToIndex(width)] = Rock
				}
			}
		}
	}

	if drawFloor {
		for x := 0; x <= width; x++ {
			tiles[(highestY+2)*width+x] = Rock
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return tiles
}
