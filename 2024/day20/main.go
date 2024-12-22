package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type coordinate struct {
	i, j int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.i, c.j)
}

func (c coordinate) getNeighbours() []coordinate {
	return []coordinate{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}}
}

func (c coordinate) sub(other coordinate) coordinate {
	return coordinate{c.i - other.i, c.j - other.j}
}
func (c coordinate) add(other coordinate) coordinate {
	return coordinate{c.i + other.i, c.j + other.j}
}

type cheat struct {
	from coordinate
	to   coordinate
}

type race struct {
	racetrack  [][]string
	start, end coordinate
	cheatCosts map[cheat]int
	path       map[coordinate]int
}

func (r race) isValid(coord coordinate) bool {
	return coord.i >= 0 && coord.j >= 0 && coord.i < len(r.racetrack) && coord.j < len(r.racetrack[coord.i])
}

func (r race) print(visited map[coordinate]bool) {
	sb := strings.Builder{}
	for i, row := range r.racetrack {
		for j, field := range row {
			if visited[coordinate{i, j}] {
				sb.WriteString("O")
			} else {
				sb.WriteString(field)
			}
		}
		sb.WriteRune('\n')
	}
	fmt.Println(sb.String())
}

func (r race) getCheatDestinations(origin coordinate, dur int) map[cheat]int {
	walls := map[coordinate]int{}
	q := []coordinate{origin}
	for i := 1; i <= dur; i++ {
		newQ := []coordinate{}
		for _, wall := range q {
			neighbours := wall.getNeighbours()
			for _, n := range neighbours {
				if _, ok := walls[n]; r.isValid(n) && !ok {
					walls[n] = i
					newQ = append(newQ, n)
				}
			}
		}
		q = newQ
	}
	possibilities := map[cheat]int{}
	for field, dist := range walls {
		if r.isValid(field) && field != origin && r.racetrack[field.i][field.j] == "." {
			possibilities[cheat{
				from: origin,
				to:   field,
			}] = dist
		}
	}
	return possibilities
}

func (r race) calcCheatSaves(cheatDist int) {
	for c, costFromStart := range r.path {
		cheats := r.getCheatDestinations(c, cheatDist)
		for cheat, duration := range cheats {
			if _, ok := r.cheatCosts[cheat]; !r.isValid(cheat.to) || ok {
				continue
			}
			costToEnd := r.path[r.end] - r.path[cheat.to]
			r.cheatCosts[cheat] = duration + costToEnd + costFromStart
		}
	}
}

func (r *race) dfsPath(curr coordinate, from coordinate, depthLvl int) {
	if !r.isValid(curr) || r.racetrack[curr.i][curr.j] == "#" {
		return
	}

	r.path[curr] = depthLvl
	if curr == r.end {
		fmt.Println("Base Found")
		return
	}

	neighbours := curr.getNeighbours()
	for _, neighbour := range neighbours {
		if !r.isValid(neighbour) || from == neighbour || r.racetrack[neighbour.i][neighbour.j] == "#" {
			continue
		}
		r.dfsPath(neighbour, curr, depthLvl+1)
	}
}

func (r race) cheatCount(minSaveAmount int) int {
	costDistribution := make(map[int]int)
	for _, cost := range r.cheatCosts {
		saved := r.path[r.end] - cost
		if saved >= minSaveAmount {
			costDistribution[saved]++
		}
	}
	savedCosts := []int{}
	for cost, _ := range costDistribution {
		savedCosts = append(savedCosts, cost)
	}
	slices.Sort(savedCosts)
	for _, cost := range savedCosts {
		fmt.Printf("There are %d cheats that save %d picoseconds.\n", costDistribution[cost], cost)
	}

	sum := 0
	for _, cost := range r.cheatCosts {
		saved := r.path[r.end] - cost
		if saved >= minSaveAmount {
			sum++
		}
	}
	return sum
}

func main() {
	file := "day20/input.txt"

	start := time.Now()
	rctrck := readFile(file)
	rctrck.dfsPath(rctrck.start, coordinate{}, 0)
	rctrck.calcCheatSaves(2)
	fmt.Println("Part 1:", rctrck.cheatCount(100))
	fmt.Println("P1 Finished in", time.Since(start))

	start = time.Now()
	rctrck.calcCheatSaves(20)
	fmt.Println("Part 2:", rctrck.cheatCount(100))
	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) race {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	race := race{cheatCosts: map[cheat]int{}, path: map[coordinate]int{}}
	racetrack := make([][]string, len(split))
	for i := range racetrack {
		racetrack[i] = make([]string, len(split[i]))
		for j, s := range strings.Split(split[i], "") {
			racetrack[i][j] = s
			if s == "S" {
				racetrack[i][j] = "."
				race.start = coordinate{i, j}
			}
			if s == "E" {
				racetrack[i][j] = "."
				race.end = coordinate{i, j}
			}
		}
	}
	race.racetrack = racetrack
	return race
}
