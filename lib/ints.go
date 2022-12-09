package lib

import "strconv"

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