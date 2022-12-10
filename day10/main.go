package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"github.com/nicklanng/aoc22/lib"
	"github.com/nicklanng/aoc22/lib/vdu"
	"strings"
)

//go:embed input
var input []byte

func main() {
	memory := vdu.NewMemory(4096)
	cpu := vdu.NewCPU(memory)
	crt := vdu.NewCRT(cpu, 40, 6)

	// load program into memory
	program := parseProgram()
	memory.LoadProgram(program)

	// execute program
	var signalStrength int
	for  {
		x := cpu.XRegister()

		crt.ExecuteCycle()
		if !cpu.ExecuteCycle() {
			break
		}

		if cpu.Cycle() % 40 == 20 {
			signalStrength += cpu.Cycle() * x
		}
	}

	// output
	fmt.Printf("Part 1 signal strength: %d\n", signalStrength)
	fmt.Println("Part 2 render")
	crt.Render()
}

func parseProgram() []byte {
	var program []byte
	scanner := bufio.NewScanner(bytes.NewReader(input))

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		switch fields[0] {
		case "noop":
			program = append(program, vdu.InstNoop)
		case "addx":
			program = append(program, vdu.InstAddX)
			program = append(program, byte(lib.MustParseInt(fields[1])))
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return program
}
