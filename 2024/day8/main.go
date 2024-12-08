package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

type frequency string

type city struct {
	cityMap  [][]string
	antennas map[frequency][]coordinate
}

func (c city) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(c.cityMap) && coord.j < len(c.cityMap[coord.i])
}

func (c city) String() string {
	sb := strings.Builder{}
	for _, row := range c.cityMap {
		for _, s := range row {
			sb.WriteString(s)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c city) StringWithAntinodes(antinodes map[coordinate]bool) string {
	sb := strings.Builder{}
	for i, row := range c.cityMap {
		for j, s := range row {
			toBeWritten := s
			if _, ok := antinodes[coordinate{i, j}]; ok {
				if toBeWritten == "." {
					toBeWritten = "#"
				} else {
					toBeWritten = "*"
				}
			}
			sb.WriteString(toBeWritten)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func createCity(lines []string) *city {
	cityMap := make([][]string, len(lines))
	antennas := make(map[frequency][]coordinate)
	for i, line := range lines {
		cityRow := make([]string, len(line))
		for j, char := range line {
			cityRow[j] = string(char)
			if char != '.' {
				if _, ok := antennas[frequency(cityRow[j])]; ok {
					antennas[frequency(cityRow[j])] = append(antennas[frequency(cityRow[j])], coordinate{i, j})
				} else {
					antennas[frequency(cityRow[j])] = []coordinate{{i, j}}
				}
			}
		}
		cityMap[i] = cityRow
	}
	return &city{cityMap, antennas}
}

func (c city) discoverAntinodes(enableResonantHarmonics bool) map[coordinate]bool {
	antinodes := make(map[coordinate]bool)

	for _, antennaLocations := range c.antennas {
		for i := 0; i < len(antennaLocations); i++ {
			for j := i + 1; j < len(antennaLocations); j++ {
				var calcAntinodes []coordinate
				if enableResonantHarmonics {
					calcAntinodes = c.calculateAntinodesResonating(antennaLocations[i], antennaLocations[j])
				} else {
					calcAntinodes = c.calculateAntinodes(antennaLocations[i], antennaLocations[j])
				}
				for _, antinode := range calcAntinodes {
					if c.isValid(antinode) {
						antinodes[antinode] = true
					}
				}
			}
		}
	}

	return antinodes
}

func (c city) calculateAntinodes(antenna1, antenna2 coordinate) []coordinate {
	xD := antenna1.j - antenna2.j
	yD := antenna1.i - antenna2.i

	antinode2 := coordinate{antenna2.i - yD, antenna2.j - xD}
	antinode1 := coordinate{antenna1.i + yD, antenna1.j + xD}

	return []coordinate{antinode1, antinode2}
}

func (c city) calculateAntinodesResonating(antenna1, antenna2 coordinate) []coordinate {
	xD := antenna1.j - antenna2.j
	yD := antenna1.i - antenna2.i

	antinodes := c.spreadAntinodes(antenna2, -yD, -xD)
	antinodes = append(antinodes, c.spreadAntinodes(antenna1, yD, xD)...)
	antinodes = append(antinodes, coordinate{antenna1.i + yD, antenna1.j + xD})

	return antinodes
}

func (c city) spreadAntinodes(fromAntenna coordinate, yD int, xD int) []coordinate {
	var antinodes []coordinate
	iterator := 0
	nextAntinode := coordinate{fromAntenna.i + yD*iterator, fromAntenna.j + xD*iterator}
	for c.isValid(nextAntinode) {
		antinodes = append(antinodes, nextAntinode)
		iterator++
		nextAntinode = coordinate{fromAntenna.i + yD*iterator, fromAntenna.j + xD*iterator}
	}
	return antinodes
}

func main() {
	start := time.Now()

	calcSolutions("day8/input.txt")

	fmt.Println("Finished in", time.Since(start))
}

func calcSolutions(file string) {
	input := readFile(file)
	city := createCity(input)
	fmt.Println(city)

	antinodes := city.discoverAntinodes(false)
	fmt.Println(city.StringWithAntinodes(antinodes))
	fmt.Println("Part 1:", len(antinodes))

	harmonicAntinodes := city.discoverAntinodes(true)
	fmt.Println(city.StringWithAntinodes(harmonicAntinodes))
	fmt.Println("Part 2:", len(harmonicAntinodes))
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
