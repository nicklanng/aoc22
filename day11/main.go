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
	TestFunc func(int) int
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

	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].inspectCount > monkeys[j].inspectCount })
	fmt.Printf("Part 1 output: %d\n", monkeys[0].inspectCount*monkeys[1].inspectCount)


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

	sort.Slice(monkeys, func(i, j int) bool { return monkeys[i].inspectCount > monkeys[j].inspectCount })
	fmt.Printf("Part 2 output: %d\n", monkeys[0].inspectCount*monkeys[1].inspectCount)
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

			var targetMonkey = monkeys[thisMonkey.TestFunc(thisMonkey.Items[0])]
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
			var monkey Monkey

			// parse monkey ID
			fmt.Sscanf(scanner.Text(), "Monkey %d:", &monkey.MonkeyID)

			// parse starting items
			scanner.Scan()
			var startingItemsStr string
			startingItemsStr = strings.TrimPrefix(scanner.Text(), "  Starting items: ")
			startingItemsStrSlice := strings.Split(startingItemsStr, ",")
			for _, str := range startingItemsStrSlice {
				monkey.Items = append(monkey.Items, lib.MustParseInt(strings.TrimSpace(str)))
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
				monkey.InspectOperation = func(old int) int { return old*old }
			default:
				operandInt := lib.MustParseInt(operand)
				switch operator {
				case "+":
					monkey.InspectOperation = func(old int) int { return old+operandInt }
				case "*":
					monkey.InspectOperation = func(old int) int { return old*operandInt }
				}
			}

			// parse test
			var divisor, trueTarget, falseTarget int
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "  Test: divisible by %d", &divisor)
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "    If true: throw to monkey %d", &trueTarget)
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "    If false: throw to monkey %d", &falseTarget)

			monkey.TestDivisor = divisor
			monkey.TestFunc = func(i int) int {
				if i%monkey.TestDivisor == 0 {
					return trueTarget
				}
				return falseTarget
			}

			monkeys = append(monkeys, &monkey)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return monkeys
}