package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
)

type Block struct {
	bits   uint32
	height int
}

var blocks = []Block{
	{
		0b00000000000000000000000000011110,
		1,
	},
	{
		0b00000000000010000001110000001000,
		3,
	},
	{
		0b00000000000111000000010000000100,
		3,
	},
	{
		0b00010000000100000001000000010000,
		4,
	},
	{
		0b00000000000000000001100000011000,
		2,
	},
}

type Board []byte

//go:embed input
var input []byte

func main() {
	part1 := runSim(2022)
	fmt.Println(part1)
}

func runSim(rocksToFall int) int {
	board := Board{127, 127, 127, 127}

	// floor
	inputIndex := 0
	currentBlockIndex := 0
	rockCount := 0
	highestY := 4

	// spawn new block
	for {
		if rockCount == rocksToFall {
			return highestY - 4
		}

		block := blocks[currentBlockIndex]
		currentBlockIndex = (currentBlockIndex + 1) % len(blocks)

		blockY := highestY + 3 + block.height
		if blockY > len(board) {
			board = append(board, make([]byte, blockY-len(board))...)
		}

		for {
			// move sideways
			var bytes [4]byte
			binary.BigEndian.PutUint32(bytes[:], block.bits)

			switch input[inputIndex] {
			case '<':
				if bytes[0]&64 == 0 && bytes[1]&64 == 0 && bytes[2]&64 == 0 && bytes[3]&64 == 0 && (block.bits<<1)&binary.BigEndian.Uint32(board[blockY-4:blockY]) == 0 {
					block.bits <<= 1
				}
			case '>':
				if bytes[0]&1 == 0 && bytes[1]&1 == 0 && bytes[2]&1 == 0 && bytes[3]&1 == 0 && (block.bits>>1)&binary.BigEndian.Uint32(board[blockY-4:blockY]) == 0 {
					block.bits >>= 1
				}
			}
			inputIndex = (inputIndex + 1) % len(input)

			// move down
			if block.bits&binary.BigEndian.Uint32(board[blockY-5:blockY-1]) != 0 {
				binary.BigEndian.PutUint32(board[blockY-4:blockY], block.bits|binary.BigEndian.Uint32(board[blockY-4:blockY]))

				if blockY > highestY {
					highestY = blockY
				}
				rockCount++
				break
			}
			blockY--
		}
	}
}
