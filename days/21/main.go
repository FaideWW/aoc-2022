package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Node interface {
	GetType() string
	GetName() string
	Eval() int
	Contains(string) bool
}

type Value struct {
	name  string
	value int
}

type Variable struct {
	name  string
	value string
}

type Operation struct {
	name        string
	operator    string
	left, right Node
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	input := strings.TrimSpace(readInputFile(os.Args[1]))

	nodes := parseInput(input)
	root := replaceRoot(nodes, "root")

	reorder(&root, "humn")

	result := root.right.Eval()
	fmt.Printf("humn=%d\n", result)
}

func readInputFile(filename string) string {
	dat, err := os.ReadFile(filename)
	check(err)
	return string(dat)
}

func parseInput(input string) map[string]Node {
	lines := strings.Split(input, "\n")
	nodes := make(map[string]Node)

	for _, line := range lines {
		parts := strings.Split(line, ": ")
		nodes[parts[0]] = parseNode(parts[0], parts[1])
	}

	return nodes
}

func parseNode(name string, input string) Node {
	components := strings.Split(input, " ")
	if len(components) == 1 {
		if name == "humn" {
			return Variable{name: name, value: name}
		}
		value, _ := strconv.Atoi(components[0])
		return Value{name, value}
	} else if len(components) == 3 {
		var left, right Node
		leftValue, err := strconv.Atoi(components[0])
		if err != nil {
			left = Variable{components[0], components[0]}
		} else {
			left = Value{"", leftValue}
		}
		rightValue, _ := strconv.Atoi(components[2])
		if err != nil {
			right = Variable{components[2], components[2]}
		} else {
			right = Value{"", rightValue}
		}

		operator := components[1]
		if name == "root" {
			operator = "="
		}

		return Operation{name: name, operator: operator, left: left, right: right}
	} else {
		panic(errors.New("unrecognized expression"))
	}
}

func replaceRoot(nodes map[string]Node, root string) Operation {
	var replace func(node Node) Node
	replace = func(node Node) Node {
		switch node.GetType() {
		case "value":
			{
				return node
			}
		case "variable":
			{
				value := nodes[node.(Variable).value]
				if value.GetType() != "variable" {
					return replace(value)
				}
				return node
			}
		case "operation":
			{
				op := node.(Operation)
				return Operation{
					name:     op.name,
					operator: op.operator,
					left:     replace(op.left),
					right:    replace(op.right),
				}
			}
		default:
			{
				panic(errors.New("unrecognized node"))
			}
		}
	}

	return replace(nodes[root]).(Operation)
}

func reorder(root *Operation, target string) {
	if root.operator != "=" {
		panic(errors.New("root must be an equality"))
	}

	leftContains := root.left.Contains(target)
	rightContains := root.right.Contains(target)
	if !leftContains && !rightContains {
		panic(errors.New("target not found"))
	}

	if rightContains {
		root.left, root.right = root.right, root.left
	}

	inverseOperators := map[string]string{
		"+": "-",
		"-": "+",
		"*": "/",
		"/": "*",
	}

	// Rotate the tree until the left side contains just the target node
	for root.left.GetName() != target {
		leftSubchildContainsTarget := root.left.(Operation).left.Contains(target)
		var toMove Node
		if leftSubchildContainsTarget {
			toMove = root.left.(Operation).right
		} else {
			toMove = root.left.(Operation).left
		}

		op := root.left.(Operation).operator
		if op == "+" || op == "*" || leftSubchildContainsTarget {
			root.right = Operation{"", inverseOperators[op], root.right, toMove}
		} else {
			root.right = Operation{"", op, toMove, root.right}
		}

		if leftSubchildContainsTarget {
			root.left = root.left.(Operation).left
		} else {
			root.left = root.left.(Operation).right
		}
	}
}

func (v Value) Eval() int {
	return v.value
}

func (v Value) GetType() string {
	return "value"
}

func (v Value) GetName() string {
	return v.name
}

func (v Value) Contains(name string) bool {
	return v.GetName() == name
}

func (v Variable) Eval() int {
	panic(errors.New("tried to eval a variable; it must be replaced first"))
}

func (v Variable) GetType() string {
	return "variable"
}

func (v Variable) GetName() string {
	return v.name
}

func (v Variable) Contains(name string) bool {
	return v.GetName() == name
}

func (o Operation) Eval() int {
	switch o.operator {
	case "+":
		{
			return o.left.Eval() + o.right.Eval()
		}
	case "-":
		{
			return o.left.Eval() - o.right.Eval()
		}
	case "*":
		{
			return o.left.Eval() * o.right.Eval()
		}
	case "/":
		{
			return o.left.Eval() / o.right.Eval()
		}
	case "=":
		{
			equal := o.left.Eval() == o.right.Eval()
			if equal {
				return 1
			} else {
				return 0
			}
		}
	default:
		{
			panic(errors.New("unrecognized operator"))
		}
	}
}

func (o Operation) GetType() string {
	return "operation"
}

func (o Operation) GetName() string {
	return o.name
}

func (o Operation) Contains(name string) bool {
	return o.GetName() == name || o.left.Contains(name) || o.right.Contains(name)
}
