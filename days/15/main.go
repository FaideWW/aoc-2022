package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Position struct {
	x int
	y int
}

type Sensor struct {
	pos    Position
	radius int
}

type Cavern struct {
	sensors  map[Position]Sensor
	beacons  map[Position]bool
	minRange int
	maxRange int
	maxDepth int
}

type Interval struct {
	min int
	max int
}

type Level struct {
	min       int
	max       int
	intervals []Interval
}

// part 1
// const Y_LEVEL = 10
// const SEARCH_AREA = 20
// const TUNING_CONSTANT = 4000000

// part 2
const Y_LEVEL = 2000000
const SEARCH_AREA = 4000000
const TUNING_CONSTANT = 4000000

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	fmt.Println(input)

	cavern := parseCavern(input)

	coveredTiles := cavern.findLevelCoverage(Y_LEVEL)

	fmt.Printf("%d covered tiles at y=%d\n", coveredTiles, Y_LEVEL)

	beacon := cavern.findMissingBeacon(SEARCH_AREA)
	tuningFreq := beacon.x*TUNING_CONSTANT + beacon.y
	fmt.Printf("beacon found at %+v. tuning frequency: %d\n", beacon, tuningFreq)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseCavern(input string) Cavern {
	r := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)

	lines := strings.Split(input, "\n")
	sensors := make(map[Position]Sensor)
	beacons := make(map[Position]bool)

	minRange := math.MaxInt
	maxRange := math.MinInt
	maxDepth := math.MinInt

	for _, line := range lines {
		data := r.FindStringSubmatch(line)
		sx, _ := strconv.Atoi(data[1])
		sy, _ := strconv.Atoi(data[2])
		bx, _ := strconv.Atoi(data[3])
		by, _ := strconv.Atoi(data[4])

		if maxDepth < sy {
			maxDepth = sy
		}
		if maxDepth < by {
			maxDepth = by
		}

		if minRange > sx {
			minRange = sx
		}
		if minRange > bx {
			minRange = bx
		}

		if maxRange < sx {
			maxRange = sx
		}
		if maxRange < bx {
			maxRange = bx
		}

		sensorPos := Position{x: sx, y: sy}
		beaconPos := Position{x: bx, y: by}

		sensors[sensorPos] = Sensor{
			pos:    sensorPos,
			radius: calculateManhattanDistance(sensorPos, beaconPos),
		}

		if !beacons[beaconPos] {
			beacons[beaconPos] = true
		}
	}

	return Cavern{
		sensors:  sensors,
		beacons:  beacons,
		minRange: minRange,
		maxRange: maxRange,
		maxDepth: maxDepth,
	}
}

func calculateManhattanDistance(a Position, b Position) int {
	dx := a.x - b.x
	dy := a.y - b.y

	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}

	return dx + dy
}

func (c *Cavern) hasSensor(pos Position) bool {
	_, ok := c.sensors[pos]
	return ok
}

func (c *Cavern) hasBeacon(pos Position) bool {
	return c.beacons[pos]
}

func (c *Cavern) findLevelCoverage(y int) (coveredTiles int) {
	coverage := make(map[int]bool)
	for pos, sensor := range c.sensors {
		dy := y - pos.y
		if dy < 0 {
			dy *= -1
		}
		xRange := (sensor.radius - dy)
		// If the radius overlaps with the requested level, the covered tiles will
		// be the interval [-xRange, xRange]
		if xRange > 0 {
			for x := -xRange; x <= xRange; x++ {
				testPos := Position{
					x: sensor.pos.x + x,
					y: y,
				}
				if !coverage[testPos.x] && !c.hasBeacon(testPos) {
					coveredTiles++
					coverage[testPos.x] = true
				}
			}
		}
	}

	return
}

func (s *Sensor) isInRadius(pos Position) bool {
	dy := s.pos.y - pos.y
	if dy < 0 {
		dy *= -1
	}
	dx := s.pos.x - pos.x
	if dx < 0 {
		dx *= -1
	}
	return dy+dx <= s.radius
}

