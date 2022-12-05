package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
)

//go:embed input
var input []byte

func main() {
	var (
		partOneCount int
		partTwoCount int
	)

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := scanner.Text()

		var elf1Min, elf1Max, elf2Min, elf2Max int
		if _, err := fmt.Sscanf(s, "%d-%d,%d-%d", &elf1Min, &elf1Max, &elf2Min, &elf2Max); err != nil {
			panic(err)
		}

		assignment1 := assignment{min: elf1Min, max: elf1Max}
		assignment2 := assignment{min: elf2Min, max: elf2Max}

		if assignment1.contains(assignment2) || assignment2.contains(assignment1) {
			partOneCount++
		}

		if assignment1.intersects(assignment2) {
			partTwoCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1 Count: %d\n", partOneCount)
	fmt.Printf("Part 2 Count: %d\n", partTwoCount)
}

type assignment struct {
	min, max int
}

func (a assignment) contains(a2 assignment) bool {
	return a2.min >= a.min && a2.max <= a.max
}

func (a assignment) intersects(a2 assignment) bool {
	return !(a2.max < a.min || a2.min > a.max)
}
