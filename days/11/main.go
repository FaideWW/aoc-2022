package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const WORRY_DIVISOR = 1

type Operation struct {
	operator string
	left     string
	right    string
	leftInt  int
	rightInt int
}

type Test struct {
	condition int
	trueCase  int
	falseCase int
}

type Monkey struct {
	items       []int
	operation   Operation
	test        Test
	inspections int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))
	fmt.Println(input)
	monkeys := parseInput(input)
	printMonkeys(&monkeys)
	for i := 0; i < 10000; i++ {
		doRound(&monkeys)
	}
	printMonkeys(&monkeys)

	for i, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", i, m.inspections)
	}

	monkeyBusiness := calculateMonkeyBusiness(monkeys)
	fmt.Printf("monkey business: %d\n", monkeyBusiness)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) []Monkey {
	blocks := strings.Split(input, "\n\n")
	monkeys := make([]Monkey, len(blocks))

	for i, block := range blocks {
		monkeys[i] = parseMonkey(block)
	}

	return monkeys
}

func parseMonkey(input string) Monkey {
	lines := strings.Split(input, "\n")
	// structure of monkey notes:
	// Monkey [n]:
	//   Starting items: [...]
	//   Operation: [...]
	//   Test: divisible by [n]
	//     If true: throw to monkey [n]
	//     If false: throw to monkey [n]

	if lines[0][:6] != "Monkey" {
		panic(errors.New("error in input"))
	}

	items := parseItems(lines[1])
	operation := parseOperation(lines[2])
	test := parseTest(lines[3:])

	return Monkey{
		items:       items,
		operation:   operation,
		test:        test,
		inspections: 0,
	}
}

func parseItems(itemsStr string) (items []int) {
	for _, itemStr := range strings.Split(itemsStr[18:], ", ") {
		item, _ := strconv.Atoi(itemStr)
		items = append(items, item)
	}

	return
}

func parseOperation(operationStr string) Operation {
	tokens := strings.Split(operationStr[13:], " ")
	var leftInt, rightInt int
	leftInt, _ = strconv.Atoi(tokens[2])
	rightInt, _ = strconv.Atoi(tokens[4])

	return Operation{
		operator: tokens[3],
		left:     tokens[2],
		right:    tokens[4],
		leftInt:  leftInt,
		rightInt: rightInt,
	}
}

func parseTest(lines []string) Test {
	testStr := lines[0][8:]
	testTrueStr := lines[1][13:]
	testFalseStr := lines[2][14:]

	divisor, _ := strconv.Atoi(testStr[13:])
	ifTrue, _ := strconv.Atoi(testTrueStr[16:])
	ifFalse, _ := strconv.Atoi(testFalseStr[16:])
	return Test{
		condition: divisor,
		trueCase:  ifTrue,
		falseCase: ifFalse,
	}
}

func doRound(monkeys *[]Monkey) {
	modulo := calculateMonkeyModulo(*monkeys)

	for i, monkey := range *monkeys {
		for _, item := range monkey.items {
			newWorry := (doOperation(monkey.operation, item) / WORRY_DIVISOR) % modulo
			modWorry := newWorry % monkey.test.condition
			if modWorry == 0 {
				(*monkeys)[monkey.test.trueCase].items = append((*monkeys)[monkey.test.trueCase].items, newWorry)
			} else {
				(*monkeys)[monkey.test.falseCase].items = append((*monkeys)[monkey.test.falseCase].items, newWorry)
			}
			(*monkeys)[i].inspections++
		}

		(*monkeys)[i].items = []int{}
	}
}

func printMonkeys(monkeys *[]Monkey) {
	for i, monkey := range *monkeys {
		fmt.Printf("Monkey %d: ", i)
		for _, item := range monkey.items {
			fmt.Printf("%d, ", item)
		}
		fmt.Printf("\n")
	}
}

func doOperation(op Operation, oldValue int) (newValue int) {
	var left, right int

	if op.left == "old" {
		left = oldValue
	} else {
		left = op.leftInt
	}

	if op.right == "old" {
		right = oldValue
	} else {
		right = op.rightInt
	}

	switch op.operator {
	case "+":
		{
			newValue = left + right
		}
	case "-":
		{
			newValue = left - right
		}
	case "*":
		{
			newValue = left * right
		}
	case "/":
		{
			newValue = left / right
		}
	default:
		{
			panic(errors.New("unknown operator"))
		}
	}

	return
}

func calculateMonkeyBusiness(monkeys []Monkey) int {
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections < monkeys[j].inspections
	})

	return monkeys[len(monkeys)-1].inspections * monkeys[len(monkeys)-2].inspections
}

func calculateMonkeyModulo(monkeys []Monkey) int {
	product := 1
	for _, monkey := range monkeys {
		product *= monkey.test.condition
	}

	return product
}
