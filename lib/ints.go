package lib

import (
	"math"
	"strconv"
)

func MustParseInt(s string) int {
	number, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return number
}

func IntAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IntMax(x ...int) int {
	max := math.MinInt
	for i := range x {
		if x[i] > max {
			max = x[i]
		}
	}
	return max
}
