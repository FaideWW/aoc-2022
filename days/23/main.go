package main

import (
	//	"errors"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
)

type Direction int

const (
	NORTH Direction = 0
	SOUTH Direction = 1
	WEST  Direction = 2
	EAST  Direction = 3
)

type Position struct {
	x int
	y int
}

type Grove struct {
	elves map[Position]bool
	min   Position
	max   Position
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])
	grove := parseInput(input)
	turnsCompleted := grove.run(math.MaxInt)
	fmt.Printf("completed %d rounds\n", turnsCompleted+1)
	emptyTiles := grove.computeEmptyTiles()
	fmt.Printf("empty tiles: %d\n", emptyTiles)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) Grove {
	lines := strings.Split(input, "\n")
	elves := make(map[Position]bool)

	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			if line[x] == '#' {
				p := Position{x, y}
				elves[p] = true
			}
		}
	}

	g := Grove{
		elves: elves,
	}
	g.recomputeBoundingBox()

	return g
}

type Move struct {
	to   Position
	from Position
}

func (g *Grove) run(turns int) int {
	i := 0
	for ; i < turns; i++ {
		seenDestinations := make(map[Position][]Position)
		for elfPos := range g.elves {
			if g.isAlone(elfPos) {
				continue
			}
			for j := 0; j < 4; j++ {
				direction := (i + j) % 4
				candidates := elfPos.getCandidates(Direction(direction))
				occupied := false
				for _, candidate := range candidates {
					if g.elves[candidate] {
						occupied = true
					}
				}

				if !occupied {
					destination := candidates[1]
					_, ok := seenDestinations[destination]
					if !ok {
						seenDestinations[destination] = make([]Position, 1)
						seenDestinations[destination][0] = elfPos
					} else {
						seenDestinations[destination] = append(seenDestinations[destination], elfPos)
					}
					break
				}
			}
		}

		// If no one tried to move, we're done
		if len(seenDestinations) == 0 {
			break
		}

		for dest, origin := range seenDestinations {
			if len(origin) > 1 {
				// If more than one elf tried to move to this position, no one moves
				continue
			}

			pos := origin[0]
			g.elves[dest] = true
			delete(g.elves, pos)

			g.recomputeBoundingBox()
		}
	}
	return i
}

func (g *Grove) recomputeBoundingBox() {

	min := Position{math.MaxInt, math.MaxInt}
	max := Position{math.MinInt, math.MinInt}

	for elf := range g.elves {
		if min.x > elf.x {
			min.x = elf.x
		}
		if min.y > elf.y {
			min.y = elf.y
		}
		if max.x < elf.x {
			max.x = elf.x
		}
		if max.y < elf.y {
			max.y = elf.y
		}
	}

	g.min = min
	g.max = max
}

func (g *Grove) computeEmptyTiles() int {
	emptyTiles := (g.max.x + 1 - g.min.x) * (g.max.y + 1 - g.min.y)
	elfCount := len(g.elves)
	return emptyTiles - elfCount
}

func (g *Grove) isAlone(p Position) bool {
	return (!g.elves[Position{p.x - 1, p.y - 1}] &&
		!g.elves[Position{p.x, p.y - 1}] &&
		!g.elves[Position{p.x + 1, p.y - 1}] &&
		!g.elves[Position{p.x - 1, p.y}] &&
		!g.elves[Position{p.x + 1, p.y}] &&
		!g.elves[Position{p.x - 1, p.y + 1}] &&
		!g.elves[Position{p.x, p.y + 1}] &&
		!g.elves[Position{p.x + 1, p.y + 1}])
}

func (p *Position) getCandidates(direction Direction) []Position {
	switch direction {
	case NORTH:
		{
			return []Position{
				{p.x - 1, p.y - 1},
				{p.x, p.y - 1},
				{p.x + 1, p.y - 1},
			}
		}
	case SOUTH:
		{
			return []Position{
				{p.x - 1, p.y + 1},
				{p.x, p.y + 1},
				{p.x + 1, p.y + 1},
			}
		}
	case WEST:
		{
			return []Position{
				{p.x - 1, p.y - 1},
				{p.x - 1, p.y},
				{p.x - 1, p.y + 1},
			}
		}
	case EAST:
		{
			return []Position{
				{p.x + 1, p.y - 1},
				{p.x + 1, p.y},
				{p.x + 1, p.y + 1},
			}
		}
	}
	panic(errors.New("unrecognized direction"))
}

func (g *Grove) print() string {
	var output string
	for y := g.min.y; y < g.max.y+1; y++ {
		line := ""
		for x := g.min.x; x < g.max.x+1; x++ {
			if g.elves[Position{x, y}] {
				line += "#"
			} else {
				line += "."
			}
		}
		output += line + "\n"
	}

	return output
}
