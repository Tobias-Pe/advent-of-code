package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
)

var grid [][]*square

var start *square
var end *square

type square struct {
	letter              rune
	elevation           int
	position            position
	neighbours          []*square
	visitableNeighbours []*square
	stepsToReach        int
}

type position struct {
	i int
	j int
}

func (s *square) populateElevationLevel() {
	letter := s.letter
	if unicode.IsUpper(letter) {
		if letter == 'S' {
			s.elevation = 0
			s.stepsToReach = 0
			start = s
			fmt.Println("Start found!")
		} else {
			s.elevation = 25
			end = s
			fmt.Println("End found!")
		}
	} else {
		ascii := int(letter)
		ascii -= 97
		s.elevation = ascii
	}
}

func (s *square) getNeighbours() []*square {
	var neighbours []*square
	if s.position.i > 0 {
		neighbours = append(neighbours, grid[s.position.i-1][s.position.j])
	}
	if s.position.i < len(grid)-1 {
		neighbours = append(neighbours, grid[s.position.i+1][s.position.j])
	}
	if s.position.j > 0 {
		neighbours = append(neighbours, grid[s.position.i][s.position.j-1])
	}
	if s.position.j < len(grid[s.position.i])-1 {
		neighbours = append(neighbours, grid[s.position.i][s.position.j+1])
	}
	return neighbours
}

func (s *square) isVisitable(neighbour square) bool {
	return neighbour.elevation <= s.elevation+1
}

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	splittedLines := strings.Split(lines, "\n")
	grid = make([][]*square, len(splittedLines))
	for i, line := range splittedLines {
		grid[i] = make([]*square, len(line))
		for j, letter := range line {
			s := square{
				letter: letter,
				position: position{
					i: i,
					j: j,
				},
				stepsToReach: -1,
			}
			s.populateElevationLevel()
			grid[i][j] = &s
		}
	}
}

func populateNeighbours() {
	for _, squares := range grid {
		for _, s := range squares {
			s.neighbours = s.getNeighbours()
			for _, neighbour := range s.neighbours {
				if s.isVisitable(*neighbour) {
					s.visitableNeighbours = append(s.visitableNeighbours, neighbour)
				}
			}
		}
	}
}

func populateShortestPathFrom(start *square) {
	var visited = map[position]bool{}
	var toBeVisited []*square
	toBeVisited = append(toBeVisited, start)
	for len(toBeVisited) > 0 {
		visiting := toBeVisited[0]
		if isVisited, ok := visited[visiting.position]; !ok || !isVisited {
			visited[visiting.position] = true

			for _, neighbour := range visiting.visitableNeighbours {
				isVisited, ok := visited[neighbour.position]
				if !ok || !isVisited {
					toBeVisited = append(toBeVisited, neighbour)
					if neighbour.stepsToReach == -1 || visiting.stepsToReach+1 < neighbour.stepsToReach {
						neighbour.stepsToReach = visiting.stepsToReach + 1
					}
				}
			}
		}

		toBeVisited = append(toBeVisited[:0], toBeVisited[1:]...)
	}
	fmt.Println(end.stepsToReach)
}

func resetGrid() {
	for _, squares := range grid {
		for _, s := range squares {
			s.stepsToReach = -1
		}
	}
}

func findBestShortestPathFrom(possibleStarts []*square) {
	bestEnd := -1
	for _, possibleStart := range possibleStarts {
		resetGrid()
		possibleStart.stepsToReach = 0
		populateShortestPathFrom(possibleStart)
		if bestEnd == -1 || end.stepsToReach < bestEnd && end.stepsToReach != -1 {
			bestEnd = end.stepsToReach
		}
	}
	fmt.Println(bestEnd)
}

func getPossibleStarts() []*square {
	var possibleStarts []*square

	possibleStarts = append(possibleStarts, start)

	for _, squares := range grid {
		for _, s := range squares {
			if s.elevation == 0 {
				possibleStarts = append(possibleStarts, s)
			}
		}
	}

	return possibleStarts
}

func main() {
	readFile("input.txt")
	populateNeighbours()
	populateShortestPathFrom(start)
	fmt.Println(end.stepsToReach, "\n")

	findBestShortestPathFrom(getPossibleStarts())
}
