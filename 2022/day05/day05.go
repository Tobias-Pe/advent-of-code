package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Stack []string

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack) Push(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Push a new value onto the stack
func (s *Stack) AddStart(str string) {
	*s = append([]string{str}, *s...)
}

func (s *Stack) AppendToEnd(elements []string) {
	*s = append(*s, elements...)
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Peek() string {
	if s.IsEmpty() {
		return ""
	} else {
		index := len(*s) - 1 // Get the index of the top most element.
		return (*s)[index]   // Index into the slice and obtain the element.
	}
}

var stacks []*Stack
var instructions []instruction

type instruction struct {
	times int
	from  int
	to    int
}

func (instr instruction) executeAsCrateMover9000() {
	for i := 0; i < instr.times; i++ {
		popped, hadElement := stacks[instr.from].Pop()
		fmt.Println(popped, hadElement)
		if hadElement {
			stacks[instr.to].Push(popped)
		}
	}
}

func (instr instruction) executeAsCrateMover9001() {
	fromStack := *stacks[instr.from]
	var elements Stack
	if len(fromStack) >= instr.times {
		elements = fromStack[len(fromStack)-instr.times:]
	} else {
		fmt.Println("here2")
		elements = fromStack
	}
	fmt.Println(instr.from, instr.to, instr.times, elements, *stacks[instr.from], *stacks[instr.to])
	stacks[instr.to].AppendToEnd(elements)
	for i := 0; i < instr.times; i++ {
		stacks[instr.from].Pop()
	}
	fmt.Println(instr.from, instr.to, instr.times, elements, *stacks[instr.from], *stacks[instr.to])
}

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	allLines := string(content)
	lines := strings.Split(allLines, "\n")
	stackLines := lines[:8]
	for _, line := range stackLines {
		currentStackIndex := 0
		for i := 1; i < len(line); i += 4 {
			if line[i] != ' ' {
				stacks[currentStackIndex].AddStart(string(line[i]))
			}
			currentStackIndex++
		}
	}
	instructionLines := lines[10:]
	for _, line := range instructionLines {
		line = strings.ReplaceAll(line, " ", "")
		splitLine := strings.Split(line, "from")
		if len(splitLine) >= 2 {
			fromToPart := strings.Split(splitLine[1], "to")
			timesPart := strings.Split(splitLine[0], "move")[1]

			timesInt, _ := strconv.Atoi(timesPart)
			fromInt, _ := strconv.Atoi(fromToPart[0])
			toInt, _ := strconv.Atoi(fromToPart[1])
			instructions = append(instructions, instruction{
				times: timesInt,
				from:  fromInt - 1,
				to:    toInt - 1,
			})
		}
	}
}

func main() {
	stacks = make([]*Stack, 9)
	for i := range stacks {
		stacks[i] = &Stack{}
	}
	readFile()
	for _, instr := range instructions {
		instr.executeAsCrateMover9001()
	}
	for i := range stacks {
		fmt.Print(stacks[i].Peek())
	}
}
