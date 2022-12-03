package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Rucksack struct {
	sackContents string
	duplicates   []rune
}

type RucksackGroup struct {
	sacks     [3]Rucksack
	badgeItem rune
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])

	rucksacks := parseInput(input)
	priority := getTotalPriority(rucksacks)
	fmt.Println("first priority", priority)

	groups := groupRucksacks(rucksacks)
	groupPriority := getGroupPriority(groups)
	fmt.Println("group priority", groupPriority)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []Rucksack {
	lines := strings.Split(input, "\n")
	sacks := make([]Rucksack, len(lines))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		c1, c2 := splitString(line)
		sack := Rucksack{
			sackContents: strings.TrimSpace(line),
			duplicates:   findDuplicates(c1, c2),
		}
		sacks[i] = sack
	}

	return sacks
}

func splitString(input string) (string, string) {
	firstHalf := input[:len(input)/2]
	secondHalf := input[len(input)/2:]

	return firstHalf, secondHalf
}

func findDuplicates(str1 string, str2 string) []rune {
	hash := make(map[rune]bool)

	duplicates := make([]rune, 0)

	for _, r := range str1 {
		hash[r] = true
	}

	for _, r := range str2 {
		if hash[r] == true {
			duplicates = append(duplicates, r)
		}
	}

	// There's a chance that multiple duplicates of the same rune were found in
	// str2, but we don't want to double-count them so we de-dupe the array
	// before returning.
	return dedupe(duplicates)
}

func dedupe(array []rune) (uniques []rune) {
	seen := make(map[rune]bool)

	for _, r := range array {
		if !seen[r] {
			uniques = append(uniques, r)
			seen[r] = true
		}
	}

	return
}

func getTotalPriority(sacks []Rucksack) (sum int) {
	for _, sack := range sacks {
		sackSum := 0
		for _, dupe := range sack.duplicates {
			sackSum += getRunePriority(dupe)
		}
		sum += sackSum
	}

	return
}

func getRunePriority(r rune) (priority int) {
	// ASCII char codes enumerate through A-Z before a-z, but the priorities are
	// enumerated the other way around (a-z then A-Z). To account for this we
	// look for any char codes that start before 'a', and add 58 to them (2*26
	// alphabet characters + 6 special characters [\]^_`)
	priority = int(r - 'a' + 1)
	if priority < 1 {
		priority += 26*2 + 6
	}
	return
}

func groupRucksacks(sacks []Rucksack) []RucksackGroup {
	groups := make([]RucksackGroup, len(sacks)/3)
	for i := 0; i < len(sacks)/3; i++ {
		sack1 := sacks[i*3]
		sack2 := sacks[i*3+1]
		sack3 := sacks[i*3+2]
		groups[i] = RucksackGroup{
			sacks:     [3]Rucksack{sack1, sack2, sack3},
			badgeItem: findBadgeItem(sack1, sack2, sack3),
		}
	}

	return groups
}

func findBadgeItem(sack1 Rucksack, sack2 Rucksack, sack3 Rucksack) rune {
	firstDuplicates := findDuplicates(sack1.sackContents, sack2.sackContents)

	secondDuplicates := findDuplicates(string(firstDuplicates), sack3.sackContents)

	if len(secondDuplicates) != 1 {
		panic(errors.New("Badge item count was not exactly 1"))
	}
	return secondDuplicates[0]
}

func getGroupPriority(groups []RucksackGroup) (sum int) {
	for _, group := range groups {
		sum += getRunePriority(group.badgeItem)
	}

	return
}
