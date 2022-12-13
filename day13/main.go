package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

type Packet []any

func (p Packet) String() string {
	return fmt.Sprintf("%v", []any(p))
}

type PacketList []Packet

func (p PacketList) Len() int {
	return len(p)
}
func (p PacketList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PacketList) Less(i, j int) bool {
	return sortPackets(p[i], p[j]) == -1
}

func sortPackets(left, right []any) int {
	for i := range left {
		// right side ran out of items, so inputs are not in the right order
		if i >= len(right) {
			return 1
		}

		switch v1 := left[i].(type) {
		case float64:
			switch v2 := right[i].(type) {
			case float64:
				// if both elements are numbers
				if v1 < v2 {
					return -1
				}
				if v2 < v1 {
					return 1
				}
			case []any:
				// mismatched types, convert left to slice and recurse
				slice1 := []any{v1}
				if result := sortPackets(slice1, v2); result != 0 {
					return result
				}
			}

		case []any:
			switch v2 := right[i].(type) {
			case float64:
				// mismatched types, convert right to slice and recurse
				slice2 := []any{v2}
				if result := sortPackets(v1, slice2); result != 0 {
					return result
				}
			case []any:
				// If both values are lists, recurse
				if result := sortPackets(v1, v2); result != 0 {
					return result
				}
			}
		}
	}

	// left side ran out of items, so inputs are in the right order
	if len(left) < len(right) {
		return -1
	}

	// the two slices matches exactly
	return 0
}

//go:embed input
var input []byte

func main() {
	pairs := parseInput()

	// part 1
	var part1 int
	for i, p := range pairs {
		if sortPackets(p[0], p[1]) == -1 {
			part1 += i + 1
		}
	}
	fmt.Printf("Part 1 output: %d\n", part1)

	// part 2
	allPackets := PacketList{[]any{[]any{2.0}}, []any{[]any{6.0}}}
	for _, p := range pairs {
		allPackets = append(allPackets, p...)
	}

	var marker2, marker6 int
	marker2Str := allPackets[0].String()
	marker6Str := allPackets[1].String()

	sort.Sort(allPackets)

	for i, p := range allPackets {
		if fmt.Sprintf("%v", p) == marker2Str {
			marker2 = i + 1
		}
		if fmt.Sprintf("%v", p) == marker6Str {
			marker6 = i + 1
		}
	}

	fmt.Printf("Part 2 output: %d\n", marker2*marker6)
}

func parseInput() []PacketList {
	var pairs []PacketList

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		var s1 []any
		if err := json.Unmarshal([]byte(scanner.Text()), &s1); err != nil {
			panic(err)
		}

		scanner.Scan()
		var s2 []any
		if err := json.Unmarshal([]byte(scanner.Text()), &s2); err != nil {
			panic(err)
		}
		pairs = append(pairs, PacketList{s1, s2})

		scanner.Scan()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return pairs
}
