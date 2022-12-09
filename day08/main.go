package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
)

//go:embed input
var input []byte

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	var (
		width, height int
		trees         []byte
	)

	// scan input
	for scanner.Scan() {
		height++
		s := scanner.Text()
		width = len(s)
		for i := 0; i < len(s); i++ {
			trees = append(trees, s[i])
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	partOne(width, height, trees)
	partTwo(width, height, trees)

}

func partOne(width, height int, trees []byte) {
	visible := make([]bool, len(trees))

	// left to right
	for y := 0; y < height; y++ {
		var highest byte
		for x := 0; x < width; x++ {
			index := y*width + x
			if trees[index] > highest {
				highest = trees[index]
				visible[index] = true
			}
		}
	}

	// right to left
	for y := 0; y < height; y++ {
		var highest byte
		for x := width - 1; x >= 0; x-- {
			index := y*width + x
			if trees[index] > highest {
				highest = trees[index]
				visible[index] = true
			}
		}
	}

	// top to bottom
	for x := 0; x < width; x++ {
		var highest byte
		for y := 0; y < height; y++ {
			index := y*width + x
			if trees[index] > highest {
				highest = trees[index]
				visible[index] = true
			}
		}
	}

	// bottom to top
	for x := 0; x < width; x++ {
		var highest byte
		for y := height - 1; y >= 0; y-- {
			index := y*width + x
			if trees[index] > highest {
				highest = trees[index]
				visible[index] = true
			}
		}
	}

	var count int
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			if visible[index] {
				count++
			}
		}
	}
	fmt.Printf("Part one count: %d\n", count)
}

func partTwo(width int, height int, trees []byte) {
	var maxSceneScore int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var right, left, down, up int
			treeHeight := trees[y*width + x]

			// right
			if x < width -1 {
				for dx := 1; dx < width-x; dx++ {
					right = dx
					index := y*width + x + dx
					if trees[index] >= treeHeight {
						break
					}
				}
			}

			// left
			if x > 0 {
				for dx := 1; dx <= x ; dx++ {
					left = dx
					index := y*width + x - dx
					if trees[index] >= treeHeight {
						break
					}
				}
			}

			// down
			if y < height - 1 {
				for dy := 1; dy < height - y ; dy++ {
					down = dy
					index := (y+dy)*width + x
					if trees[index] >= treeHeight {
						break
					}
				}
			}

			// up
			if y > 0 {
				for dy := 1; dy <= y ; dy++ {
					up = dy
					index := (y-dy)*width + x
					if trees[index] >= treeHeight {
						break
					}
				}
			}

			if  right * left * down * up > maxSceneScore {
				maxSceneScore = right * left * down * up
			}
		}
	}

	fmt.Printf("Part two score: %d\n", maxSceneScore)
}