package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"math/big"
)

//go:embed input
var input []byte

func main() {
	cypher := parseInput()

	partOne := sumFirstCoords(decrypt(cypher, 1, 1))
	fmt.Println(partOne)

	partTwo := sumFirstCoords(decrypt(cypher, 811589153, 10))
	fmt.Println(partTwo)
}

func decrypt(file []int, scalar int, runs int) []int {
	for i := range file {
		file[i] *= scalar
	}

	indices := make([]int, len(file))
	for i := range indices {
		indices[i] = i
	}

	mixed := make([]int, len(file))
	copy(mixed, file)

	for k := 0; k < runs; k++ {
		for i, num := range file {
			if num == 0 {
				continue
			}

			currentPos, _ := lib.Find(indices, i)

			x := big.NewInt(int64(currentPos + num))
			newPos := x.Mod(x, big.NewInt(int64(len(file)-1))).Int64()

			if int(newPos) > currentPos {
				copy(mixed[currentPos:newPos], mixed[currentPos+1:newPos+1])
				mixed[newPos] = num
				copy(indices[currentPos:newPos], indices[currentPos+1:newPos+1])
				indices[newPos] = i
			} else {
				copy(mixed[newPos+1:currentPos+1], mixed[newPos:currentPos])
				mixed[newPos] = num
				copy(indices[newPos+1:currentPos+1], indices[newPos:currentPos])
				indices[newPos] = i
			}

		}
	}
	return mixed
}

func sumFirstCoords(decrypted []int) int {
	// Find the index of the first 0 in the mixed list
	index, _ := lib.Find(decrypted, 0)
	coord1 := decrypted[(index+1000)%len(decrypted)]
	coord2 := decrypted[(index+2000)%len(decrypted)]
	coord3 := decrypted[(index+3000)%len(decrypted)]

	// Calculate the sum of the three numbers
	sum := coord1 + coord2 + coord3
	return sum
}

func parseInput() []int {
	var cypher []int

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		cypher = append(cypher, lib.MustParseInt(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return cypher
}
