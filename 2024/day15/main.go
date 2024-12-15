package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

func (c coordinate) getGPS() int {
	return c.i + c.j*100
}

type warehouse struct {
	layout       [][]string
	instructions []string
}

func main() {
	start := time.Now()

	readFile("day15/input.txt")

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
