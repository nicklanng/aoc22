package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
)

//go:embed input
var input []byte

func main() {
	startPos, endPos, width, heightMap := parseMap()

	// part 1
	path := breadthFirstSearch(startPos, heightMap, width, func(next, cur int) bool { return next-cur <= 1}, func(pos lib.Vector2i, height int) bool {
		return pos == endPos
	})

	start, _ := path[endPos]
	steps := countSteps(path, start)
	fmt.Printf("Part 1 output: %d\n", steps)

	// part 2
	var foundPos lib.Vector2i
	path = breadthFirstSearch(endPos, heightMap, width, func(next, cur int) bool { return cur-next <= 1}, func(pos lib.Vector2i, height int) bool {
		if height == 0 {
			foundPos = pos
			return true
		}
		return false
	})

	start, _ = path[foundPos]
	steps = countSteps(path, start)
	fmt.Printf("Part 1 output: %d\n", steps)
}


func parseMap() (startingPos lib.Vector2i, finalPos lib.Vector2i, width int, heightMap []int){
	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		s := scanner.Text()

		width = len(s)
		for i := 0; i < width; i++ {
			switch s[i] {
			case 'S':
				startingPos = lib.Vector2i{X: i, Y: (len(heightMap)+1)/width}
				heightMap = append(heightMap, 'a'-'a')
			case 'E':
				finalPos = lib.Vector2i{X: i, Y: (len(heightMap)+1)/width}
				heightMap = append(heightMap, 'z'-'a')
			default:
				heightMap = append(heightMap, int(s[i]-'a'))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return startingPos, finalPos, width, heightMap
}

func breadthFirstSearch(
	startPos lib.Vector2i,
	heightMap []int,
	width int,
	passable func(neighborHeight int, currentHeight int) bool,
	stopSearch func(pos lib.Vector2i, height int) bool,
) map[lib.Vector2i]lib.Vector2i {

	height := len(heightMap)/width

	frontier := lib.Queue[lib.Vector2i]{}
	frontier.Push(startPos)

	cameFrom := map[lib.Vector2i]lib.Vector2i{}
	cameFrom[startPos] = lib.Vector2i{X: -1, Y: -1}

	for {
		current, ok := frontier.Pop()
		if !ok {
			break
		}

		currentHeight := heightMap[current.Y*width+current.X]

		neighbors := [...]lib.Vector2i{
			current.Add(lib.Vector2i{X: 0, Y: -1}),
			current.Add(lib.Vector2i{X: 0, Y: 1}),
			current.Add(lib.Vector2i{X: -1, Y: 0}),
			current.Add(lib.Vector2i{X: 1, Y: 0}),
		}

		for _, n := range neighbors {
			if n.X < 0 || n.X >= width || n.Y < 0 || n.Y >= height {
				continue
			}

			neighborHeight := heightMap[n.Y*width+n.X]
			if !passable(neighborHeight, currentHeight) {
				continue
			}

			if _, ok := cameFrom[n]; ok {
				continue
			}

			if stopSearch(n, neighborHeight) {
				cameFrom[n] = current
				return cameFrom
			}

			frontier.Push(n)
			cameFrom[n] = current
		}
	}
	return cameFrom
}


func countSteps(path map[lib.Vector2i]lib.Vector2i, start lib.Vector2i) int {
	current := start
	var steps int
	for {
		steps++
		prev, _ := path[current]
		if prev.X == -1 {
			break
		}
		current = prev
	}
	return steps
}
