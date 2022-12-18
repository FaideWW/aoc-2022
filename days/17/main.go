package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const CHAMBER_WIDTH = 7
const ROCK_SPAWN_X = 2
const ROCK_SPAWN_Y_GAP = 3
const UNIQUE_ROCKS = 5

type Rock struct {
	x      int
	y      int
	width  int
	height int
	shape  []uint8
}

type Surface struct {
	contour    []uint8
	baseHeight int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	// fmt.Println(input)
	// fmt.Println(len(input))
	height1 := simulate(input, 2022)
	fmt.Printf("tower height (pt 1): %d\n", height1)
	height2 := simulate(input, 1000000000000)
	fmt.Printf("tower height (pt 2): %d\n", height2)

}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

// for confirming the accuracy of the optimized solution
func bruteForce(jetPattern string, rockCount int) int {
	rockIndex := 0
	jetIndex := 0
	surface := Surface{
		contour:    make([]uint8, 0),
		baseHeight: 0,
	}
	for i := 0; i < rockCount; i++ {
		dropRock(&rockIndex, &surface, jetPattern, &jetIndex)
		fmt.Println(printSurface(&surface))
	}

	maxHeight := surface.baseHeight + len(surface.contour)
	return maxHeight
}

func simulate(jetPattern string, rockCount int) int {
	// there is a near 100% chance that after some amount of time, the rock
	// dropping cycle repeats (meaning the topography of the tower is exactly the
	// same at two different heights, in the same point in both the rock cycle
	// and the jet cycle). solving part 2 is a matter of identifying where that
	// cycle is and then extrapolating that repeating cycle out until we reach
	// the rock limit.
	//
	// - first we run the simulation until we find a cycle, at which point the
	//   amount of rocks before the cycle begins is `initialRocks`. We also keep
	//   track of the height at this point (initialHeight).
	// - then we take the number of rocks in the cycle and divide it into
	//   (rockCount - initialRocks) to determine how many loops there are
	//   (loopCount), and the marginal height gain of each loop (marginalHeight).
	// - then we take the remainder of that division, as we will need to run
	//   these iterations to work up to the final rock count and tower height
	//   (remainingRocks, remainingHeight)
	//
	// so all in all, the equation for the full tower height is:
	//   initialHeight + (marginalHeight * loopCount) + remainingHeight

	type CacheKey struct {
		contour   string
		rockIndex int
		jetIndex  int
	}

	type State struct {
		baseHeight int
		rockIndex  int
		jetIndex   int
	}

	stateCache := make(map[int]State)
	loopCache := make(map[CacheKey]int)

	rockIndex := 0
	jetIndex := 0

	surface := Surface{
		contour:    make([]uint8, 0),
		baseHeight: 0,
	}

	key := CacheKey{
		contour:   serializeSurfaceContour(&surface),
		rockIndex: rockIndex,
		jetIndex:  jetIndex,
	}

	loopCache[key] = 0
	stateCache[0] = State{
		baseHeight: 0,
		rockIndex:  0,
		jetIndex:   0,
	}

	var cycleStart int
	var cycleEnd int
	var cycleKey CacheKey
	cycleFound := false

	// Choose a high (but not too high) value to search for loops
	for i := 1; i < 10000; i++ {
		dropRock(&rockIndex, &surface, jetPattern, &jetIndex)
		contour := serializeSurfaceContour(&surface)
		key = CacheKey{
			contour:   contour,
			rockIndex: rockIndex,
			jetIndex:  jetIndex,
		}

		state := State{
			baseHeight: surface.baseHeight,
			rockIndex:  rockIndex,
			jetIndex:   jetIndex,
		}
		stateCache[i] = state

		if cacheEntry, ok := loopCache[key]; ok {
			cycleFound = true
			cycleStart = cacheEntry
			cycleKey = key
			cycleEnd = i
			break
		}

		loopCache[key] = i
	}

	if !cycleFound {
		panic(errors.New("no cycles found"))
	}
	fmt.Printf("cycle found from rock %d - rock %d (cache key: %+v)\n", cycleStart, cycleEnd, cycleKey)

	loopSize := cycleEnd - cycleStart
	loopCount := (rockCount - cycleStart) / loopSize
	remainder := (rockCount - cycleStart) % loopSize

	fmt.Printf("initial:%d - loopSize:%d - loopCount:%d - remainder:%d (sum: %d)\n", cycleStart, loopSize, loopCount, remainder, cycleStart+(loopSize*loopCount)+remainder)

	loopHeight := stateCache[cycleEnd].baseHeight - stateCache[cycleStart].baseHeight

	// Now we just need to run the remaining rocks, starting from the height of
	// the end of the last loop
	remainderRockIndex := cycleKey.rockIndex
	remainderJetIndex := cycleKey.jetIndex
	surface.baseHeight = stateCache[cycleStart].baseHeight + (loopHeight * loopCount)

	for i := 0; i < remainder; i++ {
		dropRock(&remainderRockIndex, &surface, jetPattern, &remainderJetIndex)
	}

	finalHeight := surface.baseHeight + len(surface.contour)
	return finalHeight
}

func serializeSurfaceContour(surface *Surface) (result string) {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(surface.contour)), ","), "[]")
}

func deserializeSurfaceContour(input string, surface *Surface) {
	surfaceBytes := strings.Split(input, ",")
	surface.contour = make([]uint8, len(surfaceBytes))
	for i := 0; i < len(surfaceBytes); i++ {
		h, _ := strconv.Atoi(surfaceBytes[i])
		(*surface).contour[i] = uint8(h)
	}
}

