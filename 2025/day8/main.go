package main

import (
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	jb1, jb2 *JunctionBox
}

func (conn *Connection) Distance() float64 {
	box, other := conn.jb1, conn.jb2
	return math.Sqrt(math.Pow(float64(box.x-other.x), 2) + math.Pow(float64(box.y-other.y), 2) + math.Pow(float64(box.z-other.z), 2))
}

type JunctionBox struct {
	x, y, z     int
	connectedTo []*JunctionBox
}

func (box *JunctionBox) String() string {
	return fmt.Sprintf("(%d,%d,%d)", box.x, box.y, box.z)
}

func (box *JunctionBox) PopulateCirc(in map[*JunctionBox]bool) map[*JunctionBox]bool {
	in[box] = true
	for _, b := range box.connectedTo {
		if !in[b] {
			in[b] = true
			in = b.PopulateCirc(in)
		}
	}
	return in
}

func (box *JunctionBox) Connect(other *JunctionBox) bool {
	//if len(box.connectedTo) >= 2 || len(other.connectedTo) >= 2 {
	//	return false
	//}
	if box.PopulateCirc(make(map[*JunctionBox]bool))[other] == true { // expensive
		return false
	}
	box.connectedTo = append(box.connectedTo, other)
	other.connectedTo = append(other.connectedTo, box)
	return true
}

type Playground struct {
	boxes               []*JunctionBox
	possibleConnections []*Connection
}

func NewPlayground(lines []string) *Playground {
	pg := &Playground{}

	for _, line := range lines {
		args := strings.Split(line, ",")
		x, _ := strconv.Atoi(args[0])
		y, _ := strconv.Atoi(args[1])
		z, _ := strconv.Atoi(args[2])
		jb := &JunctionBox{
			x: x,
			y: y,
			z: z,
		}
		pg.boxes = append(pg.boxes, jb)
	}

	for i, jb1 := range pg.boxes {
		for j := i + 1; j < len(pg.boxes); j++ {
			jb2 := pg.boxes[j]
			pg.possibleConnections = append(pg.possibleConnections, &Connection{
				jb1: jb1,
				jb2: jb2,
			})
		}
	}

	slices.SortFunc(pg.possibleConnections, func(a, b *Connection) int {
		return cmp.Compare(a.Distance(), b.Distance())
	})

	return pg
}

func (pg *Playground) ConnectShortestN(n int) {
	for i := 0; i < n; i++ {
		con := pg.possibleConnections[i]
		fmt.Print(i+1, ":\t", con.jb1, "\t➡️\t", con.jb2, "\t:\t")
		diff := con.jb1.Connect(con.jb2)
		cs := pg.Circs()
		if diff {
			pg.PrintCircs(cs...)
		} else {
			fmt.Println("Nothing happened")
		}
		if len(cs) == 1 {
			return
		}
	}
}

func (pg *Playground) Circs() []map[*JunctionBox]bool {
	visited := make(map[*JunctionBox]bool)
	var circs []map[*JunctionBox]bool
	for _, jb := range pg.boxes {
		if !visited[jb] {
			circ := jb.PopulateCirc(make(map[*JunctionBox]bool))
			circs = append(circs, circ)
			for cJb := range circ {
				visited[cJb] = true
			}
		}
	}

	slices.SortFunc(circs, func(a, b map[*JunctionBox]bool) int {
		return cmp.Compare(len(b), len(a))
	})
	return circs
}

func (pg *Playground) PrintCircs(cs ...map[*JunctionBox]bool) {
	sb := &strings.Builder{}
	for _, c := range cs {
		sb.WriteString("[")
		for b := range c {
			sb.WriteString(b.String())
			sb.WriteString(",")
		}
		sb.WriteString("] len")
		sb.WriteString(strconv.Itoa(len(c)))
		sb.WriteString("\t")
	}
	prod := int64(1)
	for i := 0; i < len(cs)-1 && i < 3; i++ {
		prod *= int64(len(cs[i]))
	}
	fmt.Println(prod, sb.String())
}

func main() {
	start := time.Now()

	lines := readFile("day8/input.txt")
	pg := NewPlayground(lines)
	pg.ConnectShortestN(math.MaxInt)

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
