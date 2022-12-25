package main

import (
	"fmt"
	"os"
	"strings"
)

var snafuTable = map[byte]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

var intTable = map[int]byte{
	-2: '=',
	-1: '-',
	0:  '0',
	1:  '1',
	2:  '2',
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))

	snafus := parseInput(input)

	sum := 0
	for _, s := range snafus {
		sum += snafuToInt(s)
	}

	totalSnafu := intToSnafu(sum)
	fmt.Printf("sum snafu: %s\n", totalSnafu)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []string {
	return strings.Split(input, "\n")
}

func snafuToInt(s string) int {
	value := 0
	for i := 0; i < len(s); i++ {
		digitValue := snafuTable[s[len(s)-1-i]]
		value += digitValue * intPow(5, i)
	}
	return value
}

func intToSnafu(i int) string {
	remainingValue := i
	snafu := ""
	digits := make([]int, 0)
	// currentDigit := 0
	for remainingValue != 0 {
		carry := 0
		digitAtPlace := remainingValue % 5

		if digitAtPlace > 2 {
			carry = 5
			digitAtPlace -= carry
		}
		digits = append(digits, digitAtPlace)
		nextValue := ((remainingValue + carry) / 5)
		remainingValue = nextValue
	}

	for i := 0; i < len(digits); i++ {
		snafu += string(intTable[digits[len(digits)-1-i]])

	}

	return snafu
}

func intPow(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := base
	for i := 2; i <= exp; i++ {
		result *= base
	}
	return result
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
