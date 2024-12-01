package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type sectorRange struct {
	start int
	end   int
}

type elfPair struct {
	input             string
	firstSectorRange  sectorRange
	secondSectorRange sectorRange
}

func isSectorRangePartOfOtherSectorRange(firstSectorRange sectorRange, secondSectorRange sectorRange) bool {
	return (firstSectorRange.start <= secondSectorRange.start && firstSectorRange.end >= secondSectorRange.end) ||
		(secondSectorRange.start <= firstSectorRange.start && secondSectorRange.end >= firstSectorRange.end)
}

func isSectorOverlappingWithSectorRange(firstSectorRange sectorRange, secondSectorRange sectorRange) bool {
	return (firstSectorRange.start <= secondSectorRange.start && firstSectorRange.end >= secondSectorRange.start) ||
		(secondSectorRange.start <= firstSectorRange.start && secondSectorRange.end >= firstSectorRange.start)
}

var elfPairs []elfPair

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		pair := strings.Split(line, ",")
		if len(pair) < 2 {
			return
		}
		firstSector := strings.Split(pair[0], "-")
		secondSector := strings.Split(pair[1], "-")

		firstSectorStart, _ := strconv.Atoi(firstSector[0])
		firstSectorEnd, _ := strconv.Atoi(firstSector[1])
		secondSectorStart, _ := strconv.Atoi(secondSector[0])
		secondSectorEnd, _ := strconv.Atoi(secondSector[1])
		elfPair := elfPair{
			input: line,
			firstSectorRange: sectorRange{
				start: firstSectorStart,
				end:   firstSectorEnd,
			},
			secondSectorRange: sectorRange{
				start: secondSectorStart,
				end:   secondSectorEnd,
			},
		}
		elfPairs = append(elfPairs, elfPair)
	}
}

func main() {
	readFile()

	sum := 0
	for _, pair := range elfPairs {
		if isSectorRangePartOfOtherSectorRange(pair.firstSectorRange, pair.secondSectorRange) {
			sum++
		}
	}
	fmt.Println("Part01: ", sum)

	sum = 0
	for _, pair := range elfPairs {
		if isSectorOverlappingWithSectorRange(pair.firstSectorRange, pair.secondSectorRange) {
			sum++
		}
	}
	fmt.Println("Part02: ", sum)
}
