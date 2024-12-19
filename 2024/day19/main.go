package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type onsen struct {
	towelPatterns  map[string]bool
	maxPatternLen  int
	desiredDesigns []string
	isPossibleDP   map[string]bool
	waysDP         map[string]int
}

func (o onsen) countWaysForDesigns() int {
	sum := 0
	for _, d := range o.desiredDesigns {
		sum += o.countWaysForDesign(d)
	}
	return sum
}

func (o onsen) countViableDesigns() int {
	sum := 0
	for _, d := range o.desiredDesigns {
		if o.isPossible(d) {
			sum += 1
		}
	}
	return sum
}

func (o onsen) isPossible(design string) bool {
	if len(design) == 0 {
		return true
	}
	if out, ok := o.isPossibleDP[design]; ok {
		return out
	}
	match := false
	for i := 0; i < o.maxPatternLen; i++ {
		if len(design) <= i {
			continue
		}
		currPart := design[:i+1]
		counterPart := design[i+1:]
		if o.towelPatterns[currPart] {
			match = o.isPossible(counterPart)
		}
		if match {
			break
		}
	}

	o.isPossibleDP[design] = match
	return o.isPossibleDP[design]
}

func (o onsen) countWaysForDesign(design string) int {
	if len(design) == 0 {
		return 1
	}
	if out, ok := o.waysDP[design]; ok {
		return out
	}
	matches := 0
	for i := 0; i < o.maxPatternLen; i++ {
		if len(design) <= i {
			continue
		}
		currPart := design[:i+1]
		counterPart := design[i+1:]
		if o.towelPatterns[currPart] {
			matches += o.countWaysForDesign(counterPart)
		}
	}

	o.waysDP[design] = matches
	return o.waysDP[design]
}

func main() {
	start := time.Now()

	onsen := readFile("day19/input.txt")
	fmt.Println("Part 1:", onsen.countViableDesigns())
	fmt.Println("Part 2:", onsen.countWaysForDesigns())

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) onsen {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n\n")
	onsen := onsen{isPossibleDP: make(map[string]bool), waysDP: make(map[string]int)}
	onsen.towelPatterns = make(map[string]bool)
	for _, towelPattern := range strings.Split(split[0], ", ") {
		tP := strings.TrimSpace(towelPattern)
		onsen.towelPatterns[tP] = true
		onsen.maxPatternLen = max(onsen.maxPatternLen, len(tP))
	}
	onsen.desiredDesigns = []string{}
	for _, desiredPattern := range strings.Split(split[1], "\n") {
		onsen.desiredDesigns = append(onsen.desiredDesigns, strings.TrimSpace(desiredPattern))
	}
	return onsen
}
