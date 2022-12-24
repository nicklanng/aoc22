package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math"
)

//go:embed input
var input []byte

func main() {

	// part 1

	elves := parseInput()
	//debugPrint(elves)

	for i := 0; i < 10; i++ {
		planMoves(elves, i)
		elves = move(elves)
		//debugPrint(elves)
	}

	// find minX, minY, maxX, maxY
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	for _, v := range elves {
		if v.X < minX {
			minX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}
	//debugPrint(elves)
	areaOfRectangle := (1 + maxX - minX) * (1 + maxY - minY)
	fmt.Printf("Part 1: %d\n", areaOfRectangle-len(elves))

	// part 2
	elves = parseInput()

	i := 0
	for {
		static := planMoves(elves, i)
		if static {
			break
		}
		elves = move(elves)
		i++
	}
	fmt.Printf("Part 2: %d\n", i+1)
}

func debugPrint(elves map[lib.Vector2i]lib.Vector2i) {
	for y := 0; y < 12; y++ {
		for x := 0; x < 14; x++ {
			if _, ok := elves[lib.Vector2i{X: x, Y: y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func planMoves(elves map[lib.Vector2i]lib.Vector2i, runNumber int) bool {
	// propose moves
	for k := range elves {
		neighbors := elves[k].Neighbors()

		_, nw := elves[neighbors[0]]
		_, n := elves[neighbors[1]]
		_, ne := elves[neighbors[2]]
		_, w := elves[neighbors[3]]
		_, e := elves[neighbors[4]]
		_, sw := elves[neighbors[5]]
		_, s := elves[neighbors[6]]
		_, se := elves[neighbors[7]]

		if !nw && !n && !ne && !w && !e && !sw && !s && !se {
			continue
		}

		var moved bool
		for i := 0; i < 4; i++ {
			if moved {
				break
			}

			switch (i + runNumber) % 4 {
			case 0:
				if !nw && !n && !ne {
					elves[k] = neighbors[1]
					moved = true
				}
			case 1:
				if !sw && !s && !se {
					elves[k] = neighbors[6]
					moved = true
				}
			case 2:
				if !nw && !w && !sw {
					elves[k] = neighbors[3]
					moved = true
				}
			case 3:
				if !ne && !e && !se {
					elves[k] = neighbors[4]
					moved = true
				}
			}
		}

	}

	// count elves that are moving to the same space
	requestedMoveCount := make(map[lib.Vector2i]int)
	for _, dest := range elves {
		count, ok := requestedMoveCount[dest]
		if !ok {
			requestedMoveCount[dest] = 1
		} else {
			requestedMoveCount[dest] = count + 1
		}
	}

	// stop elves that are moving to the same space
	for k, v := range elves {
		if count := requestedMoveCount[v]; count > 1 {
			elves[k] = k
		}
	}

	static := true
	for k, v := range elves {
		if k != v {
			static = false
			break
		}
	}

	return static
}

func move(elves map[lib.Vector2i]lib.Vector2i) map[lib.Vector2i]lib.Vector2i {
	newElves := make(map[lib.Vector2i]lib.Vector2i)
	for _, v := range elves {
		newElves[v] = v
	}
	return newElves
}

func parseInput() map[lib.Vector2i]lib.Vector2i {
	elves := make(map[lib.Vector2i]lib.Vector2i)

	var y int

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()
		for x := range s {
			if s[x] == '#' {
				elves[lib.Vector2i{X: x, Y: y}] = lib.Vector2i{X: x, Y: y}
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return elves
}
