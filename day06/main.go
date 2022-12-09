package main

import (
	_ "embed"
	"fmt"
)

const (
	headerMarkerLen  = 4
	messageMarkerLen = 14
)

//go:embed input
var input []byte

func main() {
	fmt.Printf("First packet header: %d\n", findFirstMarker(headerMarkerLen))
	fmt.Printf("First message header: %d\n", findFirstMarker(messageMarkerLen))
}

func findFirstMarker(length int) int {
	for i := length; i <= len(input); i++ {
		if uniqueBytes(input[i-length : i]) {
			return i
		}
	}
	return -1
}

func uniqueBytes(s []byte) bool {
	m := make([]bool, 256)
	for j := 0; j < len(s); j++ {
		if m[s[j]] {
			return false
		}
		m[s[j]] = true
	}
	return true
}
