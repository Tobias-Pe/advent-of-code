package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type calibrationEquation struct {
	result  int64
	numbers []int64
}

func (cE calibrationEquation) isValid(opsCount int) bool {
	if opsCount == 2 {
		return cE.recurseSolveTwoOps(cE.numbers)
	}
	return cE.recurseSolveThreeOps(cE.numbers)
}

func (cE calibrationEquation) recurseSolveTwoOps(queue []int64) bool {
	if len(queue) == 1 {
		return cE.result == queue[0]
	}

	return cE.recurseSolveTwoOps(append([]int64{queue[0] + queue[1]}, queue[2:]...)) ||
		cE.recurseSolveTwoOps(append([]int64{queue[0] * queue[1]}, queue[2:]...))
}

func (cE calibrationEquation) recurseSolveThreeOps(queue []int64) bool {
	if len(queue) == 1 {
		return cE.result == queue[0]
	}

	num1AsString := strconv.FormatInt(queue[0], 10)
	num2AsString := strconv.FormatInt(queue[1], 10)
	concatNum, _ := strconv.ParseInt(num1AsString+num2AsString, 10, 64)
	return cE.recurseSolveThreeOps(append([]int64{queue[0] + queue[1]}, queue[2:]...)) ||
		cE.recurseSolveThreeOps(append([]int64{queue[0] * queue[1]}, queue[2:]...)) ||
		cE.recurseSolveThreeOps(append([]int64{concatNum}, queue[2:]...))
}

func main() {
	start := time.Now()

	input := readFile("day7/input.txt")
	calibrationEquations := parseCalibrationEquations(input)

	fmt.Println("Part 1:", sumValidEquations(2, calibrationEquations))
	fmt.Println("Part 2:", sumValidEquations(3, calibrationEquations))

	fmt.Println("Finished in", time.Since(start))
}

func sumValidEquations(opCount int, equations []calibrationEquation) int64 {
	sum := int64(0)
	for _, equation := range equations {
		if equation.isValid(opCount) {
			sum += equation.result
		}
	}
	return sum
}

func parseCalibrationEquations(input []string) []calibrationEquation {
	var calibrationEquations []calibrationEquation
	for _, line := range input {
		cE := calibrationEquation{}
		args := strings.Split(line, ":")
		result, _ := strconv.Atoi(args[0])
		cE.result = int64(result)
		fields := strings.Fields(strings.TrimSpace(args[1]))
		var nums []int64
		for _, field := range fields {
			num, _ := strconv.Atoi(field)
			nums = append(nums, int64(num))
		}
		cE.numbers = nums
		calibrationEquations = append(calibrationEquations, cE)
	}
	return calibrationEquations
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
