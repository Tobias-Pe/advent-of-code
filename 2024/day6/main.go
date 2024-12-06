package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"os"
	"strings"
	"time"
)

type coord struct {
	row, col int
}

func (c *coord) add(other coord) {
	c.row += other.row
	c.col += other.col
}

type lab struct {
	roomTiles         [][]string
	guard             coord
	guardDir          coord
	tileCoordsTouched map[coord]bool
}

type game struct {
	laboratory *lab
}

func (g game) Update() error {
	if g.laboratory.tileCoordsTouched == nil {
		g.laboratory.tileCoordsTouched = make(map[coord]bool)
	}
	if g.laboratory.isValid(g.laboratory.guard) {
		g.laboratory.tileCoordsTouched[g.laboratory.guard] = true
		g.laboratory.moveGuard()
	}
	return nil
}

func (g game) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	const size = 27
	op.LineSpacing = size * 1.01
	s := g.laboratory.String()
	text.Draw(screen, s, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   size,
	}, op)
}

func (g game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//TODO implement me
	return len(g.laboratory.roomTiles[0]) * 30, len(g.laboratory.roomTiles) * 30
}

func (l *lab) isValid(c coord) bool {
	return c.row >= 0 && c.col >= 0 && c.row < len(l.roomTiles) && c.col < len(l.roomTiles[c.row])
}

func parseLab(lines []string) *lab {
	laboratory := lab{
		roomTiles: make([][]string, len(lines)),
		guard:     coord{},
	}
	for i, line := range lines {
		row := make([]string, 0, len(line))
		for j, tile := range line {
			switch {
			case string(tile) == "v":
				laboratory.guard = coord{i, j}
				laboratory.guardDir = coord{1, 0}
				row = append(row, ".")
			case string(tile) == "^":
				laboratory.guard = coord{i, j}
				laboratory.guardDir = coord{-1, 0}
				row = append(row, ".")
			case string(tile) == ">":
				laboratory.guard = coord{i, j}
				laboratory.guardDir = coord{0, 1}
				row = append(row, ".")
			case string(tile) == "<":
				laboratory.guard = coord{i, j}
				laboratory.guardDir = coord{0, -1}
				row = append(row, ".")
			default:
				row = append(row, string(tile))
			}
		}
		laboratory.roomTiles[i] = row
	}
	return &laboratory
}

func (l *lab) print(visualize bool) {
	if !visualize {
		return
	}
	fmt.Println(l)
}

func (l *lab) String() string {
	output := strings.Builder{}
	for i, tile := range l.roomTiles {
		for j, s := range tile {
			currentCord := coord{i, j}
			if i == l.guard.row && j == l.guard.col {
				output.WriteString("!")
			} else if s == "#" {
				output.WriteString("¤")
			} else if l.tileCoordsTouched[currentCord] {
				output.WriteString("★")
			} else {
				output.WriteString(".")
			}
		}
		output.WriteString("\n")
	}
	output.WriteString(fmt.Sprintf("Guard Dir: %v\n", l.guardDir))
	output.WriteString(fmt.Sprintf("Tilecount: %v\n", len(l.tileCoordsTouched)))
	return output.String()
}

func (l *lab) moveGuard() {
	nextGuardPos := coord{l.guard.row, l.guard.col}
	nextGuardPos.add(l.guardDir)
	if l.isValid(nextGuardPos) && l.roomTiles[nextGuardPos.row][nextGuardPos.col] == "#" {
		// rotation matrix clockwise 90°
		// x1 = y0 --> x = row; y = col
		// y1 = -x0
		l.guardDir.col, l.guardDir.row = -1*l.guardDir.row, l.guardDir.col
	} else {
		l.guard = nextGuardPos
	}
}

func (l *lab) simulateGuard(abortAfterIterations int, visualize bool) (touchedTiles int, aborted bool) {
	l.tileCoordsTouched = make(map[coord]bool)
	iterations := 0
	for l.isValid(l.guard) && iterations < abortAfterIterations {
		l.print(visualize)
		if visualize {
			time.Sleep(50 * time.Nanosecond)
		}
		l.tileCoordsTouched[l.guard] = true
		l.moveGuard()
		iterations++
	}
	l.print(visualize)
	return len(l.tileCoordsTouched), iterations >= abortAfterIterations
}

func (l *lab) fakeObstacles(abortAfterIterations int, visualize bool) (fakedObstacle int) {
	fakedObstacle = 0
	guardInitPos := coord{l.guard.row, l.guard.col}
	guardInitDir := coord{l.guardDir.row, l.guardDir.col}
	for i, _ := range l.roomTiles {
		for j, _ := range l.roomTiles[i] {
			l.guard = guardInitPos
			l.guardDir = guardInitDir
			isOnGuardPos := l.guard.row == i && l.guard.col == j
			isOnGuardNextStep := l.guard.row+l.guardDir.row == i && l.guard.col+l.guardDir.col == j
			if l.roomTiles[i][j] == "." && !isOnGuardPos && !isOnGuardNextStep {
				l.roomTiles[i][j] = "#"
				_, aborted := l.simulateGuard(abortAfterIterations, visualize)
				if aborted {
					fakedObstacle++
				}
				l.roomTiles[i][j] = "."
			}
		}
	}

	return fakedObstacle
}

func main() {
	start := time.Now()

	inputLines := readFile("day6/input.txt")
	laboratory := parseLab(inputLines)

	const abortAfterIterations = 6666
	tiles, aborted := laboratory.simulateGuard(abortAfterIterations, false)
	fmt.Println("Part 1:", tiles, aborted)

	laboratory = parseLab(inputLines)
	fakedObstacles := laboratory.fakeObstacles(abortAfterIterations, false)
	fmt.Println("Part 2:", fakedObstacles)

	fmt.Println("Finished in", time.Since(start))

	g := &game{
		laboratory: parseLab(inputLines),
	}
	ebiten.SetWindowTitle("Part1")
	ebiten.SetFullscreen(true)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
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
