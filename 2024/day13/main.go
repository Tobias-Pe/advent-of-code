package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

type clawMachine struct {
	btns   buttons
	target coordinate
}

type buttons struct {
	a, b coordinate
}

type coordinate struct {
	x, y float64
}

func (c *clawMachine) invertMatrix() {
	determinant := (c.btns.a.x * c.btns.b.y) - (c.btns.a.y * c.btns.b.x)
	inverseMultiplicator := (&big.Float{}).Quo(big.NewFloat(1), big.NewFloat(float64(determinant)))
	// a           b         c          d       =     d          -b           -c         a
	c.btns.a.x, c.btns.a.y, c.btns.b.x, c.btns.b.y = c.btns.b.y, -c.btns.a.y, -c.btns.b.x, c.btns.a.x
	c.btns.a.x = mulDetereminant(inverseMultiplicator, c.btns.a.x)
	c.btns.a.y = mulDetereminant(inverseMultiplicator, c.btns.a.y)
	c.btns.b.x = mulDetereminant(inverseMultiplicator, c.btns.b.x)
	c.btns.b.y = mulDetereminant(inverseMultiplicator, c.btns.b.y)
}

func (c *clawMachine) calculateButtonPresses() (buttonA, buttonB int) {
	c.invertMatrix()
	a := round(c.target.x*c.btns.a.x+c.btns.b.x*c.target.y, 2)
	b := round(c.target.x*c.btns.a.y+c.btns.b.y*c.target.y, 2)
	if math.Trunc(a) != a {
		return 0, 0
	}
	if math.Trunc(b) != b {
		return 0, 0
	}
	return int(a), int(b)
}

func mulDetereminant(determinant *big.Float, value float64) float64 {
	sol := (&big.Float{}).Mul(determinant, big.NewFloat(value))
	f, _ := sol.Float64()
	return f
}

func round(num float64, decimalPlaces int) float64 {
	multiplier := math.Pow(10, float64(decimalPlaces))
	return math.Round(num*multiplier) / multiplier
}

func main() {
	start := time.Now()

	clawMachines := readFile("day13/input.txt")
	fmt.Println("Part 1:", calcScoreSum(clawMachines, false))
	fmt.Println("Part 2:", calcScoreSum(clawMachines, true))

	fmt.Println("Finished in", time.Since(start))
}

func calcScoreSum(clawMachines []clawMachine, isPart2 bool) (scoreSum int) {
	for _, machine := range clawMachines {
		if isPart2 {
			machine.target.x += 10000000000000
			machine.target.y += 10000000000000
		}
		a, b := machine.calculateButtonPresses()
		score := a*3 + b
		scoreSum += score
		//fmt.Println(i, a, b, score)
	}
	return scoreSum
}

func readFile(file string) []clawMachine {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	inputs := strings.Split(lines, "\n\n")
	clawMachines := make([]clawMachine, len(inputs))
	for i, input := range inputs {
		gameLines := strings.Split(input, "\n")
		clawMachine := clawMachine{}
		clawMachine.btns = buttons{}

		xyStr := strings.Split(strings.TrimLeft(gameLines[0], "Button A: X+"), ", Y+")
		x, _ := strconv.Atoi(xyStr[0])
		y, _ := strconv.Atoi(xyStr[1])
		clawMachine.btns.a = coordinate{float64(x), float64(y)}

		xyStr = strings.Split(strings.TrimLeft(gameLines[1], "Button B: X+"), ", Y+")
		x, _ = strconv.Atoi(xyStr[0])
		y, _ = strconv.Atoi(xyStr[1])
		clawMachine.btns.b = coordinate{float64(x), float64(y)}

		xyStr = strings.Split(strings.TrimLeft(gameLines[2], "Prize: X="), ", Y=")
		x, _ = strconv.Atoi(xyStr[0])
		y, _ = strconv.Atoi(xyStr[1])
		clawMachine.target = coordinate{float64(x), float64(y)}
		clawMachines[i] = clawMachine
	}

	return clawMachines
}
