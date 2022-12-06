package main

import (
	_ "embed"
	"fmt"
	"sort"
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
	packet := make([]byte, len(s))
	copy(packet, s)

	sort.Slice(packet, func(i, j int) bool { return packet[i] < packet[j] })

	for j := 1; j < len(packet); j++ {
		if packet[j-1] == packet[j] {
			return false
		}
	}
	return true
}
