package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	min int
	max int
}

type RangePair [2]Range

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])
	pairs := parseInput(input)
	fullOverlaps, partialOverlaps := countOverlaps(pairs)
	fmt.Printf("full overlaps found: %d\n", fullOverlaps)
	fmt.Printf("partial overlaps found: %d\n", partialOverlaps)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []RangePair {
	lines := strings.Split(input, "\n")
	pairs := make([]RangePair, len(lines)-1)

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		ranges := strings.Split(line, ",")
		range1Str := strings.Split(ranges[0], "-")
		range2Str := strings.Split(ranges[1], "-")

		range1Min, _ := strconv.Atoi(range1Str[0])
		range1Max, _ := strconv.Atoi(range1Str[1])
		range1 := Range{min: range1Min, max: range1Max}

		range2Min, _ := strconv.Atoi(range2Str[0])
		range2Max, _ := strconv.Atoi(range2Str[1])
		range2 := Range{min: range2Min, max: range2Max}

		pairs[i] = RangePair{range1, range2}
	}
	return pairs
}

func countOverlaps(pairs []RangePair) (fullOverlaps int, partialOverlaps int) {
	for _, pair := range pairs {
		if doRangesFullyOverlap(pair[0], pair[1]) {
			fullOverlaps += 1
			partialOverlaps += 1
		} else if doRangesOverlapAtAll(pair[0], pair[1]) {
			partialOverlaps += 1
		}
	}

	return
}

func doRangesFullyOverlap(range1 Range, range2 Range) bool {
	if range1.min <= range2.min && range1.max >= range2.max ||
		range2.min <= range1.min && range2.max >= range1.max {
		return true
	}
	return false

}

func doRangesOverlapAtAll(range1 Range, range2 Range) bool {
	if range1.min <= range2.min && range1.max >= range2.min || range1.min <= range2.max && range1.max >= range2.max {
		return true
	}
	return false
}
