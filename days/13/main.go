package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Integer int
type List []Datum

type Datum interface {
	compare(d Datum) int
}

type Packet Datum

type PacketPair struct {
	left  Packet
	right Packet
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	packetPairs := parseInput(input)

	sumIndices := 0
	for i, pair := range packetPairs {
		if pair.isInOrder() {
			sumIndices += i + 1
		}
	}

	fmt.Printf("ordered pairs sum: %d\n", sumIndices)

	packets := flattenPairs(packetPairs)

	dividers := []Packet{parsePacket("[[2]]"), parsePacket("[[6]]")}

	packets = append(packets, dividers...)

	sort.Slice(packets, func(i, j int) bool {
		return packets[i].compare(packets[j]) < 0
	})

	divider1Index := findIndex(packets, dividers[0]) + 1
	divider2Index := findIndex(packets, dividers[1]) + 1

	fmt.Printf("decoder key: %d\n", divider1Index*divider2Index)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []PacketPair {
	inputPairs := strings.Split(input, "\n\n")
	pairs := make([]PacketPair, len(inputPairs))

	for i, inputPair := range inputPairs {
		inputPackets := strings.Split(inputPair, "\n")
		pairs[i] = PacketPair{
			left:  parsePacket(inputPackets[0]),
			right: parsePacket(inputPackets[1]),
		}
	}

	return pairs
}

func parsePacket(input string) Packet {
	listStack := make([]List, 1)
	rootList := make(List, 0)
	listStack[0] = rootList
	var currentList *List

	for i := 1; i < len(input)-1; {
		currentList = &listStack[len(listStack)-1]
		token := input[i]
		switch token {
		case ',':
			{
				// ignore, continue
				i++
			}
		case '[':
			{
				// start a new list
				newList := make(List, 0)
				listStack = append(listStack, newList)
				i++
			}
		case ']':
			{
				// end the current list
				if len(listStack) == 0 {
					panic(errors.New("stack is empty"))
				}

				topIndex := len(listStack) - 1
				parentIndex := len(listStack) - 2
				listStack[parentIndex] = append(listStack[parentIndex], listStack[topIndex])
				listStack = listStack[:topIndex]
				i++
			}
		default:
			{
				// consume until a comma or right paren, then push into the current list
				j := i + 1
				for ; input[j] != ',' && input[j] != ']'; j++ {
				}
				token, err := strconv.Atoi(input[i:j])
				if err != nil {
					panic(err)
				}
				*currentList = append(*currentList, Integer(token))
				i = j
			}
		}
	}

	return listStack[0]
}

func (i Integer) compare(d Datum) int {
	// if d is an int, compare the two ints
	if dInt, ok := d.(Integer); ok {
		return (int)(i - dInt)
	}

	// if d is a list, up-convert i to a list
	iList := List{i}
	return iList.compare(d)
}

func (l List) compare(d Datum) int {
	var dList List
	// if d is an int, up-convert d to a list
	if dInt, ok := d.(Integer); ok {
		dList = List{dInt}
	} else {
		// if d is a list, compare the two lists
		dList = d.(List)
	}

	var minLen int
	if len(l) < len(dList) {
		minLen = len(l)
	} else {
		minLen = len(dList)
	}

	for i := 0; i < minLen; i++ {
		result := l[i].compare(dList[i])
		if result != 0 {
			return result
		}
	}
	return len(l) - len(dList)
}

func (p PacketPair) isInOrder() bool {
	return p.left.compare(p.right) <= 0
}

func flattenPairs(pairs []PacketPair) []Packet {
	packets := make([]Packet, len(pairs)*2)
	for i, pair := range pairs {
		packets[2*i] = pair.left
		packets[2*i+1] = pair.right
	}
	return packets
}

func findIndex(packets []Packet, toFind Packet) int {
	for i, toCompare := range packets {
		if toCompare.compare(toFind) == 0 {
			return i
		}
	}
	return -1
}
