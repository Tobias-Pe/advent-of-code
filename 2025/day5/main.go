package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	start, end int
}

type IngredientDB struct {
	ranges []Range
}

func NewIngredientDB(lines []string) *IngredientDB {
	var rs []Range
	for _, line := range lines {
		args := strings.Split(line, "-")
		if len(args) != 2 {
			continue
		}
		start, _ := strconv.Atoi(args[0])
		end, _ := strconv.Atoi(args[1])
		rs = append(rs, Range{start, end})
	}
	idb := IngredientDB{ranges: rs}
	idb.MergeRanges()
	return &idb
}

func (db *IngredientDB) MergeRanges() {
	sort.Slice(db.ranges, func(i, j int) bool {
		if db.ranges[i].start == db.ranges[j].start {
			return db.ranges[i].end < db.ranges[j].end
		}
		return db.ranges[i].start < db.ranges[j].start
	})

	var mergedRanges []Range
	mergedRanges = append(mergedRanges, db.ranges[0])
	for i := 1; i < len(db.ranges); i++ {
		currRange := db.ranges[i]
		lastMRIdx := len(mergedRanges) - 1
		if currRange.start <= mergedRanges[lastMRIdx].end {
			if currRange.end >= mergedRanges[lastMRIdx].end {
				mergedRanges[lastMRIdx].end = currRange.end
			}
		} else {
			mergedRanges = append(mergedRanges, db.ranges[i])
		}
	}
	db.ranges = mergedRanges
}

func (db *IngredientDB) IsFresh(ingredID string) (bool, error) {
	if strings.Contains(ingredID, "-") {
		return false, fmt.Errorf("this is a range")
	}
	if ingredID == "" {
		return false, fmt.Errorf("empty ingredID")
	}
	ingredInt, _ := strconv.Atoi(ingredID)
	idx, found := slices.BinarySearchFunc(db.ranges, ingredInt, func(r Range, i int) int {
		return cmp.Compare(r.end, ingredInt)
	})
	if found {
		return true, nil
	}
	for i := max(0, idx-2); i < min(len(db.ranges), idx+2); i++ {
		currRange := db.ranges[i]
		if currRange.start <= ingredInt && currRange.end >= ingredInt {
			return true, nil
		}
	}
	return false, nil
}

func (db *IngredientDB) CountAvailableFreshIDs() int {
	counter := 0
	for _, rng := range db.ranges {
		counter += 1 + (rng.end - rng.start)
	}
	return counter
}

func main() {
	start := time.Now()

	input := readFile("day5/input.txt")
	idb := NewIngredientDB(input)
	freshCount := 0
	for _, line := range input {
		if isFresh, err := idb.IsFresh(line); err == nil && isFresh {
			freshCount++
		}
	}
	fmt.Println("Part01: 🍎🫧 ", freshCount)

	fmt.Println("Part02: 🍎🍏 ", idb.CountAvailableFreshIDs())

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
