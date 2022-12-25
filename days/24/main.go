package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
)

type Direction int

const (
	NORTH Direction = 0
	EAST  Direction = 1
	SOUTH Direction = 2
	WEST  Direction = 3
)

type Position struct {
	x int
	y int
}

type Blizzards map[Position]byte

type Valley struct {
	width         int
	height        int
	entrance      Position
	exit          Position
	blizzardCache map[int]Blizzards
}

type State struct {
	pos  Position
	time int
}

type PQItem struct {
	value    State
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
	valley := parseInput(input)
	fmt.Printf("%+v\n", valley)
	steps1 := valley.findPath(valley.entrance, valley.exit, 0)
	fmt.Printf("goal reached in %d steps\n", steps1)
	steps2 := valley.findPath(valley.exit, valley.entrance, steps1)
	fmt.Printf("start reached in %d steps\n", steps2)
	steps3 := valley.findPath(valley.entrance, valley.exit, steps2)
	fmt.Printf("goal reached in %d steps\n", steps3)

}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) Valley {
	lines := strings.Split(input, "\n")
	blizzards := make(Blizzards)

	var entrance Position
	var exit Position
	for y := 0; y < len(lines); y++ {
		line := lines[y]
		for x := 0; x < len(line); x++ {
			tile := line[x]
			pos := Position{x, y}
			switch tile {
			case '>':
				{
					blizzards[pos] = blizzards[pos] | (1 << EAST)
				}
			case '<':
				{
					blizzards[pos] = blizzards[pos] | (1 << WEST)
				}
			case '^':
				{
					blizzards[pos] = blizzards[pos] | (1 << NORTH)
				}
			case 'v':
				{
					blizzards[pos] = blizzards[pos] | (1 << SOUTH)
				}
			case '.':
				{
					if y == 0 {
						entrance = pos
					} else if y == len(lines)-1 {
						exit = pos
					}
				}
			default:
				{
					continue
				}
			}
		}
	}

	blizzardCache := make(map[int]Blizzards)
	blizzardCache[0] = blizzards

	return Valley{
		height:        len(lines),
		width:         len(lines[0]),
		entrance:      entrance,
		exit:          exit,
		blizzardCache: blizzardCache,
	}
}

// A* search for shortest path from start to goal, starting at initialTime
// (important for blizzard positioning)
func (v *Valley) findPath(start Position, goal Position, initialTime int) int {
	if _, ok := v.blizzardCache[initialTime]; !ok {
		v.blizzardCache[initialTime] = advanceBlizzards(v.blizzardCache[0], v.width, v.height, initialTime)
	}

	heuristic := func(nextState State) int {
		// heruistic is the manhattan distance to the goal
		return (goal.x - nextState.pos.x) + (goal.y - nextState.pos.y)
	}

	getNextPositions := func(state State) []Position {
		nextTime := (state.time + 1)
		// I believe there's a repeating cycle of blizzards based on the width of
		// the valley, since blizzards move at a constant rate in one direction
		// forever. However, this optimization doesn't seem necessary (it runs fast
		// enough without) and using the below line doesn't give the correct
		// result, so noting that there is time to be gained by figuring this out.
		// cycleTime := (state.time + 1) % (v.width - 2)
		b, ok := v.blizzardCache[nextTime]
		if !ok {
			v.blizzardCache[nextTime] = advanceBlizzards(v.blizzardCache[0], v.width, v.height, nextTime)
			b = v.blizzardCache[nextTime]
		}

		candidates := make([]Position, 0)

		north := Position{state.pos.x, state.pos.y - 1}
		east := Position{state.pos.x + 1, state.pos.y}
		south := Position{state.pos.x, state.pos.y + 1}
		west := Position{state.pos.x - 1, state.pos.y}

		if v.inBounds(north) && b[north] == 0 {
			candidates = append(candidates, north)
		}
		if v.inBounds(east) && b[east] == 0 {
			candidates = append(candidates, east)
		}
		if v.inBounds(south) && b[south] == 0 {
			candidates = append(candidates, south)
		}
		if v.inBounds(west) && b[west] == 0 {
			candidates = append(candidates, west)
		}

		if b[state.pos] == 0 {
			candidates = append(candidates, state.pos)
		}

		return candidates
	}

	startingState := State{start, initialTime}

	frontier := make(PriorityQueue, 1)
	frontier[0] = &PQItem{
		value:    startingState,
		priority: 0,
		index:    0,
	}
	heap.Init(&frontier)

	cameFrom := make(map[State]State)
	costs := make(map[State]int)

	costs[startingState] = 0

	var goalState State
	for frontier.Len() > 0 {
		current := heap.Pop(&frontier).(*PQItem).value

		if current.pos == goal {
			goalState = current
			break
		}

		candidates := getNextPositions(current)
		for _, candidate := range candidates {
			next := State{candidate, current.time + 1}
			newCost := costs[current] + 1
			foundCost, ok := costs[next]
			if !ok || foundCost < costs[next] {
				costs[next] = newCost
				priority := newCost + heuristic(next)
				heap.Push(&frontier, &PQItem{value: next, priority: priority})
				cameFrom[next] = current
			}
		}
	}

	return goalState.time
}

