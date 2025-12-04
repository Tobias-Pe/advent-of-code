package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Tobias-Pe/advent-of-code/util"
)

type PaperRollDiagram struct {
	diagram    [][]string
	iMax, jMax int
}

func NewPaperRollDiagram(lines []string) *PaperRollDiagram {
	diagram := make([][]string, len(lines))
	for i, line := range lines {
		diagram[i] = make([]string, len(line))
		for j, cell := range strings.Split(line, "") {
			diagram[i][j] = cell
		}
	}
	return &PaperRollDiagram{
		diagram: diagram,
		iMax:    len(lines),
		jMax:    len(lines[0]),
	}
}

func (p *PaperRollDiagram) IsAccessibleRoll(i, j int) bool {
	if !p.IsPaperRoll(i, j) {
		return false
	}
	nNPR := p.CountNeighbourPaperRolls(i, j)
	if nNPR < 4 {
		return true
	}
	return false
}

func (p *PaperRollDiagram) IsPaperRoll(i, j int) bool {
	return p.diagram[i][j] == "@"
}

func (p *PaperRollDiagram) CountAccessibleRolls() int {
	counter := 0
	for i, row := range p.diagram {
		for j := range row {
			if p.IsAccessibleRoll(i, j) {
				counter++
			}
		}
	}
	return counter
}

func (p *PaperRollDiagram) CountNeighbourPaperRolls(i, j int) int {
	coord := util.Coordinate{I: i, J: j}
	neighbours := coord.GetNeighbours8()
	var validNeighbours []util.Coordinate
	for _, neighbour := range neighbours {
		if !neighbour.IsValid(p.iMax, p.jMax) {
			continue
		}
		if p.IsPaperRoll(neighbour.I, neighbour.J) {
			validNeighbours = append(validNeighbours, neighbour)
		}
	}
	return len(validNeighbours)
}

func (p *PaperRollDiagram) RemoveAccessibleRolls() (removed int) {
	newDiagram := make([][]string, len(p.diagram))
	for i, row := range p.diagram {
		newDiagram[i] = make([]string, len(row))
		for j, cell := range row {
			newDiagram[i][j] = cell
			if p.IsAccessibleRoll(i, j) {
				newDiagram[i][j] = "."
				removed++
			}
		}
	}
	p.diagram = newDiagram
	return removed
}

func (p *PaperRollDiagram) RemoveAllAccessibleRolls() (total int) {
	removed := p.RemoveAccessibleRolls()
	for removed != 0 {
		total += removed
		removed = p.RemoveAccessibleRolls()
	}
	return total
}

func main() {
	start := time.Now()

	lines := readFile("day4/input.txt")
	prd := NewPaperRollDiagram(lines)
	fmt.Println("Part 01:", prd.CountAccessibleRolls())
	fmt.Println("Part 02:", prd.RemoveAllAccessibleRolls())

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
