package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Vec3 struct {
	x int
	y int
	z int
}

type Grid struct {
	size int
	grid [][][]bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	// fmt.Println(input)

	grid := newGrid(100)
	faces := parseInput(input, &grid)

	fmt.Printf("exposed faces: %d\n", faces)

	exteriorFaces := findExteriorSurface(&grid)

	fmt.Printf("exterior faces: %d\n", exteriorFaces)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func newGrid(size int) Grid {
	grid := make([][][]bool, size)
	for i := 0; i < size; i++ {
		grid[i] = make([][]bool, size)
		for j := 0; j < size; j++ {
			grid[i][j] = make([]bool, size)
		}
	}

	return Grid{
		size: size,
		grid: grid,
	}
}

func parseInput(input string, grid *Grid) int {
	exposedFaces := 0
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		coords := strings.Split(line, ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])

		// insert cube into grid
		grid.grid[x][y][z] = true
		exposedFaces += 6

		// check for neighbors
		if x > 0 && grid.grid[x-1][y][z] {
			exposedFaces -= 2
		}
		if x < grid.size && grid.grid[x+1][y][z] {
			exposedFaces -= 2
		}
		if y > 0 && grid.grid[x][y-1][z] {
			exposedFaces -= 2
		}
		if y < grid.size && grid.grid[x][y+1][z] {
			exposedFaces -= 2
		}
		if z > 0 && grid.grid[x][y][z-1] {
			exposedFaces -= 2
		}
		if z < grid.size && grid.grid[x][y][z+1] {
			exposedFaces -= 2
		}
	}
	return exposedFaces
}

func findExteriorSurface(grid *Grid) int {
	// the general idea: starting at a known outside air tile (say, 0,0,0), we
	// can implicitly find all exterior faces by flood-filling from the air tile
	// and adding a face any time flood-fill would move into a tile in the volume

	startingTile := Vec3{0, 0, 0}
	if grid.check(startingTile) {
		panic(errors.New("starting tile is in the volume; try another tile"))
	}

	getNeighbors := func(v Vec3) ([]Vec3, []Vec3) {
		candidates := []Vec3{
			{v.x - 1, v.y, v.z},
			{v.x + 1, v.y, v.z},
			{v.x, v.y - 1, v.z},
			{v.x, v.y + 1, v.z},
			{v.x, v.y, v.z - 1},
			{v.x, v.y, v.z + 1},
		}

		airNeighbors := make([]Vec3, 0)
		volumeNeighbors := make([]Vec3, 0)

		for _, v := range candidates {
			if v.x < -1 || v.x > grid.size || v.y < -1 || v.y > grid.size || v.z < -1 || v.z > grid.size {
				continue
			}
			if grid.check(v) {
				volumeNeighbors = append(volumeNeighbors, v)
			} else {
				airNeighbors = append(airNeighbors, v)
			}
		}

		return airNeighbors, volumeNeighbors
	}

	frontier := make([]Vec3, 1)
	frontier[0] = startingTile

	seen := make(map[Vec3]bool)

	exteriorFaces := 0

	// naive flood fill; if it's too slow we can look into span-fill
	for len(frontier) > 0 {
		currentIndex := len(frontier) - 1
		current := frontier[currentIndex]
		frontier = frontier[:currentIndex]

		neighbors, faces := getNeighbors(current)

		exteriorFaces += len(faces)

		for _, next := range neighbors {
			if !seen[next] {
				seen[next] = true
				frontier = append(frontier, next)
			}
		}
	}

	return exteriorFaces
}

func (g *Grid) check(v Vec3) bool {
	if v.x < 0 || v.x >= g.size || v.y < 0 || v.y >= g.size || v.z < 0 || v.z >= g.size {
		return false
	}
	return g.grid[v.x][v.y][v.z]
}
