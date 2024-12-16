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
	v    vector
	cost int
}

type vector struct {
	c   coordinate
	dir coordinate
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

func (m maze) dsa() (int, int) {
	prevs := make(map[vector][]vector)
	costGroups := make(map[coordinate][]vector)
	costs := make(map[vector]int)
	initState := state{v: vector{c: m.start, dir: coordinate{0, 1}}, cost: 0}
	costs[initState.v] = initState.cost
	pq := &stateHeap{initState}
	heap.Init(pq)
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(state)
		neigbourState := state{v: vector{c: curr.v.c.add(curr.v.dir), dir: curr.v.dir}, cost: curr.cost + 1}
		if _, ok := costs[neigbourState.v]; m.isValid(neigbourState.v.c) && !ok {
			costs[neigbourState.v] = neigbourState.cost
			heap.Push(pq, neigbourState)
			costGroups[neigbourState.v.c] = append(costGroups[neigbourState.v.c], neigbourState.v)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		} else if m.isValid(neigbourState.v.c) && costs[neigbourState.v] == neigbourState.cost {
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		}
		currDir := curr.v.dir.turnRight()
		neigbourState = state{v: vector{c: curr.v.c.add(currDir), dir: currDir}, cost: curr.cost + 1001}
		if _, ok := costs[neigbourState.v]; m.isValid(neigbourState.v.c) && !ok {
			costs[neigbourState.v] = neigbourState.cost
			heap.Push(pq, neigbourState)
			costGroups[neigbourState.v.c] = append(costGroups[neigbourState.v.c], neigbourState.v)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		} else if m.isValid(neigbourState.v.c) && costs[neigbourState.v] == neigbourState.cost {
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		}
		currDir = curr.v.dir.turnRight().turnRight()
		neigbourState = state{v: vector{c: curr.v.c.add(currDir), dir: currDir}, cost: curr.cost + 2001}
		if _, ok := costs[neigbourState.v]; m.isValid(neigbourState.v.c) && !ok {
			costs[neigbourState.v] = neigbourState.cost
			heap.Push(pq, neigbourState)
			costGroups[neigbourState.v.c] = append(costGroups[neigbourState.v.c], neigbourState.v)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		} else if m.isValid(neigbourState.v.c) && costs[neigbourState.v] == neigbourState.cost {
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		}
		currDir = curr.v.dir.turnLeft()
		neigbourState = state{v: vector{c: curr.v.c.add(currDir), dir: currDir}, cost: curr.cost + 1001}
		if _, ok := costs[neigbourState.v]; m.isValid(neigbourState.v.c) && !ok {
			costs[neigbourState.v] = neigbourState.cost
			heap.Push(pq, neigbourState)
			costGroups[neigbourState.v.c] = append(costGroups[neigbourState.v.c], neigbourState.v)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		} else if m.isValid(neigbourState.v.c) && costs[neigbourState.v] == neigbourState.cost {
			prevs[neigbourState.v] = append(prevs[neigbourState.v], prevs[curr.v]...)
			prevs[neigbourState.v] = append(prevs[neigbourState.v], curr.v)
		}
	}
	minCost := math.MaxInt64
	for _, vec := range costGroups[m.end] {
		minCost = min(minCost, costs[vec])
	}
	q := []vector{}
	for _, vec := range costGroups[m.end] {
		if costs[vec] == minCost {
			q = append(q, prevs[vec]...)
		}
	}
	unique := map[coordinate]bool{}
	for _, v := range q {
		unique[v.c] = true
	}
	unique[m.end] = true
	unique[m.start] = true
	m.print(unique)
	return minCost, len(unique)
}

func (m maze) dfsExactCost(curr vector, costToEnd int, valid map[coordinate]bool, prev map[vector]bool, connections map[vector]map[state]bool) bool {
	if present, ok := prev[curr]; ok && present {
		return false
	}
	if costToEnd < 0 {
		return false
	}
	if curr.c == m.end && costToEnd == 0 {
		valid[curr.c] = true
		return true
	}
	destinations := connections[curr]
	endIsReachable := false
	prev[curr] = true
	for dest := range destinations {
		foundTheEnd := m.dfsExactCost(dest.v, costToEnd-dest.cost, valid, prev, connections)
		if foundTheEnd {
			valid[curr.c] = true
			endIsReachable = true
		}
	}
	prev[curr] = false
	return endIsReachable
}

func (m maze) print(valid map[coordinate]bool) {
	for i, t := range m.tiles {
		for j, cell := range t {
			if valid[coordinate{i, j}] {
				fmt.Print("🦌")
			} else {
				if cell == "." {
					fmt.Print("⬜")
				} else {
					fmt.Print("🌲")
				}
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

func main() {
	start := time.Now()

	m := readFile("day16/input.txt")
	costToEnd, connections := m.dsa()
	fmt.Println("Part 1:", costToEnd)
	fmt.Println("Part 2:", connections)

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
