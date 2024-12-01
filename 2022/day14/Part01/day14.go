package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type unit struct {
	materialType int // 0 = air, 1 = stone, 2 = sand, 3 = sand_source
	i            int
	j            int
}

var cave [][]*unit

var sandSource *unit

func readFile(file string) {
	cave = make([][]*unit, 700)
	for i := range cave {
		cave[i] = make([]*unit, 700)
		for j := 0; j < len(cave[i]); j++ {
			cave[i][j] = &unit{
				materialType: 0,
				i:            i,
				j:            j,
			}
		}
	}
	sandSource = &unit{
		materialType: 3,
		i:            0,
		j:            500,
	}
	cave[0][500] = sandSource

	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	splittedLines := strings.Split(lines, "\n")
	for _, line := range splittedLines {
		connectedCoordinates := strings.Split(line, " -> ")
		for indexCoordinate := 0; indexCoordinate < len(connectedCoordinates)-1; indexCoordinate++ {
			coordinates := strings.Split(connectedCoordinates[indexCoordinate], ",")
			coordinateJ, _ := strconv.Atoi(coordinates[0])
			coordinateI, _ := strconv.Atoi(coordinates[1])
			cave[coordinateI][coordinateJ].materialType = 1

			coordinates2 := strings.Split(connectedCoordinates[indexCoordinate+1], ",")
			coordinateJ2, _ := strconv.Atoi(coordinates2[0])
			coordinateI2, _ := strconv.Atoi(coordinates2[1])
			cave[coordinateI2][coordinateJ2].materialType = 1

			createRockLine(cave[coordinateI][coordinateJ], cave[coordinateI2][coordinateJ2])
		}
	}
}

func createRockLine(firstRock *unit, secondRock *unit) {
	isHorizontal := firstRock.i-secondRock.i != 0
	startIndex := -1
	endIndex := -1

	if isHorizontal {
		if firstRock.i > secondRock.i {
			startIndex = secondRock.i
			endIndex = firstRock.i
		} else {
			startIndex = firstRock.i
			endIndex = secondRock.i
		}

		for i := startIndex; i <= endIndex; i++ {
			cave[i][firstRock.j].materialType = 1
		}
	} else {

		if firstRock.j > secondRock.j {
			startIndex = secondRock.j
			endIndex = firstRock.j
		} else {
			startIndex = firstRock.j
			endIndex = secondRock.j
		}

		for j := startIndex; j <= endIndex; j++ {
			cave[firstRock.i][j].materialType = 1
		}
	}
}

func printCave() {
	for i, units := range cave {
		if i < 15 {
			for j, unit := range units {
				if j > 450 && j < 550 {
					switch unit.materialType {
					case 0:
						fmt.Print(".")
					case 1:
						fmt.Print("#")
					case 2:
						fmt.Print("o")
					case 3:
						fmt.Print("+")
					}
				}
			}
			fmt.Print("\n")
		}
	}
	fmt.Println("\n---------------------------------------------------------------------------------------------")
}

func createSand() *unit {
	for i := sandSource.i + 1; i < len(cave); i++ {
		if cave[i][sandSource.j].materialType != 0 {
			cave[i-1][sandSource.j].materialType = 2
			return cave[i-1][sandSource.j]
		}
	}
	return nil
}

func simulateSand() {
	counter := 0
	for sand := createSand(); sand != nil; sand = createSand() {

		hasFallenOff := comeToRest(sand)
		if hasFallenOff {
			break
		}
		counter++
		//fmt.Println(counter)
		//printCave()
	}
	fmt.Println(counter)
}

// returns if the sand corn has fallen off the world
func comeToRest(sand *unit) bool {
	if sand.i+1 >= len(cave) {
		return true
	}

	if cave[sand.i+1][sand.j].materialType == 0 {
		cave[sand.i+1][sand.j].materialType = 2
		sand.materialType = 0
		return comeToRest(cave[sand.i+1][sand.j])
	} else if cave[sand.i+1][sand.j-1].materialType == 0 {
		cave[sand.i+1][sand.j-1].materialType = 2
		sand.materialType = 0
		return comeToRest(cave[sand.i+1][sand.j-1])
	} else if cave[sand.i+1][sand.j+1].materialType == 0 {
		cave[sand.i+1][sand.j+1].materialType = 2
		sand.materialType = 0
		return comeToRest(cave[sand.i+1][sand.j+1])
	}
	return false
}

func main() {
	readFile("input.txt")
	simulateSand()
}
