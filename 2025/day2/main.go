package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	start, end int64
}

func NewRange(str string) *Range {
	str = strings.TrimSpace(str)
	args := strings.Split(str, "-")
	atoiStart, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		panic(err)
	}
	atoiEnd, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return &Range{atoiStart, atoiEnd}
}

func (r *Range) InvalidIDSumP1() int64 {
	var sum int64

	for i := r.start; i <= r.end; i++ {
		if r.isInvalidP1(i) {
			sum += i
		}
	}

	return sum
}

func (r *Range) isInvalidP1(num int64) bool {
	strNum := strconv.FormatInt(num, 10)
	left := strNum[0 : len(strNum)/2]
	right := strNum[len(strNum)/2:]
	if len(left) != len(right) { // useless speed up 🏎️
		return false
	}
	return strings.Compare(left, right) == 0
}

func (r *Range) InvalidIDSumP2() int64 {
	var sum int64
	for i := r.start; i <= r.end; i++ {
		if r.isInvalidP2(i) {
			sum += i
		}
	}
	return sum
}

func (r *Range) isInvalidP2(num int64) bool {
	strNum := strconv.FormatInt(num, 10)

	currStr := ""
	for i := 0; i < len(strNum); i++ {
		currStr += string(strNum[i])
		if fullMatch, _ := regexp.MatchString("^("+currStr+")("+currStr+")+$", strNum); fullMatch { // truly devilish unhinged regex magic 🔮🪄
			return true
		}
	}
	return false
}

func main() {
	start := time.Now()

	lines := readFile("day2/input.txt")
	//assume we have only one line
	line := lines[0]
	strs := strings.Split(line, ",")
	ranges := make([]*Range, len(strs))
	for i, str := range strs {
		ranges[i] = NewRange(str)
	}
	sum := int64(0)
	for _, rng := range ranges {
		sum += rng.InvalidIDSumP1()
	}
	fmt.Println("Part01 Sum:", sum)

	sum = int64(0)
	for _, rng := range ranges {
		sum += rng.InvalidIDSumP2()
	}
	fmt.Println("Part02 Sum:", sum)

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
