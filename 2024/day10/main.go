package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

type hikingGuide struct {
	topoMap    [][]int
	trailheads []coordinate
}

func (hG hikingGuide) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(hG.topoMap) && coord.j < len(hG.topoMap[coord.i])
}

func createHikingGuide(lines []string) hikingGuide {
	hG := hikingGuide{topoMap: make([][]int, len(lines)), trailheads: []coordinate{}}
	for i, line := range lines {
		mapRow := make([]int, len(line))
		for j, heightString := range line {
			height, _ := strconv.Atoi(string(heightString))
			mapRow[j] = height
			if height == 0 {
				hG.trailheads = append(hG.trailheads, coordinate{i, j})
			}
		}
		hG.topoMap[i] = mapRow
	}
	return hG
}

func (hG hikingGuide) getScore(trailHead coordinate) int {
	score := 0
	visitedNines := make(map[coordinate]bool)
	nodeQueue := []coordinate{trailHead}
	for len(nodeQueue) > 0 {
		currNode := nodeQueue[0]
		if hG.topoMap[currNode.i][currNode.j] == 9 && !visitedNines[currNode] {
			score++
			visitedNines[currNode] = true
		}
		nodeQueue = nodeQueue[1:]
		neighbours := currNode.getNeighbours()
		for _, neighbour := range neighbours {
			if hG.isValid(neighbour) && hG.topoMap[currNode.i][currNode.j]+1 == hG.topoMap[neighbour.i][neighbour.j] {
				nodeQueue = append(nodeQueue, neighbour)
			}
		}
	}

	return score
}

func (hG hikingGuide) getRating(trailHead coordinate) int {
	rating := 0
	nodeQueue := []coordinate{trailHead}
	for len(nodeQueue) > 0 {
		currNode := nodeQueue[0]
		if hG.topoMap[currNode.i][currNode.j] == 9 {
			rating++
		}
		nodeQueue = nodeQueue[1:]
		neighbours := currNode.getNeighbours()
		for _, neighbour := range neighbours {
			if hG.isValid(neighbour) && hG.topoMap[currNode.i][currNode.j]+1 == hG.topoMap[neighbour.i][neighbour.j] {
				nodeQueue = append(nodeQueue, neighbour)
			}
		}
	}

	return rating
}

func (hG hikingGuide) getScores() int {
	sum := 0
	for _, trailhead := range hG.trailheads {
		sum += hG.getScore(trailhead)
	}
	return sum
}

func (hG hikingGuide) getRatings() int {
	rating := 0
	for _, trailhead := range hG.trailheads {
		rating += hG.getRating(trailhead)
	}
	return rating
}

func (hG hikingGuide) String() string {
	sB := strings.Builder{}
	for _, heights := range hG.topoMap {
		for _, height := range heights {
			sB.WriteString(strconv.Itoa(height))
		}
		sB.WriteString("\n")
	}
	sB.WriteString("\n")
	for _, trailhead := range hG.trailheads {
		sB.WriteString(trailhead.String())
		sB.WriteString("\n")
	}
	return sB.String()
}

func main() {
	start := time.Now()

	input := readFile("day10/input.txt")
	hG := createHikingGuide(input)
	fmt.Println("Part 1:", hG.getScores())
	fmt.Println("Part 2:", hG.getRatings())

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
