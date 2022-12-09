package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ROPE_LENGTH = 10

type Move string

const (
	Right Move = "R"
	Left  Move = "L"
	Up    Move = "U"
	Down  Move = "D"
)

type Position struct {
	x int
	y int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	moves := parseInput(input)
	tilesSeen := executeMoves(moves)
	fmt.Printf("unique tiles seen: %d\n", tilesSeen)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []Move {
	lines := strings.Split(input, "\n")
	moves := make([]Move, 0)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		amt, _ := strconv.Atoi(parts[1])

		for j := 0; j < amt; j++ {
			var move Move
			switch parts[0] {
			case "R":
				{
					move = Right
				}
			case "L":
				{
					move = Left
				}
			case "U":
				{
					move = Up
				}
			case "D":
				{
					move = Down
				}
			default:
				{
					panic(errors.New("invalid move"))
				}
			}
			moves = append(moves, move)
		}
	}
	return moves
}

func executeMoves(moves []Move) int {
	start := Position{x: 0, y: 0}
	ropeNodes := make([]Position, ROPE_LENGTH)
	for i := 0; i < len(ropeNodes); i++ {
		ropeNodes[i] = start
	}

	tailSeenTiles := make(map[Position]bool, 0)
	tailSeenTiles[start] = true
	uniqueTiles := 1

	for _, move := range moves {
		ropeNodes[0] = doMove(ropeNodes[0], move)
		for i := 1; i < len(ropeNodes); i++ {
			shouldMove, newPos := shouldTailMove(ropeNodes[i-1], ropeNodes[i])
			if shouldMove {
				ropeNodes[i] = newPos
			}
		}
		tailPos := ropeNodes[len(ropeNodes)-1]
		if !tailSeenTiles[tailPos] {
			tailSeenTiles[tailPos] = true
			uniqueTiles++
		}
	}

	return uniqueTiles
}

func doMove(pos Position, move Move) Position {
	switch move {
	case Right:
		{
			return Position{x: pos.x + 1, y: pos.y}
		}
	case Left:
		{
			return Position{x: pos.x - 1, y: pos.y}
		}
	case Up:
		{
			return Position{x: pos.x, y: pos.y - 1}
		}
	case Down:
		{
			return Position{x: pos.x, y: pos.y + 1}
		}
	default:
		{
			panic(errors.New("invalid move"))
		}
	}
}

func shouldTailMove(head Position, tail Position) (bool, Position) {
	// First check if the tail and head overlap
	if head == tail {
		return false, tail
	}

	direction := Position{x: head.x - tail.x, y: head.y - tail.y}
	shouldMove := false
	if direction.y > 1 {
		direction.y = 1
		shouldMove = true
	}
	if direction.y < -1 {
		direction.y = -1
		shouldMove = true
	}
	if direction.x > 1 {
		direction.x = 1
		shouldMove = true
	}
	if direction.x < -1 {
		direction.x = -1
		shouldMove = true
	}

	nextTail := tail
	if shouldMove {
		nextTail = Position{x: tail.x + direction.x, y: tail.y + direction.y}
	}

	return shouldMove, nextTail
}
