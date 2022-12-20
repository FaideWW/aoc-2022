package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type State struct {
	ore          int
	clay         int
	obsidian     int
	geodes       int
	oreBots      int
	clayBots     int
	obsidianBots int
	geodeBots    int
	timeLeft     int
}

type Blueprint struct {
	id                   int
	oreBotCost           int
	clayBotCost          int
	obsidianBotCostOre   int
	obsidianBotCostClay  int
	geodeBotCostOre      int
	geodeBotCostObsidian int
	maxOreCost           int
}

func max(vs ...int) int {
	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	blueprints := parseInput(input)

	sum := 0
	for _, blueprint := range blueprints {
		state := newState(24)
		geodes := findOptimalGeodePath(blueprint, state)
		qualityLevel := geodes * blueprint.id
		sum += qualityLevel
	}

	fmt.Printf("total quality level: %d\n", sum)

	product := 1

	var firstThreeBlueprints []Blueprint
	if len(blueprints) < 3 {
		firstThreeBlueprints = blueprints
	} else {
		firstThreeBlueprints = blueprints[:3]
	}

	for _, blueprint := range firstThreeBlueprints {

		state := newState(32)
		geodes := findOptimalGeodePath(blueprint, state)
		fmt.Printf("testing blueprint %d - %d geodes\n", blueprint.id, geodes)
		product *= geodes
	}
	fmt.Printf("geode product: %d\n", product)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []Blueprint {
	r := regexp.MustCompile(`Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.`)

	lines := strings.Split(input, "\n")

	blueprints := make([]Blueprint, len(lines))

	for i, line := range lines {
		data := r.FindStringSubmatch(line)
		id, _ := strconv.Atoi(data[1])
		oreBotCost, _ := strconv.Atoi(data[2])
		clayBotCost, _ := strconv.Atoi(data[3])
		obsidianBotCostOre, _ := strconv.Atoi(data[4])
		obsidianBotCostClay, _ := strconv.Atoi(data[5])
		geodeBotCostOre, _ := strconv.Atoi(data[6])
		geodeBotCostObsidian, _ := strconv.Atoi(data[7])

		maxOreCost := max(oreBotCost, clayBotCost, obsidianBotCostOre, geodeBotCostOre)

		blueprints[i] = Blueprint{
			id:                   id,
			oreBotCost:           oreBotCost,
			clayBotCost:          clayBotCost,
			obsidianBotCostOre:   obsidianBotCostOre,
			obsidianBotCostClay:  obsidianBotCostClay,
			geodeBotCostOre:      geodeBotCostOre,
			geodeBotCostObsidian: geodeBotCostObsidian,
			maxOreCost:           maxOreCost,
		}
	}
	return blueprints
}

func newState(timeLimit int) State {
	return State{
		oreBots:  1,
		timeLeft: timeLimit,
	}
}

// run a DFS to find the maximum geode count after the time limit
func findOptimalGeodePath(blueprint Blueprint, initialState State) int {

	globalMax := 0

	var dfs func(State) int
	dfs = func(current State) int {
		// check if we can beat the global max in a best case scenario (all
		// remaining turns are building new geode bots). if not, there's no point
		// in exploring this branch further
		potential := getMaxPotentialGeodes(current)
		if current.timeLeft == 0 || globalMax >= current.geodes+potential {
			return 0
		}

		// check if we are actually in the best case. if we are, we can fast
		// forward the rest of this branch
		if current.oreBots >= blueprint.geodeBotCostOre && current.obsidianBots >= blueprint.geodeBotCostObsidian {
			return potential
		}

		maxGeodes := 0
		for _, next := range getOptions(blueprint, current) {
			nextGeodes := current.geodeBots + dfs(next)
			if nextGeodes > maxGeodes {
				maxGeodes = nextGeodes
			}
		}

		if maxGeodes > globalMax {
			globalMax = maxGeodes
		}
		return maxGeodes
	}

	maxGeodes := dfs(initialState)

	return maxGeodes
}

func getOptions(blueprint Blueprint, state State) []State {
	states := make([]State, 0)

	// Stop checking paths where we make non-geode bots if we have sufficient ore
	// generation to make the most expensive bot every turn
	shouldMakeMoreOreBots := blueprint.maxOreCost > state.oreBots
	shouldMakeMoreClayBots := blueprint.obsidianBotCostClay > state.clayBots
	shouldMakeMoreObsidianBots := blueprint.geodeBotCostObsidian > state.obsidianBots

	if shouldMakeMoreOreBots {
		doNothing := state
		doNothing.collectOres()
		doNothing.timeLeft--
		states = append(states, doNothing)
	}

	if shouldMakeMoreOreBots && state.ore >= blueprint.oreBotCost {
		buildOreBot := state
		buildOreBot.collectOres()
		buildOreBot.timeLeft--
		buildOreBot.ore -= blueprint.oreBotCost
		buildOreBot.oreBots++
		states = append(states, buildOreBot)
	}

	if shouldMakeMoreClayBots && state.ore >= blueprint.clayBotCost {
		buildClayBot := state
		buildClayBot.collectOres()
		buildClayBot.timeLeft--
		buildClayBot.ore -= blueprint.clayBotCost
		buildClayBot.clayBots++
		states = append(states, buildClayBot)
	}

	if shouldMakeMoreObsidianBots && state.ore >= blueprint.obsidianBotCostOre && state.clay >= blueprint.obsidianBotCostClay {
		buildObsidianBot := state
		buildObsidianBot.collectOres()
		buildObsidianBot.timeLeft--
		buildObsidianBot.ore -= blueprint.obsidianBotCostOre
		buildObsidianBot.clay -= blueprint.obsidianBotCostClay
		buildObsidianBot.obsidianBots++
		states = append(states, buildObsidianBot)
	}

	if state.ore >= blueprint.geodeBotCostOre && state.obsidian >= blueprint.geodeBotCostObsidian {
		buildGeodeBot := state
		buildGeodeBot.collectOres()
		buildGeodeBot.timeLeft--
		buildGeodeBot.ore -= blueprint.geodeBotCostOre
		buildGeodeBot.obsidian -= blueprint.geodeBotCostObsidian
		buildGeodeBot.geodeBots++
		states = append(states, buildGeodeBot)
	}

	return states
}

// decrements the timer and collects the resources from bots
func (s *State) collectOres() {
	s.ore += s.oreBots
	s.clay += s.clayBots
	s.obsidian += s.obsidianBots
	s.geodes += s.geodeBots
}

func getMaxPotentialGeodes(s State) int {
	// Assume we have sufficient resources to build 1 geode bot per turn
	maxPotentialGeodes := 0
	for i := s.timeLeft - 1; i >= 0; i-- {
		maxPotentialGeodes += s.geodeBots
		s.geodeBots++
	}
	return maxPotentialGeodes
}
