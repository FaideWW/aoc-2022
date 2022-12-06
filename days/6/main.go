package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const MAX_LENGTH int = 14

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	fmt.Println(input)

	var marker int
	for start := 0; start < len(input)-MAX_LENGTH; start++ {
		substr := input[start : start+MAX_LENGTH]
		fmt.Printf("testing %s\n", substr)
		if allUnique(substr) {
			marker = start + MAX_LENGTH
			break
		}
	}

	fmt.Printf("first marker at %d\n", marker)

}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func allUnique(str string) bool {
	seen := make(map[rune]bool)
	for _, r := range str {
		if seen[r] {
			return false
		}
		seen[r] = true
	}

	return true
}
