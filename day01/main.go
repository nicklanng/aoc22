package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"log"
	"sort"
)

//go:embed input
var input []byte

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	var (
		currentElf int
		topElves   [3]int
	)

	for scanner.Scan() {
		switch s := scanner.Text(); s {
		case "":
			if currentElf > topElves[0] {
				topElves[0] = currentElf
				sort.Ints(topElves[:])
			}
			currentElf = 0
		default:
			i := lib.MustParseInt(s)
			currentElf += i
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Top single elf: %d\n", topElves[2])
	fmt.Printf("Top 3 elves combined: %d\n", topElves[0]+topElves[1]+topElves[2])
}
