package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"math"
)

//go:embed input
var input []byte

func main() {
	str := parseInput()

	// part 1
	var part1 int
	for _, s := range str {
		part1 += parseSNAFU(s)
	}
	fmt.Println(intToSNAFU(part1))
}

func parseSNAFU(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		result += parseSNAFUChar(s[len(s)-1-i]) * int(math.Pow(5, float64(i)))
	}
	return result
}

func parseSNAFUChar(b byte) int {
	switch b {
	case '2':
		return 2
	case '1':
		return 1
	case '0':
		return 0
	case '-':
		return -1
	case '=':
		return -2
	default:
		panic("invalid number")
	}
}

func intToSNAFU(a int) string {
	var result string

	for a != 0 {
		switch a % 5 {
		case 0:
			result = "0" + result
		case 1:
			a -= 1
			result = "1" + result
		case 2:
			a -= 2
			result = "2" + result
		case 3:
			a += 2
			result = "=" + result
		case 4:
			a += 1
			result = "-" + result
		default:
			panic("hell naw")
		}
		a /= 5
	}
	return result
}

func parseInput() []string {
	var str []string
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		str = append(str, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return str
}
