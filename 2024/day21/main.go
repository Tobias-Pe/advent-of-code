package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) sub(other coordinate) coordinate {
	return coordinate{c.i - other.i, c.j - other.j}
}

type padRes struct {
	dest coordinate
	res  []string
}

type padQuery struct {
	origin coordinate
	code   string
}

type pad struct {
	buttons map[string]coordinate
	pos     coordinate
	dp      map[padQuery]padRes
}

func creadNumPad() pad {
	nP := pad{}
	nP.buttons = make(map[string]coordinate)
	nP.buttons["7"] = coordinate{0, 0}
	nP.buttons["8"] = coordinate{0, 1}
	nP.buttons["9"] = coordinate{0, 2}
	nP.buttons["4"] = coordinate{1, 0}
	nP.buttons["5"] = coordinate{1, 1}
	nP.buttons["6"] = coordinate{1, 2}
	nP.buttons["1"] = coordinate{2, 0}
	nP.buttons["2"] = coordinate{2, 1}
	nP.buttons["3"] = coordinate{2, 2}
	nP.buttons[""] = coordinate{3, 0}
	nP.buttons["0"] = coordinate{3, 1}
	nP.buttons["A"] = coordinate{3, 2}
	nP.pos = coordinate{3, 2}
	nP.dp = make(map[padQuery]padRes)
	return nP
}

func createDirPad() pad {
	dP := pad{}
	dP.buttons = make(map[string]coordinate)
	dP.buttons[""] = coordinate{0, 0}
	dP.buttons["^"] = coordinate{0, 1}
	dP.buttons["A"] = coordinate{0, 2}
	dP.buttons["<"] = coordinate{1, 0}
	dP.buttons["v"] = coordinate{1, 1}
	dP.buttons[">"] = coordinate{1, 2}
	dP.pos = coordinate{0, 2}
	dP.dp = make(map[padQuery]padRes)
	return dP
}

type dirCombos struct {
	hori, verti int
	curr        string
}

var dp = map[dirCombos][]string{}

func combinate(hori, verti int, curr string) []string {
	if val, ok := dp[dirCombos{
		hori:  hori,
		verti: verti,
		curr:  curr,
	}]; ok {
		return val
	}
	results := []string{}
	if hori > 0 {
		results = append(results, combinate(hori-1, verti, curr+"-")...)
	}
	if verti > 0 {
		results = append(results, combinate(hori, verti-1, curr+"|")...)
	}
	if hori == 0 && verti == 0 {
		results = append(results, curr)
	}
	dp[dirCombos{
		hori:  hori,
		verti: verti,
		curr:  curr,
	}] = results
	return results
}

func (p *pad) goTo(dest string) []string {
	if val, ok := p.dp[padQuery{
		origin: p.pos,
		code:   dest,
	}]; ok {
		p.pos = val.dest
		return val.res
	}

	if dest == "" {
		panic("Not allowed empty string")
	}
	oldPos := p.pos
	diff := p.buttons[dest].sub(oldPos)
	p.pos = p.buttons[dest]

	rawCombinations := combinate(int(math.Abs(float64(diff.j))), int(math.Abs(float64(diff.i))), "")
	combinations := make([]string, len(rawCombinations))
	for i, combination := range rawCombinations {
		combinations[i] = combination
		if diff.i > 0 {
			combinations[i] = strings.ReplaceAll(combinations[i], "|", "v")
		} else {
			combinations[i] = strings.ReplaceAll(combinations[i], "|", "^")
		}
		if diff.j > 0 {
			combinations[i] = strings.ReplaceAll(combinations[i], "-", ">")
		} else {
			combinations[i] = strings.ReplaceAll(combinations[i], "-", "<")
		}
	}
	p.dp[padQuery{
		origin: oldPos,
		code:   dest,
	}] = padRes{
		dest: p.pos,
		res:  combinations,
	}
	return combinations
}

func (p pad) execute(codes map[string]bool) map[string]bool {
	results := map[string]bool{}
	for code, _ := range codes {
		for _, result := range p.executeRecurse(code, "") {
			results[result] = true
		}
	}
	return results
}

func (p pad) executeRecurse(code string, curr string) []string {
	results := []string{}
	if len(code) == 0 {
		results = append(results, curr)
		return results
	}

	res := p.goTo(string(code[0]))
	for i := 0; i < len(res); i++ {
		results = append(results, p.executeRecurse(code[1:], curr+res[i]+"A")...)
	}
	return results
}

type spaceShip struct {
	numPad pad
	robot1 pad
	robot2 pad
	me     pad
	result string
	code   string
}

func createSpaceShip(code string) spaceShip {
	spcShip := spaceShip{}

	spcShip.numPad = creadNumPad()
	spcShip.robot1 = createDirPad()
	spcShip.robot2 = createDirPad()
	spcShip.me = createDirPad()

	spcShip.code = code

	return spcShip
}

func (s *spaceShip) execute() {
	initCode := make(map[string]bool)
	initCode[s.code] = true
	resultNumPad := s.numPad.execute(initCode)
	resultRobot1 := s.robot1.execute(resultNumPad)
	resultRobot2 := s.robot2.execute(resultRobot1)
	resultMe := s.me.execute(resultRobot2)
	minRes := ""
	for s2, _ := range resultMe {
		if minRes == "" || len(s2) < len(minRes) {
			minRes = s2
		}
	}
	s.result = minRes
}

func main() {
	start := time.Now()

	codes := readFile("day21/input.txt")
	scoreSum := 0
	for i, code := range codes {
		ship := createSpaceShip(code)
		ship.execute()
		numPart, err := strconv.Atoi(strings.Trim(code, "A"))
		if err != nil {
			return
		}
		score := numPart * len(ship.result)
		fmt.Println(i, code, numPart, "*", len(ship.result), ship.result)
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
