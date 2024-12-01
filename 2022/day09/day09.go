package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type command struct {
	direction string
	amount    int
}

type position struct {
	x int
	y int
}

type head struct {
	position position
}

func (head *head) move(direction string, amount int) position {
	oldPosition := head.position
	switch direction {
	case "R":
		head.position.x += amount
	case "L":
		head.position.x -= amount
	case "U":
		head.position.y += amount
	case "D":
		head.position.y -= amount
	}
	return oldPosition
}

type tail struct {
	position position
}

func (t tail) isClose(p position) bool {
	return t.position.x+1 == p.x && t.position.y-1 == p.y ||
		t.position.x+1 == p.x && t.position.y == p.y ||
		t.position.x+1 == p.x && t.position.y+1 == p.y ||
		t.position.x == p.x && t.position.y-1 == p.y ||
		t.position.x == p.x && t.position.y == p.y ||
		t.position.x == p.x && t.position.y+1 == p.y ||
		t.position.x-1 == p.x && t.position.y-1 == p.y ||
		t.position.x-1 == p.x && t.position.y == p.y ||
		t.position.x-1 == p.x && t.position.y+1 == p.y
}

func (t *tail) moveTo(headPosition position) position {
	var p = t.position

	if t.position.x < headPosition.x {
		t.position.x += 1
	} else if t.position.x > headPosition.x {
		t.position.x -= 1
	}

	if t.position.y < headPosition.y {
		t.position.y += 1
	} else if t.position.y > headPosition.y {
		t.position.y -= 1
	}

	return p
}

var mapTailPos map[position]bool
var commands []command
var Head head
var Tail tail

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	splittedLines := strings.Split(strings.ReplaceAll(lines, "\r\n", "\n"), "\n")
	for _, line := range splittedLines {
		stringCommand := strings.Split(line, " ")
		amount, err := strconv.Atoi(stringCommand[1])
		if err != nil {
			return
		}

		commands = append(commands, command{
			direction: stringCommand[0],
			amount:    amount,
		})
	}

}
func executeCommands() {
	for _, command := range commands {
		for i := 0; i < command.amount; i++ {
			oldHeadPosition := Head.move(command.direction, 1)
			if !Tail.isClose(Head.position) {
				Tail.moveTo(oldHeadPosition)
				mapTailPos[Tail.position] = true
			}
		}
	}
}

func executeCommandsOn10KnotRope() {
	var tails []*tail
	for i := 0; i < 9; i++ {
		tails = append(tails, &tail{position: position{
			x: 0,
			y: 0,
		}})
	}

	for _, command := range commands {
		for i := 0; i < command.amount; i++ {
			Head.move(command.direction, 1)
			prevNodePos := Head.position
			for i, node := range tails {
				if !node.isClose(prevNodePos) {
					node.moveTo(prevNodePos)
					prevNodePos = node.position
					if i == 8 {
						mapTailPos[node.position] = true
					}
				} else {
					break
				}
			}
		}
	}
}

func main() {
	readFile()
	Head = head{position: position{
		x: 0,
		y: 0,
	}}
	Tail = tail{position: position{
		x: 0,
		y: 0,
	}}
	mapTailPos = map[position]bool{}
	mapTailPos[Tail.position] = true
	executeCommands()
	fmt.Println("Part01: ", len(mapTailPos))
	Head = head{position: position{
		x: 0,
		y: 0,
	}}
	mapTailPos = map[position]bool{}
	mapTailPos[Head.position] = true
	executeCommandsOn10KnotRope()
	fmt.Println("Part02: ", len(mapTailPos))
}