// Advance all of the blizzards 1 step
func advanceBlizzards(blizzards Blizzards, width int, height int, steps int) Blizzards {
	nextBlizzards := make(Blizzards)
	// There's probably an optimized bitwise operation that can instantly update
	// the entire grid, but for now we'll do it one at a time.

	// Assumption: the entrance and exit are considered "walls" for the purposes
	// of wrapping blizzards around the grid. Ie, a northward blizzard can't
	// wander into the entrance tile (it will wrap to the bottom after (1,1)
	for pos, mask := range blizzards {
		if ((mask >> NORTH) & 1) == 1 {
			north := Position{
				pos.x,
				wrapInt(pos.y-steps, 1, height-1),
			}
			nextBlizzards[north] = nextBlizzards[north] | (1 << NORTH)
		}
		if ((mask >> EAST) & 1) == 1 {
			east := Position{
				wrapInt(pos.x+steps, 1, width-1),
				pos.y,
			}
			nextBlizzards[east] = nextBlizzards[east] | (1 << EAST)
		}
		if ((mask >> SOUTH) & 1) == 1 {
			south := Position{
				pos.x,
				wrapInt(pos.y+steps, 1, height-1),
			}
			nextBlizzards[south] = nextBlizzards[south] | (1 << SOUTH)
		}
		if ((mask >> WEST) & 1) == 1 {
			west := Position{
				wrapInt(pos.x-steps, 1, width-1),
				pos.y,
			}
			nextBlizzards[west] = nextBlizzards[west] | (1 << WEST)
		}
	}

	return nextBlizzards
}

// Wraps a value around [min, max)
func wrapInt(x, min, max int) int {
	interval := max - min
	offset := x - min
	return min + (interval+offset%interval)%interval
}

func (v *Valley) inBounds(p Position) bool {
	if p == v.entrance {
		return true
	}
	if p == v.exit {
		return true
	}
	return p.x > 0 && p.x < v.width-1 && p.y > 0 && p.y < v.height-1
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

func (pq *PriorityQueue) update(item *PQItem, value State, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func (v *Valley) print(blizzards Blizzards, player Position) string {
	var output string
	for y := 0; y < v.height; y++ {
		var line string
		for x := 0; x < v.width; x++ {
			pos := Position{x, y}
			if player == pos {
				line += "E"
				continue
			}
			switch y {
			case 0:
				{
					if pos == v.entrance {
						line += "."
					} else {
						line += "#"
					}
				}
			case v.height - 1:
				{
					if pos == v.exit {
						line += "."
					} else {
						line += "#"
					}

				}
			default:
				{
					switch x {
					case 0:
						{
							line += "#"
						}
					case v.width - 1:
						{
							line += "#"
						}
					default:
						{
							activeBlizzards := 0
							var tile string
							// count how many blizzards are active
							if ((blizzards[pos] >> NORTH) & 1) == 1 {
								activeBlizzards++
								tile = "^"
							}
							if ((blizzards[pos] >> EAST) & 1) == 1 {
								activeBlizzards++
								tile = ">"
							}
							if ((blizzards[pos] >> SOUTH) & 1) == 1 {
								activeBlizzards++
								tile = "v"
							}
							if ((blizzards[pos] >> WEST) & 1) == 1 {
								activeBlizzards++
								tile = "<"
							}
							switch activeBlizzards {
							case 0:
								{
									line += "."
								}
							case 1:
								{
									line += tile
								}
							default:
								{
									line += fmt.Sprint(activeBlizzards)
								}
							}
						}
					}

				}
			}
		}
		output += fmt.Sprintf("%s\n", line)
	}
	return output
}
