package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Worksheet struct {
	cells [][]string
}

func NewWorksheet(lines []string) *Worksheet {
	cells := make([][]string, len(lines))
	for i, line := range lines {
		args := strings.Fields(line)
		cells[i] = args
	}
	return &Worksheet{cells: cells}
}

func (ws Worksheet) SolveColumn(j int) int64 {
	var nums []int64
	op := ""
	for _, row := range ws.cells {
		cell := row[j]
		switch cell {
		case "+":
			op = "+"
		case "*":
			op = "*"
		default:
			num, _ := strconv.ParseInt(cell, 10, 64)
			nums = append(nums, num)
		}
	}
	return ws.CaclColumn(nums, op)
}

func (ws Worksheet) CaclColumn(nums []int64, op string) int64 {
	sol := int64(0)
	for _, num := range nums {
		if sol == 0 {
			sol = num
			continue
		}
		if op == "*" {
			sol *= num
		} else {
			sol += num
		}
	}
	return sol
}

func (ws Worksheet) SolSum() int64 {
	sum := int64(0)
	for j := 0; j < len(ws.cells[0]); j++ {
		sum += ws.SolveColumn(j)
	}
	return sum
}

func main() {
	start := time.Now()

	lines := readFile("day6/input.txt")
	ws := NewWorksheet(lines)
	fmt.Println("Part 01: ✖️➕🟰:", ws.SolSum())

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
