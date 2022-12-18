package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
)

type Coordinate struct {
	x, y, z int
}

func (v *Coordinate) neighbors() []Coordinate {
	neighbors := []Coordinate{
		{v.x - 1, v.y, v.z},
		{v.x + 1, v.y, v.z},
		{v.x, v.y - 1, v.z},
		{v.x, v.y + 1, v.z},
		{v.x, v.y, v.z - 1},
		{v.x, v.y, v.z + 1},
	}
	return neighbors
}

func (v *Coordinate) inRange(min, max int) bool {
	return !(v.x < min || v.x > max || v.y < min || v.y > max || v.z < min || v.z > max)
}

const (
	AirVoxel byte = iota
	LavaVoxel
	ObsidianVoxel
)

//go:embed input
var input []byte

func main() {
	voxels := parseInput()

	// part 1
	exposedSides := countExposedSides(voxels)
	fmt.Printf("Part 1 output: %d\n", exposedSides)

	voxels = generateLava(voxels)
	exposedSides = countExposedSides(voxels)
	fmt.Printf("Part 2 output: %d\n", exposedSides)
}

func generateLava(voxels [30][30][30]byte) [30][30][30]byte {
	// replace all air with lava for now
	for x := 0; x < len(voxels); x++ {
		for y := 0; y < len(voxels[x]); y++ {
			for z := 0; z < len(voxels[x][y]); z++ {
				if voxels[x][y][z] != AirVoxel {
					continue
				}
				voxels[x][y][z] = LavaVoxel
			}
		}
	}

	// from the outside, fill in air
	queue := lib.Queue[Coordinate]{}
	for x := 0; x < len(voxels); x++ {
		for y := 0; y < len(voxels[0]); y++ {
			if x == 0 || x == len(voxels)-1 || y == 0 || y == len(voxels[0])-1 {
				queue.Push(Coordinate{x, y, 0})
				queue.Push(Coordinate{x, y, 29})
			}
		}
		queue.Push(Coordinate{x, 0, 0})
		queue.Push(Coordinate{x, 29, 0})
		queue.Push(Coordinate{x, 0, 29})
		queue.Push(Coordinate{x, 29, 29})
	}

	visited := map[Coordinate]bool{}
	for {
		c, ok := queue.Pop()
		if !ok {
			break
		}
		visited[c] = true

		if voxels[c.x][c.y][c.z] != LavaVoxel {
			continue
		}

		voxels[c.x][c.y][c.z] = AirVoxel
		for _, neighbor := range c.neighbors() {
			if neighbor.inRange(0, 29) && voxels[neighbor.x][neighbor.y][neighbor.z] == LavaVoxel && !visited[neighbor] {
				queue.Push(neighbor)
			}
		}
	}

	return voxels
}

func countExposedSides(voxels [30][30][30]byte) int {
	// Initialize a variable to keep track of the number of exposed sides
	exposedSides := 0

	// Iterate over the voxels in the voxels
	for x := 0; x < len(voxels); x++ {
		for y := 0; y < len(voxels[x]); y++ {
			for z := 0; z < len(voxels[x][y]); z++ {
				// Skip voxels that are air
				if voxels[x][y][z] != ObsidianVoxel {
					continue
				}

				// Check the voxel's neighbors and see if any of them are air voxels
				if x == 0 || voxels[x-1][y][z] == AirVoxel {
					exposedSides++
				}
				if x == len(voxels)-1 || voxels[x+1][y][z] == AirVoxel {
					exposedSides++
				}
				if y == 0 || voxels[x][y-1][z] == AirVoxel {
					exposedSides++
				}
				if y == len(voxels[x])-1 || voxels[x][y+1][z] == AirVoxel {
					exposedSides++
				}
				if z == 0 || voxels[x][y][z-1] == AirVoxel {
					exposedSides++
				}
				if z == len(voxels[x][y])-1 || voxels[x][y][z+1] == AirVoxel {
					exposedSides++
				}
			}
		}
	}

	// Return the number of exposed sides
	return exposedSides
}

func parseInput() [30][30][30]byte {
	var voxels [30][30][30]byte

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		s := scanner.Text()
		var x, y, z int
		fmt.Sscanf(s, "%d,%d,%d", &x, &y, &z)
		voxels[x][y][z] = ObsidianVoxel
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return voxels
}
