package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type wordSearchPuzzle struct {
	puzzle [][]string
}

func createWordSearchPuzzle(lines []string) *wordSearchPuzzle {
	puzzle := make([][]string, len(lines))
	for _, line := range lines {
		characters := strings.Split(line, "")
		puzzle = append(puzzle, characters)
	}
	return &wordSearchPuzzle{puzzle: puzzle}
}

func (p wordSearchPuzzle) sumAllXmasOccurrences(isPart2 bool) int {
	totalAmount := 0
	for i := 0; i < len(p.puzzle); i++ {
		for j := 0; j < len(p.puzzle[i]); j++ {
			amountFound := 0
			if isPart2 {
				amountFound = p.searchXShapedMaxAt(i, j)
			} else {
				amountFound = p.searchXmasAt(i, j)
			}
			totalAmount += amountFound
		}
	}
	return totalAmount
}

func (p wordSearchPuzzle) searchXmasAt(i int, j int) int {
	right := [][]int{{i, j}, {i + 1, j}, {i + 2, j}, {i + 3, j}}
	left := [][]int{{i, j}, {i - 1, j}, {i - 2, j}, {i - 3, j}}
	up := [][]int{{i, j}, {i, j + 1}, {i, j + 2}, {i, j + 3}}
	down := [][]int{{i, j}, {i, j - 1}, {i, j - 2}, {i, j - 3}}
	rightDiagUp := [][]int{{i, j}, {i + 1, j + 1}, {i + 2, j + 2}, {i + 3, j + 3}}
	leftDiagUp := [][]int{{i, j}, {i - 1, j + 1}, {i - 2, j + 2}, {i - 3, j + 3}}
	rightDiagDown := [][]int{{i, j}, {i + 1, j - 1}, {i + 2, j - 2}, {i + 3, j - 3}}
	leftDiagDown := [][]int{{i, j}, {i - 1, j - 1}, {i - 2, j - 2}, {i - 3, j - 3}}
	possibleMovements := [][][]int{right, left, up, down, rightDiagUp, leftDiagUp, rightDiagDown, leftDiagDown}

	foundMatches := 0

	for _, movementPossibilities := range possibleMovements {
		searchWord := "XMAS"
		for _, move := range movementPossibilities {
			if p.isValidCoord(move[0], move[1]) && string(searchWord[0]) == p.puzzle[move[0]][move[1]] {
				searchWord = searchWord[1:]
			} else {
				break
			}
			if len(searchWord) == 0 {
				foundMatches++
				break
			}
		}
	}
	return foundMatches
}

func (p wordSearchPuzzle) searchXShapedMaxAt(i int, j int) int {
	rightDiagTopToBottom := [][]int{{i + 1, j + 1}, {i, j}, {i - 1, j - 1}}
	rightDiagBottomToTop := [][]int{{i - 1, j - 1}, {i, j}, {i + 1, j + 1}}
	rightDiag := [][][]int{rightDiagTopToBottom, rightDiagBottomToTop}

	leftDiagTopToBottom := [][]int{{i - 1, j + 1}, {i, j}, {i + 1, j - 1}}
	leftDiagBottomToTop := [][]int{{i + 1, j - 1}, {i, j}, {i - 1, j + 1}}
	leftDiag := [][][]int{leftDiagTopToBottom, leftDiagBottomToTop}

	foundRight := false
	for _, movementPossibilities := range rightDiag {
		searchWord := "MAS"
		for _, move := range movementPossibilities {
			if p.isValidCoord(move[0], move[1]) && string(searchWord[0]) == p.puzzle[move[0]][move[1]] {
				searchWord = searchWord[1:]
			} else {
				break
			}
			if len(searchWord) == 0 {
				foundRight = true
				break
			}
		}
		if foundRight {
			break
		}
	}

	if !foundRight {
		return 0
	}

	for _, movementPossibilities := range leftDiag {
		searchWord := "MAS"
		for _, move := range movementPossibilities {
			if p.isValidCoord(move[0], move[1]) && string(searchWord[0]) == p.puzzle[move[0]][move[1]] {
				searchWord = searchWord[1:]
			} else {
				break
			}
			if len(searchWord) == 0 {
				return 1
			}
		}
	}
	return 0
}

func (p wordSearchPuzzle) isValidCoord(i, j int) bool {
	return i >= 0 && j >= 0 && i < len(p.puzzle) && j < len(p.puzzle[i])
}

func main() {
	start := time.Now()

	lines := readFile("day4/input.txt")
	wordPuzzle := createWordSearchPuzzle(lines)
	fmt.Println("Part 1:", wordPuzzle.sumAllXmasOccurrences(false))
	fmt.Println("Part 2:", wordPuzzle.sumAllXmasOccurrences(true))

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
