package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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

type game struct {
	room          *room
	secondsToPass *int
	secondsPassed *int
	roomStrings   map[int]string
}

func (g game) Update() error {
	_, dy := ebiten.Wheel()
	delta := int(dy)
	pressedKeys := inpututil.AppendJustPressedKeys([]ebiten.Key{})
	for _, key := range pressedKeys {
		if key == ebiten.KeyArrowLeft {
			delta -= g.room.tilesY
		}
		if key == ebiten.KeyArrowRight {
			delta += g.room.tilesY
		}
		if key == ebiten.KeyA {
			delta -= g.room.tilesX
		}
		if key == ebiten.KeyD {
			delta += g.room.tilesX
		}
	}
	*g.secondsPassed += delta
	g.room.moveRobots(delta)
	return nil
}

func (g game) Draw(screen *ebiten.Image) {
	const size = 30
	const sizeSmall = 24
	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f | Arrow Keys: +-%d Seconds | A & D: +-%d Seconds | Mouse Wheel: +-1 Second", ebiten.ActualTPS(), g.room.tilesY, g.room.tilesX)
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   sizeSmall,
	}, op)

	// Draw the game
	op = &text.DrawOptions{}
	op.GeoM.Translate(0, 30)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = size * 1.01
	s := g.room.String() + "| Second:" + fmt.Sprint(*g.secondsPassed)

	if _, ok := g.roomStrings[*g.secondsPassed]; !ok {
		g.roomStrings[*g.secondsPassed] = s
	}
	text.Draw(screen, s, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   size,
	}, op)
}

func (g game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//TODO implement me
	size := 32
	return g.room.tilesY * size, g.room.tilesX * size
}

type room struct {
	tilesX int
	tilesY int
	robots []*robot
}

func (r *room) String() string {
	sb := make([][]string, r.tilesY)
	for i := 0; i < r.tilesY; i++ {
		line := make([]string, r.tilesX)
		for j := 0; j < r.tilesX; j++ {
			line[j] = "."
		}
		sb[i] = line
	}
	robotMap := make(map[coordinate]int)
	for _, robot := range r.robots {
		robotMap[robot.position]++
	}
	for c, i := range robotMap {
		sb[c.y][c.x] = strconv.Itoa(i)
	}

	stringBuilder := strings.Builder{}
	for _, line := range sb {
		for _, s := range line {
			stringBuilder.WriteString(s)
		}
		stringBuilder.WriteString("\n")
	}
	stringBuilder.WriteString(fmt.Sprintf("Quadrant Count: %v | Safety Factor: %d ", r.countPerQuadrant(), r.safetyFactor()))
	return stringBuilder.String()
}

func (r *room) moveRobots(seconds int) {
	for _, robot := range r.robots {
		robot.move(seconds, r.tilesX, r.tilesY)
	}
}

func (r room) countPerQuadrant() []int {
	quadrant := make([]int, 4)
	quadrantX := r.tilesX / 2
	quadrantY := r.tilesY / 2
	for _, rob := range r.robots {
		if rob.position.x == quadrantX || rob.position.y == quadrantY {
			continue
		}
		quadIndx := 3
		if rob.position.x < quadrantX {
			quadIndx--
		}
		if rob.position.y < quadrantY {
			quadIndx -= 2
		}
		quadrant[quadIndx]++
	}
	return quadrant
}

func (r *room) safetyFactor() int {
	quadrant := r.countPerQuadrant()
	mul := 1
	for _, i := range quadrant {
		mul *= i
	}
	return mul
}

type robot struct {
	velocity coordinate
	position coordinate
}

func (r *robot) move(seconds, xMax, yMax int) {
	posX := r.position.x + r.velocity.x*seconds
	posY := r.position.y + r.velocity.y*seconds
	posX %= xMax
	posY %= yMax
	if posX < 0 {
		posX += xMax
	}
	if posY < 0 {
		posY += yMax
	}
	r.position.x = posX
	r.position.y = posY

}

type coordinate struct {
	x, y int
}

func (c coordinate) String() string {
	return fmt.Sprintf("(%d,%d)", c.x, c.y)
}

func main() {
	start := time.Now()

	room := readFile("day14/input.txt", 101, 103)
	room.moveRobots(100)
	fmt.Println(room.String())
	fmt.Println("Part 1:", room.safetyFactor())

	fmt.Println("Finished in", time.Since(start))

	secondsToPass := new(int)
	*secondsToPass = 100
	g := &game{
		room:          readFile("day14/input.txt", 101, 103),
		secondsToPass: secondsToPass,
		secondsPassed: new(int),
		roomStrings:   make(map[int]string),
	}
	ebiten.SetWindowTitle("AoC Day 14")
	ebiten.SetFullscreen(true)
	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}

func readFile(file string, sizeX, sizeY int) *room {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	var robots []*robot
	for _, robotString := range split {
		robotArgs := strings.Split(robotString, " v=")
		positionArgs := strings.Split(strings.TrimLeft(robotArgs[0], "p="), ",")
		posX, _ := strconv.Atoi(positionArgs[0])
		posY, _ := strconv.Atoi(positionArgs[1])

		velocityArgs := strings.Split(robotArgs[1], ",")
		velocityX, _ := strconv.Atoi(velocityArgs[0])
		velocityY, _ := strconv.Atoi(velocityArgs[1])

		robots = append(robots, &robot{
			velocity: coordinate{velocityX, velocityY},
			position: coordinate{posX, posY},
		})
	}
	return &room{
		tilesX: sizeX,
		tilesY: sizeY,
		robots: robots,
	}
}
