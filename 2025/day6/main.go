package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
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

type Column struct {
	op  string
	num int
}

func NewColumn(line string) *Column {
	if strings.TrimSpace(line) == "" {
		return nil
	}
	numStr := ""
	c := Column{}
	for _, char := range line {
		if unicode.IsDigit(char) {
			numStr += string(char)
		} else if !unicode.IsSpace(char) {
			c.op = string(char)
		}
	}
	num, _ := strconv.Atoi(numStr)
	c.num = num
	return &c
}

type TransposeWorksheet struct {
	cols []*Column
}

func NewTransposeWorksheet(lines []string) *TransposeWorksheet {
	tWs := TransposeWorksheet{cols: []*Column{}}
	maxLineLen := 0
	for _, line := range lines {
		maxLineLen = max(maxLineLen, len(line))
	}
	for col := 0; col < maxLineLen; col++ {
		transposedLine := ""
		for _, row := range lines {
			if col >= len(row) {
				continue
			}
			transposedLine += string(row[col])
		}
		newColumn := NewColumn(transposedLine)
		if newColumn != nil {
			tWs.cols = append(tWs.cols, newColumn)
		}
	}
	return &tWs
}

func (ws TransposeWorksheet) SolSum() int64 {
	sum := int64(0)
	tmpSol := int64(0)
	tmpOp := ""
	for _, col := range ws.cols {
		if col.op != "" {
			//fmt.Println("🟰", tmpSol)
			sum += tmpSol
			tmpOp = col.op
			tmpSol = int64(col.num)
			//fmt.Print(col.op, "\t")
		} else if tmpOp == "*" {
			tmpSol *= int64(col.num)
		} else {
			tmpSol += int64(col.num)
		}
		//fmt.Print(col.num, "\t")
	}

	//fmt.Println("🟰", tmpSol)
	return sum + tmpSol
}

func main() {
	start := time.Now()

	lines := readFile("day6/input.txt")
	ws := NewWorksheet(lines)
	fmt.Println("Part 01: ✖️➕🟰:", ws.SolSum())
	tWs := NewTransposeWorksheet(lines)
	fmt.Println("Part 02: ", tWs.SolSum())

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	split := strings.Split(lines, "\n")
	return split
}
