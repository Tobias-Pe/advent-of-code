package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	lines := readFile("day1/input.txt")

	distance, err := sumOfDistances(lines)
	fmt.Println("Part1 sumOfDistances Result:", distance, err)

	similarityScore, err := sumOfSimilarityScores(lines)
	fmt.Println("Part2 sumOfSimilarityScores Result:", similarityScore, err)

	fmt.Println("Finished in", time.Since(start))
}

func sumOfDistances(lines []string) (int, error) {
	var leftNums []int
	var rightNums []int
	for _, line := range lines {
		nums := strings.Fields(line)
		leftNum, err := strconv.Atoi(nums[0])
		if err != nil {
			return 0, err
		}
		rightNum, err := strconv.Atoi(nums[1])
		if err != nil {
			return 0, err
		}
		leftNums = append(leftNums, leftNum)
		rightNums = append(rightNums, rightNum)
	}
	slices.Sort(leftNums)
	slices.Sort(rightNums)
	if len(leftNums) != len(rightNums) {
		return 0, fmt.Errorf("left and right numbers length not equal")
	}

	sumDistances := 0

	for i := 0; i < len(leftNums); i++ {
		distance := leftNums[i] - rightNums[i]
		sumDistances += int(math.Abs(float64(distance)))
	}

	return sumDistances, nil
}

func sumOfSimilarityScores(lines []string) (int, error) {
	var leftNums []int
	rightNumsCounter := make(map[int]int)
	for _, line := range lines {
		nums := strings.Fields(line)
		leftNum, err := strconv.Atoi(nums[0])
		if err != nil {
			return 0, err
		}
		rightNum, err := strconv.Atoi(nums[1])
		if err != nil {
			return 0, err
		}
		leftNums = append(leftNums, leftNum)
		rightNumsCounter[rightNum] = rightNumsCounter[rightNum] + 1
	}
	sumSimilarityScores := 0

	for i := 0; i < len(leftNums); i++ {
		similarityScore := leftNums[i] * rightNumsCounter[leftNums[i]]
		sumSimilarityScores += int(math.Abs(float64(similarityScore)))
	}

	return sumSimilarityScores, nil
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
