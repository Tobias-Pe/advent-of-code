package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Tobias-Pe/advent-of-code/util"
)

type TachyonManifolds struct {
	tachyonBeam map[util.Coordinate]bool
	diagram     [][]string
	splitter    map[util.Coordinate]bool
	start       util.Coordinate
	splitCount  int
	dp          map[util.Coordinate]int
}

func NewTachyonManifolds(lines []string) *TachyonManifolds {
	tm := TachyonManifolds{}
	tm.tachyonBeam = make(map[util.Coordinate]bool)
	tm.splitter = make(map[util.Coordinate]bool)
	tm.dp = make(map[util.Coordinate]int)
	tm.diagram = make([][]string, len(lines))
	for i, line := range lines {
		tm.diagram[i] = strings.Split(line, "")
		for j := range tm.diagram[i] {
			if tm.diagram[i][j] == "S" {
				tm.start = util.Coordinate{I: i, J: j}
			} else if tm.diagram[i][j] == "^" {
				tm.splitter[util.Coordinate{I: i, J: j}] = true
			}
		}
	}
	return &tm
}

func (tm *TachyonManifolds) RunBeam(pos util.Coordinate) {
	if tm.tachyonBeam[pos] {
		return
	}
	if tm.splitter[pos] {
		tm.splitCount++
		tm.RunBeam(pos.Left())
		tm.RunBeam(pos.Right())
		return
	}
	if !pos.IsValid(len(tm.diagram), len(tm.diagram[0])) {
		return
	}
	tm.tachyonBeam[pos] = true
	tm.RunBeam(pos.Down())
}

func (tm *TachyonManifolds) Print() {
	out := ""
	for i, line := range tm.diagram {
		for j := range line {
			if tm.tachyonBeam[util.Coordinate{I: i, J: j}] {
				out += "|"
			} else {
				out += tm.diagram[i][j]
			}
		}
		out += "\n"
	}
	fmt.Println(out)
}

func (tm *TachyonManifolds) PossibilityCount(pos util.Coordinate) int {
	if c, ok := tm.dp[pos]; ok {
		return c
	}
	if tm.splitter[pos] {
		posCount := tm.PossibilityCount(pos.Left()) + tm.PossibilityCount(pos.Right())
		tm.dp[pos] = posCount
		return posCount
	}
	if !pos.IsValid(len(tm.diagram), len(tm.diagram[0])) {
		return 1
	}
	return tm.PossibilityCount(pos.Down())
}

func main() {
	start := time.Now()

	lines := readFile("day7/input.txt")
	tm := NewTachyonManifolds(lines)
	tm.RunBeam(tm.start)
	fmt.Println("Part 01:", tm.splitCount)
	fmt.Println("Part 02:", tm.PossibilityCount(tm.start))

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
