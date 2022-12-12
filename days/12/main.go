package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"strings"
)

type Position struct {
	x int
	y int
}

type Grid struct {
	start    Position
	end      Position
	tiles    [][]int
	lowTiles []Position
	costs    map[Position]int
}

type PQItem struct {
	value    Position
	priority int
	index    int
}

type PriorityQueue []*PQItem

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	grid := parseInput(input)

	minPathLength := math.MaxInt

	// we can brute force this by just running pathfinding for every starting node and caching the costs across runs
	for _, start := range grid.lowTiles {
		path, ok := grid.findShortestPath(start, grid.end)
		if !ok {
			continue
		}
		pathLength := len(path)
		if pathLength < minPathLength {
			fmt.Printf("shorter path found starting at %+v (length: %d)\n", start, pathLength)
			minPathLength = pathLength
		}
	}

	fmt.Printf("shortest path is %d steps\n", minPathLength)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) Grid {
	lines := strings.Split(input, "\n")
	rows := make([][]int, len(lines))

	var start Position
	var end Position
	lowTiles := make([]Position, 0)

	for y, line := range lines {
		row := make([]int, len(line))
		for x, tile := range line {
			pos := Position{x: x, y: y}
			if tile == 'S' {
				start = pos
			}
			if tile == 'E' {
				end = pos
			}
			row[x] = runeToHeight(tile)
			if row[x] == 0 {
				lowTiles = append(lowTiles, pos)
			}
		}

		rows[y] = row
	}

	return Grid{
		start:    start,
		end:      end,
		tiles:    rows,
		lowTiles: lowTiles,
		costs:    make(map[Position]int),
	}
}

func runeToHeight(r rune) int {
	if r == 'S' {
		return 0
	}
	if r == 'E' {
		return 25
	}
	return int(r - 'a')
}

func (g *Grid) findShortestPath(start Position, end Position) ([]Position, bool) {

	frontier := make(PriorityQueue, 1)
	frontier[0] = &PQItem{
		value:    start,
		priority: 0,
		index:    0,
	}

	heap.Init(&frontier)

	g.costs[start] = 0
	cameFrom := make(map[Position]Position)
	cameFrom[start] = Position{x: -1, y: -1}

	for frontier.Len() > 0 {
		current := heap.Pop(&frontier).(*PQItem)
		if current.value == end {
			break
		}

		for _, next := range g.getNeighbors(current.value) {
			newCost := g.costs[current.value] + 1
			if foundCost, found := g.costs[next]; !found || newCost < foundCost {
				g.costs[next] = newCost
				heap.Push(&frontier, &PQItem{value: next, priority: newCost})
				cameFrom[next] = current.value
			}
		}
	}

	path := make([]Position, 0)
	current := end

	for current != start {
		if last, ok := cameFrom[current]; ok {
			path = append(path, current)
			current = last
		} else {
			return path, false
		}
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, true
}

func (g *Grid) getNeighbors(pos Position) []Position {
	maxWidth := len(g.tiles[pos.y]) - 1
	maxHeight := len(g.tiles) - 1

	neighbors := make([]Position, 0)
	var neighbor Position
	if pos.x > 0 {
		neighbor = Position{x: pos.x - 1, y: pos.y}
		if g.isReachable(pos, neighbor) {
			neighbors = append(neighbors, neighbor)
		}
	}
	if pos.x < maxWidth {
		neighbor = Position{x: pos.x + 1, y: pos.y}
		if g.isReachable(pos, neighbor) {
			neighbors = append(neighbors, neighbor)
		}
	}
	if pos.y > 0 {
		neighbor = Position{x: pos.x, y: pos.y - 1}
		if g.isReachable(pos, neighbor) {
			neighbors = append(neighbors, neighbor)
		}
	}
	if pos.y < maxHeight {
		neighbor = Position{x: pos.x, y: pos.y + 1}
		if g.isReachable(pos, neighbor) {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

func (g *Grid) isReachable(a Position, b Position) bool {
	aHeight := g.tiles[a.y][a.x]
	bHeight := g.tiles[b.y][b.x]

	if bHeight-aHeight > 1 {
		return false
	}

	return true
}

// priority queue implementation (https://pkg.go.dev/container/heap)
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// lower value == higher priority
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PQItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *PQItem, value Position, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
