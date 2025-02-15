package main

import (
	"fmt"
	"github.com/Tobias-Pe/advent-of-code/util"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type query struct {
	from string
	to   string
}

type pad struct {
	buttons      [][]string
	buttonsMap   map[string]util.Coordinate
	combinations map[query][]string
}

func (p *pad) calcCombinations() {
	p.buttonsMap = map[string]util.Coordinate{}
	buttons := []string{}
	for i, row := range p.buttons {
		for j, cell := range row {
			if p.buttons[i][j] == "" {
				continue
			}
			p.buttonsMap[cell] = util.Coordinate{I: i, J: j}
			buttons = append(buttons, cell)
		}
	}
	p.combinations = make(map[query][]string)
	for _, buttonFrom := range buttons {
		for _, buttonTo := range buttons {
			q := query{buttonFrom, buttonTo}
			if buttonFrom == buttonTo {
				p.makeCombEntry(q, "")
				continue
			}

			possiblePaths := p.dfs(q)
			for _, possiblePath := range possiblePaths {
				p.makeCombEntry(q, possiblePath)
			}
		}
	}
}

func (p *pad) makeCombEntry(q query, s string) {
	if _, exists := p.combinations[q]; !exists {
		p.combinations[q] = []string{}
	}
	p.combinations[q] = append(p.combinations[q], s)
}

type dfsResult struct {
	way     string
	currPos util.Coordinate
}

func (p *pad) dfs(q query) []string {
	queue := []dfsResult{{
		way:     "",
		currPos: p.buttonsMap[q.from],
	}}
	visited := map[util.Coordinate]bool{}
	var sols []dfsResult
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.currPos == p.buttonsMap[q.to] {
			sols = append(sols, curr)
			continue
		}

		if p.buttons[curr.currPos.I][curr.currPos.J] == "" {
			continue
		}

		if _, ok := visited[curr.currPos.Left()]; !ok && curr.currPos.Left().IsValid(len(p.buttons), len(p.buttons[0])) {
			queue = append(queue, dfsResult{curr.way + "<", curr.currPos.Left()})
		}
		if _, ok := visited[curr.currPos.Right()]; !ok && curr.currPos.Right().IsValid(len(p.buttons), len(p.buttons[0])) {
			queue = append(queue, dfsResult{curr.way + ">", curr.currPos.Right()})
		}
		if _, ok := visited[curr.currPos.Up()]; !ok && curr.currPos.Up().IsValid(len(p.buttons), len(p.buttons[0])) {
			queue = append(queue, dfsResult{curr.way + "^", curr.currPos.Up()})
		}
		if _, ok := visited[curr.currPos.Down()]; !ok && curr.currPos.Down().IsValid(len(p.buttons), len(p.buttons[0])) {
			queue = append(queue, dfsResult{curr.way + "v", curr.currPos.Down()})
		}
		visited[curr.currPos] = true
	}
	minSolLen := math.MaxInt32
	for _, sol := range sols {
		minSolLen = min(minSolLen, len(sol.way))
	}
	var result []string
	for _, sol := range sols {
		if len(sol.way) > minSolLen {
			continue
		}
		result = append(result, sol.way)
	}
	return result
}

func (p *pad) Execute(code string) []string {
	var sols []string
	curr := "A"
	for _, i32Letter := range code {
		letter := string(i32Letter)
		combs := p.combinations[query{curr, letter}]
		if len(sols) == 0 {
			sols = append([]string{}, combs...)
		} else {
			var tmp []string
			for _, sol := range sols {
				for _, comb := range combs {
					tmp = append(tmp, sol+comb)
				}
			}
			sols = tmp
		}
		for i := range sols {
			sols[i] = sols[i] + "A"
		}
		curr = letter
	}
	return sols
}

func createNumPad() pad {
	p := pad{}
	p.buttons = [][]string{
		{"7", "8", "9"},
		{"4", "5", "6"},
		{"1", "2", "3"},
		{"", "0", "A"},
	}
	p.calcCombinations()
	return p
}

func createDirPad() pad {
	p := pad{}
	p.buttons = [][]string{
		{"", "^", "A"},
		{"<", "v", ">"},
	}
	p.calcCombinations()
	return p
}

type spaceShip struct {
	numPad pad
	robot1 pad
	robot2 pad
	code   string
}

func createSpaceShip(code string) spaceShip {
	spcShip := spaceShip{}

	spcShip.numPad = createNumPad()
	spcShip.robot1 = createDirPad()
	spcShip.robot2 = createDirPad()

	spcShip.code = code

	return spcShip
}

func (s *spaceShip) executeOn(codes []string, target *pad) []string {
	var resultCodes []string
	for _, neededInstruction := range codes {
		res := target.Execute(neededInstruction)
		resultCodes = append(resultCodes, res...)
	}
	minLen := math.MaxInt32
	for _, instructionCode := range resultCodes {
		minLen = min(minLen, len(instructionCode))
	}
	var minResultCodes []string
	for _, instructionCode := range resultCodes {
		if len(instructionCode) == minLen {
			minResultCodes = append(minResultCodes, instructionCode)
		}
	}
	return minResultCodes
}

func (s *spaceShip) execute() []string {
	instructionCodes := s.executeOn([]string{s.code}, &s.numPad)
	instructionCodes = s.executeOn(instructionCodes, &s.robot1)
	instructionCodes = s.executeOn(instructionCodes, &s.robot2)
	return instructionCodes
}

func (s *spaceShip) executeAndPickRandom() string {
	codes := s.execute()
	return codes[0]
}

func main() {
	start := time.Now()

	codes := readFile("day21/input.txt")
	scoreSum := 0
	for i, code := range codes {
		ship := createSpaceShip(code)
		result := ship.executeAndPickRandom()
		numPart, err := strconv.Atoi(strings.Trim(code, "A"))
		if err != nil {
			return
		}
		score := numPart * len(result)
		fmt.Println(i, code, numPart, "*", len(result), result)
		scoreSum += score
	}
	fmt.Println("Part 1:", scoreSum)

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
