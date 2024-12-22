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
	baseCost   int
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

func (r race) getCheatDestinations(origin coordinate, dur int, visited map[coordinate]bool) map[cheat]int {
	walls := map[coordinate]int{}
	q := []coordinate{origin}
	for i := 1; i < dur; i++ {
		newQ := []coordinate{}
		for _, wall := range q {
			neighbours := wall.getNeighbours()
			for _, n := range neighbours {
				if _, ok := walls[n]; r.isValid(n) && r.racetrack[n.i][n.j] == "#" && !ok {
					walls[n] = i
					newQ = append(newQ, n)
				}
			}
		}
		q = newQ
	}
	possibilities := map[cheat]int{}
	for wall, dist := range walls {
		neighbours := wall.getNeighbours()
		for _, n := range neighbours {
			if n == origin || visited[n] {
				continue
			}
			if r.isValid(n) && r.racetrack[n.i][n.j] == "." {
				possibilities[cheat{
					from: origin,
					to:   n,
				}] = dist + 1
			}
		}
	}
	return possibilities
}

func (r *race) dfs(cheatDist int, curr coordinate, visited map[coordinate]bool, depthLvl int, currCheat *cheat) {
	if visited[curr] || !r.isValid(curr) || r.racetrack[curr.i][curr.j] == "#" || (r.baseCost != 0 && cheatDist > r.baseCost) {
		return
	}

	if curr == r.end {
		if currCheat != nil {
			r.cheatCosts[*currCheat] = depthLvl
		} else {
			r.baseCost = depthLvl
		}
		return
	}

	visited[curr] = true
	neighbours := curr.getNeighbours()
	for _, neighbour := range neighbours {
		if !r.isValid(neighbour) || visited[neighbour] {
			continue
		}
		r.dfs(cheatDist, neighbour, visited, depthLvl+1, currCheat)

		if currCheat != nil || r.racetrack[neighbour.i][neighbour.j] == "#" {
			continue
		}

		cheats := r.getCheatDestinations(curr, cheatDist, visited)
		for cheat, duration := range cheats {
			if _, ok := r.cheatCosts[cheat]; !r.isValid(cheat.to) || visited[cheat.to] || ok {
				continue
			}
			r.dfs(cheatDist, cheat.to, visited, depthLvl+duration, &cheat)
		}
	}
	delete(visited, curr)
}

func (r race) cheatCount(minSaveAmount int) int {
	costDistribution := make(map[int]int)
	for _, cost := range r.cheatCosts {
		saved := r.baseCost - cost
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
		saved := r.baseCost - cost
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
	rctrck.dfs(2, rctrck.start, make(map[coordinate]bool), 0, nil)
	fmt.Println("Part 1:", rctrck.cheatCount(100))
	fmt.Println("P1 Finished in", time.Since(start))

	start = time.Now()
	rctrck2 := readFile(file)
	rctrck2.dfs(20, rctrck2.start, make(map[coordinate]bool), 0, nil)
	fmt.Println("Part 2:", rctrck2.cheatCount(100))
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
	race := race{cheatCosts: map[cheat]int{}}
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
