package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	cmd string
	arg int
}

type CPU struct {
	x             int
	history       []int
	displayWidth  int
	displayHeight int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	instructions := parseInput(input)

	cpu := CPU{
		x:             1,
		history:       make([]int, 0),
		displayWidth:  40,
		displayHeight: 6,
	}

	cpu.executeInstructions(instructions)

	timestamps := []int{20, 60, 100, 140, 180, 220}
	sum := 0

	for _, t := range timestamps {
		strength := cpu.getSignalStrength(t)
		fmt.Printf("signal strength at %d: %d\n", t, strength)
		sum += strength
	}

	fmt.Printf("total signal strength: %d\n", sum)

	cpu.render()
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []Instruction {
	lines := strings.Split(input, "\n")

	instructions := make([]Instruction, len(lines))

	for i, line := range lines {
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "addx":
			{
				value, _ := strconv.Atoi(parts[1])
				instructions[i] = Instruction{cmd: "addx", arg: value}
			}
		case "noop":
			{
				instructions[i] = Instruction{cmd: "noop"}
			}
		default:
			{
				panic(errors.New("unrecognized instruction"))
			}
		}
	}

	return instructions
}

func (c *CPU) executeInstructions(instructions []Instruction) {
	c.history = append(c.history, c.x)
	for _, instruction := range instructions {
		switch instruction.cmd {

		case "addx":
			{
				c.history = append(c.history, c.x)
				c.x += instruction.arg
				c.history = append(c.history, c.x)
			}
		case "noop":
			{
				c.history = append(c.history, c.x)
			}
		default:
			{
				panic(errors.New("unrecognized instruction"))
			}
		}
	}
}

func (c *CPU) getSignalStrength(t int) int {
	if t > len(c.history) {
		panic(errors.New("time is outside history"))
	}

	return c.history[t-1] * t
}

func (c *CPU) render() {
	currentCycle := 0
	for y := 0; y < c.displayHeight; y++ {
		line := ""
		for x := 0; x < c.displayWidth; x++ {
			currentCycle = y*c.displayWidth + x
			spritePosition := c.history[currentCycle]
			if x-spritePosition < 2 && spritePosition-x < 2 {
				line += "#"
			} else {
				line += "."
			}
		}
		fmt.Println(line)
	}
}
