package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Safe struct {
	currDial                     int
	maxDial                      int
	zeroPointers                 int
	isPasswordMethod0x434C49434B bool
}

func (s *Safe) Left(amount int) {
	s.ZeroPtrNew(amount, true)
	amount = amount % s.maxDial
	s.currDial += s.maxDial
	s.currDial -= amount
	s.currDial = s.currDial % s.maxDial
	s.ZeroPtr()
}

func (s *Safe) Right(amount int) {
	s.ZeroPtrNew(amount, false)
	amount = amount % s.maxDial
	s.currDial += amount
	s.currDial = s.currDial % s.maxDial
	s.ZeroPtr()
}

func (s *Safe) ZeroPtr() {
	if s.isPasswordMethod0x434C49434B {
		return
	}
	if s.currDial == 0 {
		s.zeroPointers++
	}
}

func (s *Safe) ZeroPtrNew(amount int, isLeft bool) {
	if !s.isPasswordMethod0x434C49434B {
		return
	}
	s.zeroPointers += amount / s.maxDial
	amount = amount % s.maxDial
	if s.currDial == 0 {
		return
	}
	if isLeft && s.currDial-amount <= 0 {
		s.zeroPointers++
	} else if !isLeft && s.currDial+amount >= s.maxDial {
		s.zeroPointers++
	}
}

func NewSafe(currDial int, maxDial int, isPasswordMethod0x434C49434B bool) *Safe {
	return &Safe{
		currDial:                     currDial,
		maxDial:                      maxDial,
		isPasswordMethod0x434C49434B: isPasswordMethod0x434C49434B,
	}
}

type Instruction struct {
	dir    int
	amount int
}

func ParseInstruction(line string) Instruction {
	dir := -1
	switch line[0] {
	case 'L':
		dir = 0
	case 'R':
		dir = 1
	default:
		panic("Unrecognized instruction: " + line)
	}
	amount, err := strconv.Atoi(line[1:])
	if err != nil {
		panic(err)
	}
	return Instruction{dir, amount}
}

func Part01(instructions []Instruction) {
	safe := NewSafe(50, 100, false)
	for _, inst := range instructions {
		if inst.dir == 0 {
			safe.Left(inst.amount)
		} else {
			safe.Right(inst.amount)
		}
	}
	fmt.Println("Part01:", safe.zeroPointers)
}

func Part02(instructions []Instruction) {
	safe := NewSafe(50, 100, true)
	for _, inst := range instructions {
		if inst.dir == 0 {
			safe.Left(inst.amount)
		} else {
			safe.Right(inst.amount)
		}
	}
	fmt.Println("Part02:", safe.zeroPointers)
}

func main() {
	start := time.Now()

	lines := readFile("day1/input.txt")
	instructions := make([]Instruction, len(lines))
	for i, line := range lines {
		instructions[i] = ParseInstruction(line)
	}

	Part01(instructions)
	Part02(instructions)

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	return split
}
