package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) add(other coordinate) coordinate {
	return coordinate{c.i + other.i, c.j + other.j}
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

type corruptingMemory struct {
	fallingBytes []coordinate
	memorySpace  [][]bool
}

func (cM *corruptingMemory) String() string {
	sb := strings.Builder{}
	for _, c := range cM.memorySpace {
		for _, b := range c {
			if b {
				sb.WriteString("⬜")
			} else {
				sb.WriteString("🕳️")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (cM corruptingMemory) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(cM.memorySpace) && coord.j < len(cM.memorySpace[coord.i]) && cM.memorySpace[coord.i][coord.j]
}

func (cM *corruptingMemory) letTimePass() {
	fallenByte := cM.fallingBytes[0]
	cM.fallingBytes = cM.fallingBytes[1:]
	cM.memorySpace[fallenByte.i][fallenByte.j] = false
}

func (cM *corruptingMemory) letTimePassTimes(c int) {
	for range c {
		cM.letTimePass()
	}
}

func (cM *corruptingMemory) bfs() (int, bool) {
	start := coordinate{0, 0}
	end := coordinate{len(cM.memorySpace) - 1, len(cM.memorySpace[len(cM.memorySpace)-1]) - 1}

	costMap := map[coordinate]int{}
	costMap[start] = 0
	q := []coordinate{start}
	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		neighbours := current.getNeighbours()
		for _, neighbour := range neighbours {
			if _, ok := costMap[neighbour]; cM.isValid(neighbour) && !ok {
				costMap[neighbour] = costMap[current] + 1
				q = append(q, neighbour)
			}
		}
		if cost, ok := costMap[end]; ok {
			return cost, true
		}
	}
	cost, ok := costMap[end]
	return cost, ok
}

func findPathBlockingByte(file string) coordinate {
	fallingBytes := readFile(file, 71, 71).fallingBytes
	left := 0
	right := len(fallingBytes) - 1
	lastFound := -1
	for left < right {
		mid := ((right - left) / 2) + left
		mem := readFile(file, 71, 71)
		mem.letTimePassTimes(mid)
		_, found := mem.bfs()
		if found {
			left = mid + 1
			lastFound = mid
		} else {
			right = mid - 1
		}
	}
	return coordinate{fallingBytes[lastFound].j, fallingBytes[lastFound].i}
}

func main() {
	start := time.Now()

	//memory := readFile("day18/input_exp.txt", 7, 7)
	//memory.letTimePassTimes(12)
	file := "day18/input.txt"
	memory := readFile(file, 71, 71)
	memory.letTimePassTimes(1024)
	cost, _ := memory.bfs()
	fmt.Println("Part 1:", cost)
	fmt.Println(memory.String())
	fmt.Println("Part 2:", findPathBlockingByte(file))
	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string, sizeX, sizeY int) corruptingMemory {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	fallingBytes := make([]coordinate, len(split))
	for i, s := range split {
		c := coordinate{}
		coords := strings.Split(s, ",")
		atoiX, _ := strconv.Atoi(coords[0])
		atoiY, _ := strconv.Atoi(coords[1])
		c.i = atoiY
		c.j = atoiX
		fallingBytes[i] = c
	}
	memorySpace := make([][]bool, sizeY)
	for i := range memorySpace {
		memorySpace[i] = make([]bool, sizeX)
		for j := range memorySpace[i] {
			memorySpace[i][j] = true
		}
	}
	return corruptingMemory{
		fallingBytes: fallingBytes,
		memorySpace:  memorySpace,
	}
}
