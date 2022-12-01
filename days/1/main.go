package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])
	counts := readCalories(input)
	max := maxInSlice(counts)
	fmt.Println("max", max)
	top3 := topThreeInSlice(counts)
	fmt.Println("sum of top 3", top3)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func readCalories(input string) []int {
	lines := strings.Split(input, "\n")

	counts := make([]int, len(lines))

	currentTotal := 0
	currentIndex := 0
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			counts[currentIndex] = currentTotal
			currentIndex++
			currentTotal = 0
		} else {
			lineNum, _ := strconv.Atoi(line)
			currentTotal += lineNum
		}
	}

	return counts
}

func maxInSlice(arr []int) int {
	max := math.MinInt
	for _, num := range arr {
		if num > max {
			max = num
		}
	}

	return max
}

func topThreeInSlice(arr []int) int {
	sort.Ints(arr)

	sum := 0
	for _, n := range arr[len(arr)-3:] {
		sum += n
	}

	return sum
}
