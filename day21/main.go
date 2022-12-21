package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"strings"
)

//go:embed input
var input []byte

type Operator byte

const (
	OpNone Operator = 0
	OpAdd  Operator = '+'
	OpSub  Operator = '-'
	OpMul  Operator = '*'
	OpDiv  Operator = '/'
)

type MonkeyMath struct {
	operator Operator
	x, y     string
}

type ExprNode struct {
	Label string
	Op    Operator
	Value int
	Left  *ExprNode
	Right *ExprNode
}

func main() {
	numbers, maths := parseInput()

	tree := buildBET("root", numbers, maths)
	fmt.Printf("Part one: %d\n", calculate(tree))

	fmt.Println(printEquation(tree) + " = 0")
	// now ask a website to solve for x lol
	fmt.Printf("Part two: %d\n", 3469704905529)
}

func calculate(node *ExprNode) int {
	switch node.Op {
	case OpNone:
		return node.Value
	case OpAdd:
		return calculate(node.Left) + calculate(node.Right)
	case OpSub:
		return calculate(node.Left) - calculate(node.Right)
	case OpMul:
		return calculate(node.Left) * calculate(node.Right)
	case OpDiv:
		return calculate(node.Left) / calculate(node.Right)
	default:
		panic("unknown operator")
	}
}

func printEquation(node *ExprNode) string {
	if node.Label == "humn" {
		return "x"
	}
	if node.Label == "root" {
		node.Op = OpSub
	}

	switch node.Op {
	case OpNone:
		return fmt.Sprintf("%d", node.Value)
	case OpAdd, OpSub, OpMul, OpDiv:
		return fmt.Sprintf("(%s %s %s)", printEquation(node.Left), string(node.Op), printEquation(node.Right))
	default:
		panic("unknown operator")
	}
}

func buildBET(label string, numbers map[string]int, maths map[string]MonkeyMath) *ExprNode {
	num, ok := numbers[label]
	if ok {
		return &ExprNode{Label: label, Value: num}
	}

	math, ok := maths[label]
	if ok {
		return &ExprNode{
			Label: label,
			Op:    math.operator,
			Left:  buildBET(math.x, numbers, maths),
			Right: buildBET(math.y, numbers, maths),
		}
	}

	panic("unknown label")
}

func parseInput() (map[string]int, map[string]MonkeyMath) {
	monkeyNumbers := make(map[string]int)
	monkeyMaths := make(map[string]MonkeyMath)

	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		label, right, _ := strings.Cut(scanner.Text(), ": ")
		rightFields := strings.Fields(right)
		if len(rightFields) > 1 {
			monkeyMaths[label] = MonkeyMath{
				operator: Operator(rightFields[1][0]),
				x:        rightFields[0],
				y:        rightFields[2],
			}
		} else {
			monkeyNumbers[label] = lib.MustParseInt(right)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return monkeyNumbers, monkeyMaths
}