func (c *Cavern) findMissingBeacon(maxCoord int) Position {
	fullLevels := Level{min: 0, max: maxCoord, intervals: make([]Interval, 0)}
	levels := make(map[int]*Level, 0)

	for _, sensor := range c.sensors {
		for y := sensor.pos.y - sensor.radius; y <= sensor.pos.y+sensor.radius; y++ {
			if y < 0 || y > maxCoord {
				continue
			}
			if fullLevels.covers(y) {
				continue
			}
			_, ok := levels[y]
			if !ok {
				levels[y] = &Level{min: 0, max: maxCoord, intervals: make([]Interval, 0)}
			}
			dy := y - sensor.pos.y
			if dy < 0 {
				dy *= -1
			}
			xRange := sensor.radius - dy
			xInterval := Interval{
				min: sensor.pos.x - xRange,
				max: sensor.pos.x + xRange,
			}
			filled := levels[y].addInterval(xInterval)

			if filled {
				fullLevels.addSingle(y)
				delete(levels, y)
			}
		}
	}

	if len(levels) != 1 {
		panic(errors.New("found more or fewer than 1 candidate level"))
	}

	y := -1
	for i := range levels {
		y = i
	}

	x := -1

	if len(levels[y].intervals) == 1 {
		if levels[y].intervals[0].min == 1 {
			x = 0
		} else {
			x = maxCoord
		}
	} else {
		x = levels[y].intervals[0].max + 1
	}

	return Position{x, y}
}

// returns whether the entire range [0, l.max] is covered
func (l *Level) isFull() bool {
	return len(l.intervals) == 1 && l.intervals[0] == Interval{min: l.min, max: l.max}
}

// returns whether the value is within one of the level's intervals
func (l *Level) covers(value int) bool {
	for _, i := range l.intervals {
		// intervals are sorted, so if the first interval is larger than the value we can exit early
		if i.min > value {
			return false
		}
		if i.min <= value && i.max >= value {
			return true
		}
	}

	return false
}

// adds a new interval to the level, merging existing levels where possible
func (l *Level) addInterval(toAdd Interval) (isFilled bool) {

	// clamp the interval to the min and max
	if toAdd.min < l.min {
		toAdd.min = l.min
	}
	if toAdd.max > l.max {
		toAdd.max = l.max
	}

	insertAt := -1
	for i := 0; i < len(l.intervals); i++ {
		interval := l.intervals[i]
		if interval.min < toAdd.min {
			continue
		} else if interval.min > toAdd.min {
			insertAt = i
			break
		} else if interval.min == toAdd.min {
			// If the intervals exactly match, we don't need to add it
			if interval.max == toAdd.max {
				return l.isFull()
			}

			insertAt = i
			break
		}
	}
	if insertAt == -1 {
		l.intervals = append(l.intervals, toAdd)
	} else {
		l.intervals = append(l.intervals[:insertAt+1], l.intervals[insertAt:]...)
		l.intervals[insertAt] = toAdd
	}

	l.mergeLevels()
	return l.isFull()
}

func (l *Level) addSingle(value int) (isFilled bool) {
	return l.addInterval(Interval{min: value, max: value})
}

func (l *Level) mergeLevels() {
	merged := make([]Interval, 0)
	for _, interval := range l.intervals {
		if len(merged) == 0 || merged[len(merged)-1].max+1 < interval.min {
			merged = append(merged, interval)
		} else {
			if merged[len(merged)-1].max < interval.max {
				merged[len(merged)-1].max = interval.max
			}
		}
	}

	l.intervals = merged
}

func (c *Cavern) print() string {
	var printout string
	// print headers. assume all headers are 3 digits at most
	rangeAxisHeight := len(fmt.Sprint(c.maxRange)) + 1
	depthAxisLength := len(fmt.Sprint(c.maxDepth)) + 1
	for y := 0; y < rangeAxisHeight; y++ {
		var line string
		for x := c.minRange - (depthAxisLength + 1); x < c.maxRange+1; x++ {
			xStr := fmt.Sprint(x)
			if x >= 0 && x%5 == 0 {
				digit := len(xStr) - (rangeAxisHeight - y)
				if digit >= 0 {
					line += string(xStr[digit])
				} else {
					line += " "
				}
			} else {
				line += " "
			}
		}
		printout += fmt.Sprintln(line)
	}

	for y := 0; y < c.maxDepth+1; y++ {
		var line string
		currentDepthSize := len(fmt.Sprint(y))
		for i := 0; i < depthAxisLength-currentDepthSize; i++ {
			line += " "
		}
		line += fmt.Sprintf("%d ", y)

		for x := c.minRange; x < c.maxRange+1; x++ {
			pos := Position{x: x, y: y}
			if c.hasSensor(pos) {
				line += "S"
			} else if c.hasBeacon(pos) {
				line += "B"
			} else {
				line += "."
			}
		}
		printout += fmt.Sprintln(line)
	}

	return printout
}
