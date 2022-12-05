package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CargoLane []byte

type Instruction struct {
	amount int
	from   int
	to     int
}

type CargoState struct {
	lanes              []CargoLane
	instructions       []Instruction
	instructionPointer int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := readInputFile(os.Args[1])
	state := parseInput(input)
	state.executeInstructions()
	state.printTopCrates()
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) CargoState {
	lines := strings.Split(input, "\n")
	split := 0
	for i, line := range lines {
		if len(line) == 0 {
			split = i
			break
		}
	}

	return CargoState{
		lanes:              parseCrates(lines[:split]),
		instructions:       parseInstructions(lines[split+1:]),
		instructionPointer: 0,
	}
}

func parseCrates(input []string) []CargoLane {
	numLanes := len(input[0])/4 + 1
	stackSize := len(input) - 1

	lanes := make([]CargoLane, numLanes)

	for i := 0; i < numLanes; i++ {
		lane := make(CargoLane, 0)
		lanes[i] = lane
	}

	// Fill the lanes from the bottom up
	for i := stackSize - 1; i >= 0; i-- {
		height := (stackSize - 1) - i
		for lane := 0; lane < numLanes; lane++ {
			cratePointer := lane * 4
			crate := input[i][cratePointer+1]
			if crate != ' ' {
				lanes[lane] = append(lanes[lane], crate)
				fmt.Printf("Crate %s at height %d in lane %d\n", string(crate), height, lane)
			}
		}
	}

	return lanes
}

func parseInstructions(input []string) []Instruction {
	r := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	instructions := make([]Instruction, len(input)-1)

	for i, line := range input {
		if len(line) == 0 {
			continue
		}
		data := r.FindStringSubmatch(line)

		amount, _ := strconv.Atoi(data[1])
		from, _ := strconv.Atoi(data[2])
		to, _ := strconv.Atoi(data[3])

		instructions[i] = Instruction{amount: amount, from: from, to: to}
	}

	return instructions
}

func (s *CargoState) executeInstructions() {
	for s.hasNextInstruction() {
		s.executeNextInstruction()
	}
}

func (s *CargoState) hasNextInstruction() bool {
	return s.instructionPointer < len(s.instructions)
}

func (s *CargoState) executeNextInstruction() {

	if s.instructionPointer >= len(s.instructions) {
		panic(errors.New("instruction pointer past the end of the instruction list"))
	}

	instruction := s.instructions[s.instructionPointer]

	// s.executeInstruction(instruction)
	s.executeEnhancedInstruction(instruction)

	// Increment the instruction pointer
	s.instructionPointer++
}

func (s *CargoState) executeEnhancedInstruction(instruction Instruction) {
	fmt.Printf("moving %d box(es) from lane %d to lane %d\n", instruction.amount, instruction.from, instruction.to)

	poppedCrates, err := s.lanes[instruction.from-1].popMany(instruction.amount)
	if err != nil {
		panic(err)
	}

	s.lanes[instruction.to-1].pushMany(poppedCrates)
}

func (s *CargoState) executeInstruction(instruction Instruction) {
	fmt.Printf("moving 1 box from lane %d to lane %d\n", instruction.from, instruction.to)

	poppedCrate, err := s.lanes[instruction.from-1].pop()
	if err != nil {
		panic(err)
	}

	s.lanes[instruction.to-1].push(poppedCrate)

	// Recur if we need to move more than one crate
	if instruction.amount > 1 {
		nextInstruction := Instruction{
			amount: instruction.amount - 1,
			from:   instruction.from,
			to:     instruction.to,
		}
		s.executeInstruction(nextInstruction)
	}
}

func (l *CargoLane) pushMany(arr []byte) {
	*l = append(*l, arr...)
}

func (l *CargoLane) popMany(amount int) ([]byte, error) {
	if len(*l) < amount {
		return []byte{}, errors.New("tried to move more crates than are in the lane")
	}

	newLast := len(*l) - amount
	crates := (*l)[newLast:]
	*l = (*l)[:newLast]

	return crates, nil
}

func (l *CargoLane) push(b byte) {
	*l = append(*l, b)
}

func (l *CargoLane) pop() (byte, error) {
	if len(*l) == 0 {
		return 0, errors.New("lane is empty")
	}

	lastIndex := len(*l) - 1
	crate := (*l)[lastIndex]
	*l = (*l)[:lastIndex]
	return crate, nil
}

func (s *CargoState) printTopCrates() {
	fmt.Printf("Top crates: ")
	for _, lane := range s.lanes {
		fmt.Printf("%s", string(lane[len(lane)-1]))
	}

	fmt.Printf("\n")
}
