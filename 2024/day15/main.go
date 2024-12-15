package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) getNextPos(instruction string) coordinate {
	switch instruction {
	case "^":
		return coordinate{c.i - 1, c.j}
	case "v":
		return coordinate{c.i + 1, c.j}
	case ">":
		return coordinate{c.i, c.j + 1}
	case "<":
		return coordinate{c.i, c.j - 1}
	default:
		return c
	}
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

func (c coordinate) getGPS() int {
	return c.i*100 + c.j
}

type warehouse struct {
	layout       [][]string
	instructions []string
	robot        coordinate
}

func (w warehouse) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(w.layout) && coord.j < len(w.layout[coord.i])
}

func (w *warehouse) executeInstruction() {
	w.move(w.robot, w.instructions[0])
	w.instructions = w.instructions[1:]
}

func (w *warehouse) executeInstructions(out io.Writer) {
	for len(w.instructions) > 0 {
		fmt.Fprintln(out, w)
		w.executeInstruction()
	}
	fmt.Fprintln(out, w)
}

func (w warehouse) findWholeBox(c coordinate) (coordinate, coordinate) {
	if !w.isValid(c) {
		panic("Invalid coordinate")
	}
	if w.layout[c.i][c.j] == "[" {
		return c, coordinate{c.i, c.j + 1}
	}
	if w.layout[c.i][c.j] != "]" {
		panic("Tried to find whole box without using a box coordinate")
	}
	return coordinate{c.i, c.j - 1}, c
}

func (w warehouse) findOtherBox(c coordinate) coordinate {
	if !w.isValid(c) {
		panic("Invalid coordinate")
	}
	if w.layout[c.i][c.j] == "[" {
		return coordinate{c.i, c.j + 1}
	}
	if w.layout[c.i][c.j] != "]" {
		panic("Tried to find whole box without using a box coordinate")
	}
	return coordinate{c.i, c.j - 1}
}

func (w warehouse) copy() ([][]string, coordinate) {
	dupe := make([][]string, len(w.layout))
	for i := range w.layout {
		dupe[i] = make([]string, len(w.layout[i]))
		copy(dupe[i], w.layout[i])
	}
	return dupe, coordinate{w.robot.i, w.robot.j}
}

func (w *warehouse) move(c coordinate, instruction string) bool {
	nextPos := c.getNextPos(instruction)
	if !w.isValid(nextPos) {
		return false
	}
	if w.layout[nextPos.i][nextPos.j] == "#" {
		return false
	}
	if w.layout[nextPos.i][nextPos.j] == "O" {
		w.move(nextPos, instruction)
	} else if w.layout[nextPos.i][nextPos.j] == "[" || w.layout[nextPos.i][nextPos.j] == "]" {
		other := w.findOtherBox(nextPos)
		layoutCopy, robotCopy := w.copy()
		if !w.move(other, instruction) {
			return false
		}
		if !w.move(nextPos, instruction) {
			w.layout = layoutCopy
			w.robot = robotCopy
			return false
		}
	}
	if w.layout[nextPos.i][nextPos.j] == "." {
		w.layout[nextPos.i][nextPos.j] = w.layout[c.i][c.j]
		w.layout[c.i][c.j] = "."
		if c.i == w.robot.i && c.j == w.robot.j {
			w.robot = coordinate{nextPos.i, nextPos.j}
		}
		return true
	}
	return false
}

func (w warehouse) String() string {
	sb := strings.Builder{}
	for i, line := range w.layout {
		for j, cell := range line {
			if w.robot.i == i && w.robot.j == j {
				sb.WriteString("🤖")
			} else {
				if cell == "#" {
					sb.WriteString("🧱")
				} else if cell == "O" {
					sb.WriteString("📦")
				} else if cell == "[" {
					sb.WriteString("🫷")
				} else if cell == "]" {
					sb.WriteString("🫸")
				} else {
					sb.WriteString("⬜")
				}
			}
		}
		sb.WriteString("\n")
	}
	if w.instructions != nil && len(w.instructions) > 0 {
		sb.WriteString(fmt.Sprintf("Next move: %s\n", w.instructions[0]))
	}
	return sb.String()
}

func (w warehouse) sumGPS() int {
	sum := 0
	for i, line := range w.layout {
		for j, cell := range line {
			if cell == "O" || cell == "[" {
				sum += coordinate{i, j}.getGPS()
			}
		}
	}
	return sum
}

func main() {
	start := time.Now()

	wh := readFile("day15/input.txt")
	wh.sumGPS()
	wh.executeInstructions(io.Discard)
	fmt.Println("Part 1:", wh.sumGPS())
	whP2 := readFileP2("day15/input.txt")
	whP2.executeInstructions(io.Discard)
	fmt.Println("Part 2:", whP2.sumGPS())

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) warehouse {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	inputs := strings.Split(lines, "\n\n")
	warehouseLines := strings.Split(inputs[0], "\n")
	warehouse := warehouse{
		layout:       make([][]string, len(warehouseLines)),
		instructions: []string{},
	}
	for i, line := range warehouseLines {
		tileRow := make([]string, len(line))
		for j, cell := range line {
			if cell == '@' {
				warehouse.robot = coordinate{i, j}
				tileRow[j] = "."
			} else {
				tileRow[j] = string(cell)
			}
		}
		warehouse.layout[i] = tileRow
	}
	instructionLines := strings.Split(inputs[1], "\n")
	for _, line := range instructionLines {
		for _, cell := range line {
			warehouse.instructions = append(warehouse.instructions, string(cell))
		}
	}
	return warehouse
}

func readFileP2(file string) warehouse {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	inputs := strings.Split(lines, "\n\n")
	warehouseLines := strings.Split(inputs[0], "\n")
	warehouse := warehouse{
		layout:       make([][]string, len(warehouseLines)),
		instructions: []string{},
	}
	for i, line := range warehouseLines {
		tileRow := []string{}
		for _, cell := range line {
			if cell == '@' {
				warehouse.robot = coordinate{i, len(tileRow)}
				tileRow = append(tileRow, ".")
				tileRow = append(tileRow, ".")
			} else if cell == 'O' {
				tileRow = append(tileRow, "[")
				tileRow = append(tileRow, "]")
			} else {
				tileRow = append(tileRow, string(cell))
				tileRow = append(tileRow, string(cell))
			}
		}
		warehouse.layout[i] = tileRow
	}
	instructionLines := strings.Split(inputs[1], "\n")
	for _, line := range instructionLines {
		for _, cell := range line {
			warehouse.instructions = append(warehouse.instructions, string(cell))
		}
	}
	return warehouse
}
