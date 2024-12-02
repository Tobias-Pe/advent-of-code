package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type report struct {
	levels []int
}

func isSafe(levels []int) bool {
	if len(levels) <= 1 {
		return true
	}
	var lastDistance int
	for i := 1; i < len(levels); i++ {
		distance := levels[i] - levels[i-1]
		if i > 1 && math.Signbit(float64(distance)) != math.Signbit(float64(lastDistance)) {
			return false
		}
		if distance == 0 || math.Abs(float64(distance)) > 3 {
			return false
		}
		lastDistance = distance
	}
	return true
}

func isSafeWithRemovalOnce(levels []int) bool {
	if isSafe(levels) {
		return true
	}
	for i := 0; i < len(levels); i++ {
		tmp := append([]int{}, levels[:i]...)
		tmp = append(tmp, levels[i+1:]...)
		if isSafe(tmp) {
			return true
		}
	}
	return false
}

func main() {
	start := time.Now()

	lines := readFile("day2/input.txt")
	reports, err := parseReports(lines)
	if err != nil {
		panic(err)
	}

	sumSafeReports(reports)
	sumSafeReportsWithRetry(reports)

	fmt.Println("Finished in", time.Since(start))
}

func sumSafeReports(reports []report) {
	sumSafeReports := 0
	for i, report := range reports {
		isSafe := isSafe(report.levels)
		fmt.Println(i, report, isSafe, sumSafeReports)
		if isSafe {
			sumSafeReports++
		}
	}
	fmt.Println("Part 1:", sumSafeReports)
}

func sumSafeReportsWithRetry(reports []report) {
	sumSafeReports := 0
	for i, report := range reports {
		isSafe := isSafeWithRemovalOnce(report.levels)
		fmt.Println(i, report, isSafe, sumSafeReports)
		if isSafe {
			sumSafeReports++
		}
	}
	fmt.Println("Part 2:", sumSafeReports)
}

func parseReports(lines []string) ([]report, error) {
	reports := []report{}

	for _, line := range lines {
		fields := strings.Fields(line)
		levels := []int{}
		for _, field := range fields {
			level, err := strconv.Atoi(field)
			if err != nil {
				return reports, err
			}
			levels = append(levels, level)
		}
		reports = append(reports, report{levels: levels})
	}

	return reports, nil
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
