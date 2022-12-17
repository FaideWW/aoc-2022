package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Cavern struct {
	valves            []string
	usefulValves      []string
	valveFlowRates    map[string]int
	tunnelConnections map[string][]string
	distanceMatrix    map[string]map[string]int
}

type PressureEvent struct {
	state    string
	pressure int
}

const STARTING_LOCATION = "AA"
const STARTING_MINUTES = 26

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	fmt.Println(input)

	cavern := parseCavern(input)
	cavern.computeDistanceMatrix()

	maxPressure := cavern.findMaxPressure(STARTING_LOCATION, STARTING_MINUTES)
	fmt.Printf("max pressure: %d\n", maxPressure)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseCavern(input string) Cavern {
	lines := strings.Split(input, "\n")

	valveFlowRates := make(map[string]int)
	tunnelConnections := make(map[string][]string)
	closedValves := make([]string, 0)
	usefulValves := make([]string, 0)

	r := regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d+); tunnels? leads? to valves? (.+)`)
	for _, line := range lines {
		matches := r.FindStringSubmatch(line)

		currentValve := matches[1]
		flowRate, _ := strconv.Atoi(matches[2])
		connections := strings.Split(matches[3], ", ")

		closedValves = append(closedValves, currentValve)
		if flowRate > 0 {
			usefulValves = append(usefulValves, currentValve)
		}
		valveFlowRates[currentValve] = flowRate
		tunnelConnections[currentValve] = connections
	}

	cavern := Cavern{
		valves:            closedValves,
		usefulValves:      usefulValves,
		valveFlowRates:    valveFlowRates,
		tunnelConnections: tunnelConnections,
	}
	return cavern
}

func (c *Cavern) computeDistanceMatrix() {
	distances := make(map[string]map[string]int)
	for i := 0; i < len(c.valves); i++ {
		u := c.valves[i]
		distances[u] = make(map[string]int)
		for j := 0; j < len(c.valves); j++ {
			v := c.valves[j]
			distances[u][v] = 100
		}
	}

	for u, edgeList := range c.tunnelConnections {
		distances[u][u] = 0
		for _, v := range edgeList {
			distances[u][v] = 1
		}
	}

	for kI := 0; kI < len(c.valves); kI++ {
		k := c.valves[kI]
		for iI := 0; iI < len(c.valves); iI++ {
			i := c.valves[iI]
			for jI := 0; jI < len(c.valves); jI++ {
				j := c.valves[jI]
				if distances[i][j] > distances[i][k]+distances[k][j] {
					distances[i][j] = distances[i][k] + distances[k][j]
				}
			}
		}
	}

	(*c).distanceMatrix = distances
}

func sortedAppend(array []string, item string) []string {
	insertAt := -1
	for i := 0; i < len(array); i++ {
		if array[i] < item {
			continue
		} else if array[i] >= item {
			insertAt = i
			break
		}
	}

	if insertAt == -1 {
		return append(array, item)
	}
	res := append(array[:insertAt+1], array[insertAt:]...)
	res[insertAt] = item
	return res
}

// Assume all valves are sorted alphabetically already so the order is guaranteed
func serializeValveState(position string, valves map[string]bool) string {
	sortedValves := make([]string, 0)
	for valve := range valves {
		sortedValves = sortedAppend(sortedValves, valve)
	}
	return position + ":" + strings.Join(sortedValves, ",")
}

func (c *Cavern) findMaxPressure(startingLocation string, timeLimit int) int {
	memo := make(map[int]map[string]int)
	for i := 0; i < timeLimit; i++ {
		memo[i] = make(map[string]int)
	}

	var compute func(int, string, map[string]bool)
	compute = func(timeTaken int, position string, openValves map[string]bool) {
		valveState := serializeValveState(position, openValves)
		for _, nextValve := range c.usefulValves {
			// if the valve is already open, skip it
			if openValves[nextValve] {
				continue
			}

			// if opening this valve will take too long, then we can memoize the
			// value at the time limit to be equal to now (since there's nothing
			// to be gained by moving)
			timeToOpen := timeTaken + c.distanceMatrix[position][nextValve] + 1
			if timeToOpen >= timeLimit {
				if memo[timeLimit-1][valveState] < memo[timeTaken][valveState] {
					memo[timeLimit-1][valveState] = memo[timeTaken][valveState]
				}
			} else {
				additionalPressure := (timeLimit - timeToOpen) * c.valveFlowRates[nextValve]
				nextOpenValves := make(map[string]bool)
				for k, v := range openValves {
					nextOpenValves[k] = v
				}

				nextOpenValves[nextValve] = true
				nextValveState := serializeValveState(nextValve, nextOpenValves)

				nextPressure := additionalPressure + memo[timeTaken][valveState]

				if memo[timeToOpen][nextValveState] < nextPressure {
					memo[timeToOpen][nextValveState] = nextPressure
				}

				compute(timeToOpen, nextValve, nextOpenValves)
			}

		}
		// Additionally, if we've iterated over all the valves, memoize this pressure at the time limit to simulate waiting at this step
		if memo[timeLimit-1][valveState] < memo[timeTaken][valveState] {
			memo[timeLimit-1][valveState] = memo[timeTaken][valveState]
		}
	}

	// fill in the memo table with all possible permutations of valve openings
	compute(0, startingLocation, make(map[string]bool))

	// find the permutations at time=timeLimit with the highest pressure
	finalPressures := make([]PressureEvent, 0)
	for k, v := range memo[timeLimit-1] {
		finalPressures = append(finalPressures, PressureEvent{k, v})
	}

	sort.Slice(finalPressures, func(i, j int) bool {
		return finalPressures[i].pressure > finalPressures[j].pressure
	})

	// part 1
	// return finalPressures[0].pressure

	// part 2
	// now that we've memoized all states, find the two ending states where two actors open unique sets of valves that add up to the highest combined total
	maxPressure := 0
	for i := 0; i < len(finalPressures); i++ {
		firstPressure := finalPressures[i]
		for j := i + 1; j < len(finalPressures); j++ {
			secondPressure := finalPressures[j]
			if !doStatesOverlap(firstPressure.state, secondPressure.state) {
				sum := firstPressure.pressure + secondPressure.pressure
				if sum > maxPressure {
					maxPressure = sum
				}
				break
			}
		}
	}

	return maxPressure
}

func doStatesOverlap(state1 string, state2 string) bool {
	parts1 := strings.Split(strings.Split(state1, ":")[1], ",")
	parts2 := strings.Split(strings.Split(state2, ":")[1], ",")

	seen := make(map[string]bool)
	for _, p := range parts1 {
		seen[p] = true
	}

	for _, p := range parts2 {
		if seen[p] {
			return true
		}
	}

	return false
}
