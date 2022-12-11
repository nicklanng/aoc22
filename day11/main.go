package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"sort"
	"strings"
)

type Operation func(old int) int

type Monkey struct {
	MonkeyID  int
	Items            []int
	InspectOperation Operation
	TestDivisor      int
	TestTrueTarget int
	TestFalseTarget int

	inspectCount int
}

//go:embed input
var input []byte

func main() {
	// part 1
	monkeys := parseMonkeys()
	for roundIndex := 0; roundIndex < 20; roundIndex++ {
		performRound(monkeys, func(i int) int { return i / 3 })
	}

	var topInspects [2]int
	for i := 0; i < len(monkeys); i++ {
		if monkeys[i].inspectCount > topInspects[0] {
			topInspects[0] = monkeys[i].inspectCount
			sort.Slice(topInspects[:], func(i, j int) bool { return topInspects[i] < topInspects[j] })
		}
	}
	fmt.Printf("Part 1 output: %d\n", topInspects[0]*topInspects[1])


	// part 2
	monkeys = parseMonkeys()

	// find common modulo
	// I hate these puzzles, if you didn't know about this modulo trick it would be crazy hard
	modulo := 1
	for _, monkey := range monkeys {
		modulo *= monkey.TestDivisor
	}

	for roundIndex := 0; roundIndex < 10000; roundIndex++ {
		performRound(monkeys, func(i int) int { return i % modulo })
	}

	topInspects[0] = 1
	topInspects[1] = 1
	for i := 0; i < len(monkeys); i++ {
		if monkeys[i].inspectCount > topInspects[0] {
			topInspects[0] = monkeys[i].inspectCount
			sort.Slice(topInspects[:], func(i, j int) bool { return topInspects[i] < topInspects[j] })
		}
	}
	fmt.Printf("Part 2 output: %d\n", topInspects[0]*topInspects[1])
}

func performRound(monkeys []*Monkey, worryHandler func(int) int ) {
	for monkeyID := 0; monkeyID < len(monkeys); monkeyID++ {
		thisMonkey := monkeys[monkeyID]
		for {
			if len(thisMonkey.Items) == 0 {
				break
			}

			thisMonkey.inspectCount++
			thisMonkey.Items[0] = thisMonkey.InspectOperation(thisMonkey.Items[0])
			thisMonkey.Items[0] = worryHandler(thisMonkey.Items[0])

			var targetMonkey *Monkey
			if thisMonkey.Items[0]%thisMonkey.TestDivisor == 0 {
				targetMonkey = monkeys[thisMonkey.TestTrueTarget]
			} else {
				targetMonkey = monkeys[thisMonkey.TestFalseTarget]
			}

			targetMonkey.Items = append(targetMonkey.Items, thisMonkey.Items[0])
			thisMonkey.Items = thisMonkey.Items[1:]
		}
	}
}

func parseMonkeys() []*Monkey {
	var monkeys []*Monkey

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		switch scanner.Text() {
		case "":
			continue
		default:
			var turn Monkey

			// parse monkey ID
			fmt.Sscanf(scanner.Text(), "Monkey %d:", &turn.MonkeyID)

			// parse starting items
			scanner.Scan()
			var startingItemsStr string
			startingItemsStr = strings.TrimPrefix(scanner.Text(), "  Starting items: ")
			startingItemsStrSlice := strings.Split(startingItemsStr, ",")
			for _, str := range startingItemsStrSlice {
				turn.Items = append(turn.Items, lib.MustParseInt(strings.TrimSpace(str)))
			}

			// parse operation
			scanner.Scan()
			var (
				operator string
				operand string
			)
			fmt.Sscanf(scanner.Text(), "  Operation: new = old %s %s", &operator, &operand)
			switch operand {
			case "old":
				turn.InspectOperation = func(old int) int { return old*old }
			default:
				operandInt := lib.MustParseInt(operand)
				switch operator {
				case "+":
					turn.InspectOperation = func(old int) int { return old+operandInt }
				case "*":
					turn.InspectOperation = func(old int) int { return old*operandInt }
				}
			}

			// parse test
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "  Test: divisible by %d", &turn.TestDivisor)
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "    If true: throw to monkey %d", &turn.TestTrueTarget)
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "    If false: throw to monkey %d", &turn.TestFalseTarget)

			monkeys = append(monkeys, &turn)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return monkeys
}