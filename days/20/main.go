package main

import (
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const DECRYPTION_KEY = 811589153
const MIX_COUNT = 10

// const DECRYPTION_KEY = 811589153
// const MIX_COUNT = 10

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))

	data := parseInput(input)
	decrypt(data, DECRYPTION_KEY, MIX_COUNT)

	coord1, coord2, coord3 := computeCoordinates(data)

	sum := coord1 + coord2 + coord3
	fmt.Printf("coordinate sum: %d\n", sum)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) *list.List {
	lines := strings.Split(input, "\n")
	l := list.New()
	for _, line := range lines {
		value, _ := strconv.Atoi(line)
		l.PushBack(value)
	}

	return l
}

func getNodePosition(l *list.List, e *list.Element) int {
	n := l.Front()
	for i := 0; i < l.Len(); i++ {
		if e == n {
			return i
		}
		n = n.Next()
	}

	return -1
}

func getNodeAt(l *list.List, idx int) *list.Element {
	n := l.Front()
	for i := 0; i < idx; i++ {
		n = n.Next()
	}

	return n
}

func decrypt(data *list.List, key int, mixCount int) *list.List {
	// Since the order of the actual list will be shuffled around as we move
	// nodes, we want to remember the original order of these nodes so that we
	// can iterate through them effectively. We can also apply the decryption key
	// here
	orderedNodes := make([]*list.Element, 0)
	for n := data.Front(); n != nil; n = n.Next() {
		orderedNodes = append(orderedNodes, n)
		n.Value = n.Value.(int) * key
	}

	for i := 0; i < mixCount; i++ {
		for j, node := range orderedNodes {
			move := node.Value.(int)
			if move > 0 {
				// remove the node from the list before rotation
				initialPosition := getNodePosition(data, node)
				data.Remove(node)

				// calculate the new position in the list
				newPosition := (initialPosition + node.Value.(int)) % data.Len()
				newNext := getNodeAt(data, newPosition)
				// since we removed the original node from the list, we need to update
				// the reference for the next round
				orderedNodes[j] = data.InsertBefore(node.Value, newNext)
			} else {
				// remove the node from the list before rotation
				initialPosition := getNodePosition(data, node)
				data.Remove(node)

				// calculate the new position in the list
				newPosition := ((initialPosition + node.Value.(int)) % data.Len())
				if newPosition <= 0 {
					newPosition += data.Len()
				}
				// we're finding the node before our new position, since we're
				// inserting after it
				newPrev := getNodeAt(data, newPosition-1)
				orderedNodes[j] = data.InsertAfter(node.Value, newPrev)
			}
		}
	}
	return data
}

const COORD1_LOCATION = 1000
const COORD2_LOCATION = 2000
const COORD3_LOCATION = 3000

func computeCoordinates(l *list.List) (int, int, int) {
	var coord1, coord2, coord3 int

	node := l.Front()
	// scan to 0
	for node.Value != 0 {
		node = node.Next()
	}

	for i := 0; i < COORD3_LOCATION+1; i++ {
		switch i {
		case COORD1_LOCATION:
			{
				coord1 = node.Value.(int)
				fmt.Printf("coord1: %d\n", coord1)
			}
		case COORD2_LOCATION:
			{
				coord2 = node.Value.(int)
				fmt.Printf("coord2: %d\n", coord2)
			}
		case COORD3_LOCATION:
			{
				coord3 = node.Value.(int)
				fmt.Printf("coord3: %d\n", coord3)
			}
		default:
			{
			}
		}

		if node.Next() == nil {
			node = l.Front()
		} else {
			node = node.Next()
		}
	}
	return coord1, coord2, coord3
}

func printList(l *list.List) string {
	output := ""
	for e := l.Front(); e != nil; e = e.Next() {
		output += fmt.Sprintf("%d, ", e.Value)
	}
	return output[:len(output)-2]
}
