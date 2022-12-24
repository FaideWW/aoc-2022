package main

import (
	//	"errors"
	"container/heap"
	"errors"
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
	width     int
	height    int
	entrance  Position
	exit      Position
	blizzards Blizzards
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
	steps := valley.findPath()
	fmt.Printf("goal reached in %d steps\n", steps)
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

	return Valley{
		height:    len(lines),
		width:     len(lines[0]),
		entrance:  entrance,
		exit:      exit,
		blizzards: blizzards,
	}
}

func (v *Valley) findPath() int {
	start := v.entrance
	goal := v.exit

	// memoized list of blizzard positions at a given time
	blizzardCache := make(map[int]Blizzards)
	blizzardCache[0] = v.blizzards

	getNextPositions := func(state State) []Position {
		cycleTime := state.time % (v.width - 2)
		b, ok := blizzardCache[cycleTime]
		if !ok {
			blizzardCache[cycleTime] = advanceBlizzards(blizzardCache[0], v.width, v.height, state.time)
			b = blizzardCache[cycleTime]
		}

		candidates := make([]Position, 0)
		// fmt.Printf(" computing candidates for %+v - blizzard dict %+v\n", state, b)

		north := Position{state.pos.x, state.pos.y - 1}
		east := Position{state.pos.x + 1, state.pos.y}
		south := Position{state.pos.x, state.pos.y + 1}
		west := Position{state.pos.x - 1, state.pos.y}

		// fmt.Printf("  testing north (%v) - inbounds %t blizzards? %04b\n", north, v.inBounds(north), b[north])
		// fmt.Printf("  testing east (%v) - inbounds %t blizzards? %04b\n", east, v.inBounds(east), b[east])
		// fmt.Printf("  testing south (%v) - inbounds %t blizzards? %04b\n", south, v.inBounds(south), b[south])
		// fmt.Printf("  testing west (%v) - inbounds %t blizzards? %04b\n", west, v.inBounds(west), b[west])

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

	startingState := State{start, 0}

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

	explored := make(map[State]bool)

	i := 0
	for frontier.Len() > 0 {
		i++
		current := frontier[0]
		frontier = frontier[1:]

		// if i%1000 == 0 {
		// 	fmt.Printf("iteration %d - testing %+v\n", i, current)
		// }
		if i > 500 {
			break
		}

		if current.pos == goal {
			return current.time
		}

		fmt.Printf("exploring %+v\n", current)

		candidates := getNextPositions(current)
		fmt.Printf(" candidates: %+v\n", candidates)
		for _, candidate := range candidates {
			next := State{candidate, current.time + 1}
			if !explored[next] {
				explored[next] = true
				frontier = append(frontier, next)
			}
		}
	}

	panic(errors.New("no path found"))
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

func (v *Valley) print() string {
	var output string
	for y := 0; y < v.height; y++ {
		var line string
		for x := 0; x < v.width; x++ {
			pos := Position{x, y}
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
							if ((v.blizzards[pos] >> NORTH) & 1) == 1 {
								activeBlizzards++
								tile = "^"
							}
							if ((v.blizzards[pos] >> EAST) & 1) == 1 {
								activeBlizzards++
								tile = ">"
							}
							if ((v.blizzards[pos] >> SOUTH) & 1) == 1 {
								activeBlizzards++
								tile = "v"
							}
							if ((v.blizzards[pos] >> WEST) & 1) == 1 {
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
