package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type region struct {
	flowers    []coordinate
	flowerType string
}

type coordinate struct {
	i, j int
}

type gardenArrangement struct {
	garden  [][]string
	regions []region
}

func (ga gardenArrangement) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(ga.garden) && coord.j < len(ga.garden[coord.i])
}

func createGardenArrangement(lines []string) gardenArrangement {
	garden := make([][]string, len(lines))
	for i, line := range lines {
		garden[i] = strings.Split(line, "")
	}
	return gardenArrangement{garden: garden}
}

func main() {
	start := time.Now()

	readFile("day12/input.txt")

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
