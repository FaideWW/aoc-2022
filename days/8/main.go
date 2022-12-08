package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type TreeGrid [][]int

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	grid := parseInput(input)
	visibleTrees := countVisibleTrees(grid)
	fmt.Printf("there are %d visible trees\n", visibleTrees)

	topScore := findBestTree(grid)
	fmt.Printf("max scenic score: %d\n", topScore)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) TreeGrid {
	lines := strings.Split(input, "\n")
	grid := make([][]int, len(lines))

	for i, line := range lines {
		row := make([]int, len(line))
		for n, rune := range line {
			row[n] = int(rune - '0')
		}

		grid[i] = row
	}

	return grid
}

func countVisibleTrees(grid TreeGrid) (sum int) {
	gridHeight := len(grid)
	gridWidth := len(grid[0])

	// Outer trees are always visible
	sum += gridWidth*2 + (gridHeight-2)*2

	for y := 1; y < gridHeight-1; y++ {
		for x := 1; x < gridWidth-1; x++ {
			if grid.isTreeVisible(x, y) {
				sum++
			}
		}
	}

	return
}

func (g *TreeGrid) isTreeVisible(x int, y int) bool {
	return (g.isTreeVisibleFromLeft(x, y) ||
		g.isTreeVisibleFromRight(x, y) ||
		g.isTreeVisibleFromTop(x, y) ||
		g.isTreeVisibleFromBottom(x, y))
}

func (g *TreeGrid) isTreeVisibleFromLeft(x0 int, y0 int) bool {
	if x0 == 0 {
		return true
	}
	height := (*g)[y0][x0]
	for x := x0 - 1; x >= 0; x-- {
		if (*g)[y0][x] >= height {
			return false
		}
	}

	return true
}

func (g *TreeGrid) isTreeVisibleFromRight(x0 int, y0 int) bool {
	rowSize := len((*g)[y0])
	if x0 == rowSize-1 {
		return true
	}
	height := (*g)[y0][x0]
	for x := x0 + 1; x < rowSize; x++ {
		if (*g)[y0][x] >= height {
			return false
		}
	}

	return true
}

func (g *TreeGrid) isTreeVisibleFromTop(x0 int, y0 int) bool {
	if y0 == 0 {
		return true
	}
	height := (*g)[y0][x0]
	for y := y0 - 1; y >= 0; y-- {
		if (*g)[y][x0] >= height {
			return false
		}
	}

	return true
}

func (g *TreeGrid) isTreeVisibleFromBottom(x0 int, y0 int) bool {
	colSize := len(*g)
	if y0 == colSize-1 {
		return true
	}
	height := (*g)[y0][x0]
	for y := y0 + 1; y < colSize; y++ {
		if (*g)[y][x0] >= height {
			return false
		}
	}

	return true
}

func findBestTree(grid TreeGrid) (maxScore int) {
	gridHeight := len(grid)
	gridWidth := len(grid[0])

	maxScore = math.MinInt
	for y := 1; y < gridHeight-1; y++ {
		for x := 1; x < gridWidth-1; x++ {
			scenicScore := grid.getScenicScore(x, y)
			if scenicScore > maxScore {
				maxScore = scenicScore
			}
		}
	}

	return
}

func (g *TreeGrid) getScenicScore(x int, y int) int {
	return (g.getVisibilityLeft(x, y) *
		g.getVisibilityRight(x, y) *
		g.getVisibilityTop(x, y) *
		g.getVisibilityBottom(x, y))
}

func (g *TreeGrid) getVisibilityLeft(x0 int, y0 int) (visibility int) {
	if x0 == 0 {
		return 0
	}

	height := (*g)[y0][x0]
	for x := x0 - 1; x >= 0; x-- {
		visibility++
		if (*g)[y0][x] >= height {
			break
		}
	}

	return
}

func (g *TreeGrid) getVisibilityRight(x0 int, y0 int) (visibility int) {
	rowSize := len((*g)[y0])
	if x0 == rowSize-1 {
		return 0
	}

	height := (*g)[y0][x0]
	for x := x0 + 1; x < rowSize; x++ {
		visibility++
		if (*g)[y0][x] >= height {
			break
		}
	}

	return
}

func (g *TreeGrid) getVisibilityTop(x0 int, y0 int) (visibility int) {
	if y0 == 0 {
		return 0
	}

	height := (*g)[y0][x0]
	for y := y0 - 1; y >= 0; y-- {
		visibility++
		if (*g)[y][x0] >= height {
			break
		}
	}

	return
}

func (g *TreeGrid) getVisibilityBottom(x0 int, y0 int) (visibility int) {
	colSize := len(*g)
	if y0 == colSize-1 {
		return 0
	}

	height := (*g)[y0][x0]
	for y := y0 + 1; y < colSize; y++ {
		visibility++
		if (*g)[y][x0] >= height {
			break
		}
	}

	return
}
