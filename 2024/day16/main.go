package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
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

func (c coordinate) turnRight() coordinate {
	// rotation matrix clockwise 90°
	// i1 = j0 --> i = row; j = col
	// j1 = -i0
	return coordinate{c.j, -c.i}
}

func (c coordinate) turnLeft() coordinate {
	return coordinate{-c.j, c.i}
}

type destination struct {
	c    coordinate
	cost int64
}

func (d destination) String() string {
	return fmt.Sprintf("%s,%d", d.c.String(), d.cost)
}

type state struct {
	c    coordinate
	dir  coordinate
	cost int64
}

type maze struct {
	tiles [][]string
	start coordinate
	end   coordinate
}

func (m maze) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(m.tiles) && coord.j < len(m.tiles[coord.i]) && m.tiles[coord.i][coord.j] != "#"
}

func (m maze) String() string {
	sb := strings.Builder{}
	for _, t := range m.tiles {
		for _, cell := range t {
			sb.WriteString(cell)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m maze) dsa() (int64, map[coordinate]float64) {
	costs := make(map[coordinate]float64)
	for i := range m.tiles {
		for j := range m.tiles[i] {
			costs[coordinate{i, j}] = math.Inf(1)
		}
	}
	costs[m.start] = 0
	pq := &stateHeap{state{m.start, coordinate{0, 1}, 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(state)
		currDir := curr.dir
		neighbour := curr.c.add(currDir)
		if m.isValid(neighbour) && costs[neighbour] > costs[curr.c]+1 {
			costs[neighbour] = costs[curr.c] + 1
			heap.Push(pq, state{neighbour, currDir, int64(costs[neighbour])})
		}
		currDir = curr.dir.turnRight()
		neighbour = curr.c.add(currDir)
		if m.isValid(neighbour) && costs[neighbour] > costs[curr.c]+1001 {
			costs[neighbour] = costs[curr.c] + 1001
			heap.Push(pq, state{neighbour, currDir, int64(costs[neighbour])})
		}
		currDir = currDir.turnRight()
		neighbour = curr.c.add(currDir)
		if m.isValid(neighbour) && costs[neighbour] > costs[curr.c]+2001 {
			costs[neighbour] = costs[curr.c] + 2001
			heap.Push(pq, state{neighbour, currDir, int64(costs[neighbour])})
		}
		currDir = curr.dir.turnLeft()
		neighbour = curr.c.add(currDir)
		if m.isValid(neighbour) && costs[neighbour] > costs[curr.c]+1001 {
			costs[neighbour] = costs[curr.c] + 1001
			heap.Push(pq, state{neighbour, currDir, int64(costs[neighbour])})
		}
		//m.progress(costs)
	}
	return int64(costs[m.end]), costs
}

func (m maze) print(valid map[coordinate]bool) {
	for i, t := range m.tiles {
		for j, cell := range t {
			if valid[coordinate{i, j}] {
				fmt.Print("O")
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}
}

func (m maze) progress(mapC map[coordinate]float64) {
	for i, t := range m.tiles {
		for j, cell := range t {
			if !math.IsInf(mapC[coordinate{i, j}], 1) {
				fmt.Print("O")
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}
}

func dfsExactCost(graph map[coordinate][]coordinate, curr coordinate, m maze, valid *map[coordinate]bool, prev map[coordinate]bool) bool {
	if present, ok := prev[curr]; ok && present {
		return false
	}
	if curr == m.end {
		(*valid)[curr] = true
		return true
	}
	destinations := graph[curr]
	endIsReachable := false
	prev[curr] = true
	//m.print(prev)
	for _, dest := range destinations {
		foundTheEnd := dfsExactCost(graph, dest, m, valid, prev)
		if foundTheEnd {
			(*valid)[curr] = true
			endIsReachable = true
		}
	}
	prev[curr] = false
	return endIsReachable
}

func main() {
	start := time.Now()

	m := readFile("day16/input_exp_1.txt") // TODO delte me: p1 sol 108504
	costToEnd, _ := m.dsa()
	fmt.Println("Part 1:", costToEnd)

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) maze {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	maze := maze{}
	maze.tiles = make([][]string, len(split))
	for i, line := range split {
		tileRow := strings.Split(line, "")
		for j, char := range tileRow {
			if char == "S" {
				maze.start = coordinate{i, j}
			} else if char == "E" {
				maze.end = coordinate{i, j}
			}
		}
		maze.tiles[i] = tileRow
	}
	return maze
}