func dropRock(rockIndex *int, surface *Surface, jetPattern string, jetIndex *int) {
	maxHeight := surface.baseHeight + len(surface.contour)
	rock := getRock(*rockIndex)
	rock.x = ROCK_SPAWN_X
	rock.y = maxHeight + ROCK_SPAWN_Y_GAP

	settled := false
	for !settled {
		jet := jetPattern[*jetIndex]
		*jetIndex = ((*jetIndex) + 1) % len(jetPattern)
		applyJet(&rock, jet, surface)
		// Attempt to move the rock down one. If the rock's shape overlaps with
		// any of the heights, move it back up one and consider it settled.

		rock.y--
		if doesRockOverlap(&rock, surface) {
			rock.y++
			settled = true
		}
	}

	settleRock(&rock, surface)
	*rockIndex = ((*rockIndex) + 1) % UNIQUE_ROCKS
}

func settleRock(rock *Rock, surface *Surface) {

	rockOffset := rock.y - surface.baseHeight

	newRows := rockOffset + rock.height - len(surface.contour)
	if newRows > 0 {
		surface.contour = append(surface.contour, make([]uint8, newRows)...)
	}

	for y := 0; y < rock.height; y++ {
		rockBitOffset := CHAMBER_WIDTH - rock.width - rock.x
		surfaceY := rockOffset + (rock.height - y) - 1
		surface.contour[surfaceY] = surface.contour[surfaceY] | (rock.shape[y] << uint8(rockBitOffset))
	}

	// check if there is a new base height
	newSurfaceBase := 9999
	for x := 0; x < CHAMBER_WIDTH; x++ {
		hasBaseTile := false
		newColumnBase := 9999
		for y := len(surface.contour) - 1; y >= 0; y-- {
			tile := (surface.contour[y] >> (CHAMBER_WIDTH - 1 - x)) & 1

			if tile == 1 && y < newColumnBase {
				hasBaseTile = true
				newColumnBase = y
				break
			}
		}
		if !hasBaseTile {
			newColumnBase = 0
		}

		if newSurfaceBase > newColumnBase {
			newSurfaceBase = newColumnBase
		}
	}

	if newSurfaceBase > 0 {
		surface.contour = surface.contour[newSurfaceBase:]
		surface.baseHeight += newSurfaceBase
	}
}

func doesRockOverlap(rock *Rock, surface *Surface) bool {

	rockOffset := rock.y - surface.baseHeight

	if rockOffset < 0 {
		return true
	}

	rockBitOffset := CHAMBER_WIDTH - rock.width - rock.x
	for y := 0; y < rock.height; y++ {
		surfaceY := rockOffset + (rock.height - y) - 1

		if surfaceY >= len(surface.contour) || surfaceY < 0 {
			continue
		}
		overlap := surface.contour[surfaceY] & (rock.shape[y] << uint8(rockBitOffset))
		if overlap != 0 {
			return true
		}
	}

	return false
}

func applyJet(rock *Rock, direction byte, surface *Surface) {
	switch direction {
	case '<':
		{
			if rock.x > 0 {
				(*rock).x--
				if doesRockOverlap(rock, surface) {
					(*rock).x++
				}

			}
		}
	case '>':
		{
			if rock.x+rock.width < CHAMBER_WIDTH {
				(*rock).x++
				if doesRockOverlap(rock, surface) {
					(*rock).x--
				}
			}
		}
	default:
		{
			panic(errors.New("unrecognized jet rune"))
		}
	}
}

func getRock(rockIndex int) Rock {
	switch rockIndex {
	case 0:
		{
			// ####
			return Rock{
				width:  4,
				height: 1,
				shape: []uint8{
					0b1111,
				},
			}
		}
	case 1:
		{
			// .#.
			// ###
			// .#.
			return Rock{
				width:  3,
				height: 3,
				shape: []uint8{
					0b010,
					0b111,
					0b010,
				},
			}
		}
	case 2:
		{
			// ..#
			// ..#
			// ###
			return Rock{
				width:  3,
				height: 3,
				shape: []uint8{
					0b001,
					0b001,
					0b111,
				},
			}
		}
	case 3:
		{
			// #
			// #
			// #
			// #
			return Rock{
				width:  1,
				height: 4,
				shape: []uint8{
					0b1,
					0b1,
					0b1,
					0b1,
				},
			}
		}
	case 4:
		{
			// ##
			// ##
			return Rock{
				width:  2,
				height: 2,
				shape: []uint8{
					0b11,
					0b11,
				},
			}
		}
	default:
		{
			panic(errors.New("how did we get here??"))
		}
	}
}

func printSurface(surface *Surface) (output string) {
	maxHeight := surface.baseHeight + len(surface.contour)
	maxHeightStringSize := len(fmt.Sprint(maxHeight))

	for y := len(surface.contour) - 1; y >= 0; y-- {
		output += fmt.Sprintf("%*d |", maxHeightStringSize, y+surface.baseHeight)
		for x := CHAMBER_WIDTH - 1; x >= 0; x-- {
			switch (surface.contour[y] >> x) & 1 {
			case 0:
				{
					output += "."
				}
			case 1:
				{
					output += "#"
				}
			}
		}
		output += "|\n"
	}

	return
}
