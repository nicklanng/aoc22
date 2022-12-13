package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
)

type Pair struct {
	left, right []any
}

type Result byte

const (
	ResultUnknown Result = iota
	ResultOrdered
	ResultUnordered
)

func sortPackets(left, right []any) Result {
	for i := range left {
		// right side ran out of items, so inputs are not in the right order
		if i >= len(right) {
			return ResultUnordered
		}

		switch v1 := left[i].(type) {
		case float64:
			switch v2 := right[i].(type) {
			case float64:
				// if both elements are numbers
				if v1 < v2 {
					return ResultOrdered
				}
				if v2 < v1 {
					return ResultUnordered
				}
			case []any:
				// mismatched types, convert left to slice and recurse
				slice1 := []any{v1}
				if result := sortPackets(slice1, v2); result != ResultUnknown {
					return result
				}
			}

		case []any:
			switch v2 := right[i].(type) {
			case float64:
				// mismatched types, convert right to slice and recurse
				slice2 := []any{v2}
				if result := sortPackets(v1, slice2); result != ResultUnknown {
					return result
				}
			case []any:
				// If both values are lists, recurse
				if result := sortPackets(v1, v2); result != ResultUnknown {
					return result
				}
			}
		}
	}

	// left side ran out of items, so inputs are in the right order
	if len(left) < len(right) {
		return ResultOrdered
	}

	// the two slices matches exactly
	return ResultUnknown
}

//go:embed input
var input []byte

func main() {
	pairs := parseInput()

	// part 1
	var part1 int
	for i, p := range pairs {
		if sortPackets(p.left, p.right) == ResultOrdered {
			part1 += i + 1
		}
	}
	fmt.Printf("Part 1 output: %d\n", part1)

	// part 2
	var allPackets [][]any
	for _, p := range pairs {
		allPackets = append(allPackets, p.left)
		allPackets = append(allPackets, p.right)
	}
	allPackets = append(allPackets, []any{[]any{2.0}})
	allPackets = append(allPackets, []any{[]any{6.0}})

	sort.Slice(allPackets, func(i, j int) bool {
		return sortPackets(allPackets[i], allPackets[j]) == ResultOrdered
	})

	var marker2, marker6 int
	marker2Str := fmt.Sprintf("%v", []any{[]any{2.0}})
	marker6Str := fmt.Sprintf("%v", []any{[]any{6.0}})

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

func parseInput() []Pair {
	var pairs []Pair

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		var s1 []any
		scanner.Scan()
		var s2 []any
		scanner.Scan()

		if err := json.Unmarshal([]byte(scanner.Text()), &s1); err != nil {
			panic(err)
		}
		if err := json.Unmarshal([]byte(scanner.Text()), &s2); err != nil {
			panic(err)
		}
		pairs = append(pairs, Pair{s1, s2})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return pairs
}
