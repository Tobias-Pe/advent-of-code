package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type region struct {
	flowers            []coordinate
	coordinatePresence map[coordinate]bool
	flowerType         string
}

func (r region) getFencePricing(ga gardenArrangement) int {
	return r.calcPerimeter(ga) * len(r.flowers)
}

func (r region) getFencePricingBulkDiscount(ga gardenArrangement) int {
	return r.calcSides(ga) * len(r.flowers)
}

func (r region) String() string {
	sb := strings.Builder{}
	sb.WriteString(r.flowerType)
	sb.WriteString("> ")
	for _, flower := range r.flowers {
		sb.WriteString(flower.String())
		sb.WriteString(" ")
	}
	sb.WriteString(";")
	return sb.String()
}

func (r region) calcPerimeter(ga gardenArrangement) int {
	perimenter := 0
	queue := append([]coordinate{}, r.flowers...)
	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]
		neighbours := coord.getNeighbours()
		for _, neighbour := range neighbours {
			if !ga.isValid(neighbour) || ga.garden[neighbour.i][neighbour.j] != r.flowerType {
				perimenter++
			}
		}
	}
	return perimenter
}

func (r region) calcSides(ga gardenArrangement) int {
	sides := 0
	queue := append([]coordinate{}, r.flowers...)
	alreadyHasFence := make(map[vector]bool)
	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]
		neighbours := coord.getNeighbours()
		for _, neighbour := range neighbours {
			if !ga.isValid(neighbour) || ga.garden[neighbour.i][neighbour.j] != r.flowerType {
				success := r.growSide(neighbour, coord, alreadyHasFence, ga)
				if success {
					sides++
				}
			}
		}
	}
	fmt.Println(r.flowerType, sides)
	return sides
}

func (r region) growSide(newFence coordinate, origin coordinate, alreadyHasFence map[vector]bool, ga gardenArrangement) bool {
	if alreadyHasFence[vector{origin: origin, target: newFence}] {
		return false
	}
	dI, dJ := newFence.i-origin.i, newFence.j-origin.j
	coord := origin
	for r.coordinatePresence[coord] {
		fenceLoc := coordinate{coord.i + dI, coord.j + dJ}
		if ga.isValid(fenceLoc) && ga.garden[fenceLoc.i][fenceLoc.j] == r.flowerType {
			break
		}
		alreadyHasFence[vector{origin: coord, target: fenceLoc}] = true
		coord = coordinate{coord.i + dJ, coord.j + dI}
	}
	coord = coordinate{origin.i - dJ, origin.j - dI}
	for r.coordinatePresence[coord] {
		fenceLoc := coordinate{coord.i + dI, coord.j + dJ}
		if ga.isValid(fenceLoc) && ga.garden[fenceLoc.i][fenceLoc.j] == r.flowerType {
			break
		}
		alreadyHasFence[vector{origin: coord, target: fenceLoc}] = true
		coord = coordinate{coord.i - dJ, coord.j - dI}
	}
	return true
}

type vector struct {
	origin, target coordinate
}

type coordinate struct {
	i, j int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

type gardenArrangement struct {
	garden  [][]string
	regions []region
}

func (ga gardenArrangement) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(ga.garden) && coord.j < len(ga.garden[coord.i])
}

func (ga gardenArrangement) spreadRegion(c coordinate, visited map[coordinate]bool) *region {
	if visited[c] {
		return nil
	}
	reg := region{
		flowers:            []coordinate{},
		flowerType:         "",
		coordinatePresence: map[coordinate]bool{},
	}

	queue := []coordinate{c}
	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]
		if visited[coord] {
			continue
		}
		visited[coord] = true
		if reg.flowerType == "" {
			reg.flowerType = ga.garden[coord.i][coord.j]
		}
		if reg.flowerType == ga.garden[coord.i][coord.j] {
			neighbours := coord.getNeighbours()
			for _, neighbour := range neighbours {
				if !visited[neighbour] && ga.isValid(neighbour) && reg.flowerType == ga.garden[neighbour.i][neighbour.j] {
					queue = append(queue, neighbour)
				}
			}
		}
		reg.flowers = append(reg.flowers, coord)
		reg.coordinatePresence[coord] = true
	}

	return &reg
}

func (ga gardenArrangement) calculateFencePricing() int {
	sum := 0
	for _, r := range ga.regions {
		sum += r.getFencePricing(ga)
	}
	return sum
}

func (ga gardenArrangement) calculateFencePricingBulkDiscount() int {
	sum := 0
	for _, r := range ga.regions {
		sum += r.getFencePricingBulkDiscount(ga)
	}
	return sum
}

func createGardenArrangement(lines []string) gardenArrangement {
	garden := make([][]string, len(lines))
	for i, line := range lines {
		garden[i] = strings.Split(line, "")
	}
	arrangement := gardenArrangement{garden: garden}
	regions := []region{}
	visited := map[coordinate]bool{}
	for i := 0; i < len(garden); i++ {
		for j := 0; j < len(garden[i]); j++ {
			region := arrangement.spreadRegion(coordinate{i, j}, visited)
			if region != nil {
				regions = append(regions, *region)
			}
		}
	}
	arrangement.regions = regions
	return arrangement
}

func main() {
	start := time.Now()

	input := readFile("day12/input.txt")
	arrangement := createGardenArrangement(input)
	fmt.Println("Part 1:", arrangement.calculateFencePricing())
	fmt.Println("Part 2:", arrangement.calculateFencePricingBulkDiscount())

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
