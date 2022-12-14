package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	x int
	y int
}

type Cavern struct {
	rocks      map[Position]bool
	sand       map[Position]bool
	maxDepth   int
	minRange   int
	maxRange   int
	floorDepth int
}

const SAND_ORIGIN_X = 500
const SAND_ORIGIN_Y = 0

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	fmt.Println(input)
	cavern := parseCavern(input)
	fmt.Printf("%+v\n", cavern)
	cavern.print()

	sandCount := -1
	settled := true
	for settled {
		sandCount++
		settled = cavern.produceSand()
	}

	cavern.print()
	fmt.Printf("total grains settled: %d\n", sandCount)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseCavern(input string) Cavern {
	lines := strings.Split(input, "\n")
	rocks := make(map[Position]bool)

	cavern := Cavern{
		maxDepth: 0,
		minRange: math.MaxInt,
		maxRange: math.MinInt,
		sand:     make(map[Position]bool),
	}

	for _, line := range lines {
		vertices := strings.Split(line, " -> ")
		lastVertex := parsePosition(vertices[0])
		rocks[lastVertex] = true
		cavern.updateBoundaries(lastVertex)
		for i := 1; i < len(vertices); i++ {
			currentVertex := parsePosition(vertices[i])
			cavern.updateBoundaries(currentVertex)
			for _, rock := range makeRockRun(lastVertex, currentVertex) {
				rocks[rock] = true
			}

			lastVertex = currentVertex
		}
	}

	cavern.rocks = rocks
	cavern.floorDepth = cavern.maxDepth + 2

	return cavern
}

func (c *Cavern) updateBoundaries(p Position) {
	if c.maxDepth < p.y {
		c.maxDepth = p.y
	}
	if c.minRange > p.x {
		c.minRange = p.x
	}
	if c.maxRange < p.x {
		c.maxRange = p.x
	}
}

func parsePosition(input string) Position {
	coords := strings.Split(input, ",")
	x, _ := strconv.Atoi(coords[0])
	y, _ := strconv.Atoi(coords[1])

	return Position{x: x, y: y}
}

func makeRockRun(from Position, to Position) []Position {
	delta := to.sub(from)
	var direction Position
	var count int

	if delta.y != 0 {
		// Vertical run
		if delta.y < 0 {
			// Up
			count = delta.y * -1
			direction.y = -1
		} else {
			// Down
			count = delta.y
			direction.y = 1
		}
	} else {
		// Horizontal run
		if delta.x < 0 {
			// Left
			count = delta.x * -1
			direction.x = -1
		} else {
			// Right
			count = delta.x
			direction.x = 1
		}
	}

	rocks := make([]Position, count)
	for i := 0; i < len(rocks); i++ {
		rocks[i] = from.add(direction.times(i + 1))
	}

	return rocks
}

func (p Position) add(o Position) Position {
	return Position{
		x: p.x + o.x,
		y: p.y + o.y,
	}
}

func (p Position) sub(o Position) Position {
	return Position{
		x: p.x - o.x,
		y: p.y - o.y,
	}
}

func (p Position) times(s int) Position {
	return Position{
		x: p.x * s,
		y: p.y * s,
	}
}

func (c *Cavern) hasRock(p Position) bool {
	if p.y == c.floorDepth {
		return true
	}

	_, ok := c.rocks[p]
	return ok
}

func (c *Cavern) hasSand(p Position) bool {
	_, ok := c.sand[p]
	return ok
}

func (c *Cavern) findNextObstacleDown(p Position) (Position, bool) {
	for y := 0; y <= c.maxDepth; y++ {
		nextPosition := p.add(Position{x: 0, y: y + 1})
		if c.hasRock(nextPosition) || c.hasSand(nextPosition) {
			return nextPosition, true
		}
	}

	// If no obstacles are found, return ok=false
	return p, false
}

func createSand() Position {
	return Position{x: SAND_ORIGIN_X, y: SAND_ORIGIN_Y}
}

// Create a sand particle and calculate where it settles. Returns true if the
// sand was able to settle, or false if the sand fell into the abyss (part 1)
// or was blocked at the source (part 2)
func (c *Cavern) produceSand() bool {
	sand := createSand()
	if c.hasRock(sand) || c.hasSand(sand) {
		return false
	}

	settled := false
	for !settled {

		obstacle, ok := c.findNextObstacleDown(sand)
		if !ok {
			return false
		}

		left := obstacle.add(Position{x: -1, y: 0})
		right := obstacle.add(Position{x: 1, y: 0})
		if !c.hasRock(left) && !c.hasSand(left) {
			// check left of the obstacle
			sand = left
		} else if !c.hasRock(right) && !c.hasSand(right) {
			// check right of the obstacle
			sand = right
		} else {
			settledAt := obstacle.add(Position{x: 0, y: -1})
			c.sand[settledAt] = true
			settled = true
		}
	}

	return true
}

func (c *Cavern) print() {
	// print headers. assume all headers are 3 digits at most
	depthAxisLength := len(fmt.Sprint(c.maxDepth)) + 1
	fmt.Println()
	for y := 0; y < 3; y++ {
		var line string
		for x := c.minRange - depthAxisLength; x < c.maxRange+1; x++ {
			switch x {
			case c.minRange:
				{
					line += string(fmt.Sprint(c.minRange)[y])
				}
			case c.maxRange:
				{
					line += string(fmt.Sprint(c.maxRange)[y])

				}
			case SAND_ORIGIN_X:
				{
					line += string(fmt.Sprint(SAND_ORIGIN_X)[y])
				}
			default:
				{
					line += " "
				}
			}
		}
		fmt.Println(line)
	}

	for y := 0; y < c.floorDepth+1; y++ {
		line := fmt.Sprintf("%d", y)
		currentDepthSize := len(fmt.Sprint(y))
		for i := 0; i < depthAxisLength-currentDepthSize; i++ {
			line += " "
		}

		for x := c.minRange; x < c.maxRange+1; x++ {
			pos := Position{x: x, y: y}
			if x == SAND_ORIGIN_X && y == SAND_ORIGIN_Y {
				line += "+"
			} else if c.hasRock(pos) {
				line += "#"
			} else if c.hasSand(pos) {
				line += "o"
			} else {
				line += "."
			}
		}
		fmt.Println(line)
	}
}
