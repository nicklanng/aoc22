package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strconv"
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
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			currentElf += i
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Top single elf: %d\n", topElves[2])
	fmt.Printf("Top 3 elves combined: %d\n", topElves[0]+topElves[1]+topElves[2])
}
