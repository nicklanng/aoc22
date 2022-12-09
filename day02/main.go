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

// Hand represents a move in Rock, Paper, Scissors
type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

func byteToHand(b byte) Hand {
	switch b {
	case 'A', 'X':
		return Rock
	case 'B', 'Y':
		return Paper
	case 'C', 'Z':
		return Scissors
	default:
		panic("unknown hand")
	}
}

// Result represent an end game state of Rock, Paper, Scissors
type Result int

const (
	Loss Result = 0
	Draw Result = 3
	Win  Result = 6
)

func byteToResult(b byte) Result {
	switch b {
	case 'X':
		return Loss
	case 'Y':
		return Draw
	case 'Z':
		return Win
	default:
		panic("unknown result")
	}
}

func main() {
	scanner := bufio.NewScanner(bytes.NewReader(input))

	var (
		part1Score int
		part2Score int
	)

	for scanner.Scan() {
		s := scanner.Text()
		part1Score += calculatePart1Score(s[0], s[2])
		part2Score += calculatePart2Score(s[0], s[2])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1 Score: %d\n", part1Score)
	fmt.Printf("Part 2 Score: %d\n", part2Score)
}

func calculatePart1Score(firstPosition, secondPosition byte) int {
	opponent := byteToHand(firstPosition)
	me := byteToHand(secondPosition)
	result := playGame(opponent, me)
	return scoreGame(me, result)
}

func calculatePart2Score(firstPosition, secondPosition byte) int {
	opponent := byteToHand(firstPosition)
	result := byteToResult(secondPosition)
	me := handToThrowForResult(opponent, result)
	return scoreGame(me, result)
}

func playGame(opponent, me Hand) Result {
	if opponent == me {
		return Draw
	}

	// if this is extended later to rock paper scissors lizard spock
	// then we probably need a lookup table of results or something
	if opponent == Rock && me == Paper ||
		opponent == Paper && me == Scissors ||
		opponent == Scissors && me == Rock {
		return Win
	}

	return Loss
}

func handToThrowForResult(opponent Hand, result Result) Hand {
	switch result {
	case Loss:
		return (opponent - 1 + 3) % 3
	case Draw:
		return opponent
	case Win:
		return (opponent + 1) % 3
	default:
		panic("unknown result")
	}
}

func scoreGame(me Hand, result Result) int {
	return 1 + int(me) + int(result)
}
