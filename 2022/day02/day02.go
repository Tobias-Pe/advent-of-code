package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type gameConstellation struct {
	opponentsChoice string
	myChoice        string
	points          int
	state           string
}

var gameConstellations []gameConstellation

func readFilePart01() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		choices := strings.Split(line, " ")
		if len(choices) != 2 {
			return
		}
		gameConstellation := gameConstellation{opponentsChoice: choices[0], myChoice: choices[1]}
		calcPointsPart01(&gameConstellation)
		gameConstellations = append(gameConstellations, gameConstellation)
	}
}

func calcPointsPart01(gameConstellation *gameConstellation) {
	if gameConstellation.myChoice == "Z" && gameConstellation.opponentsChoice == "C" ||
		gameConstellation.myChoice == "Y" && gameConstellation.opponentsChoice == "B" ||
		gameConstellation.myChoice == "X" && gameConstellation.opponentsChoice == "A" {
		gameConstellation.state = "draw"
		gameConstellation.points = 3
	} else if gameConstellation.myChoice == "Z" && gameConstellation.opponentsChoice == "A" ||
		gameConstellation.myChoice == "Y" && gameConstellation.opponentsChoice == "C" ||
		gameConstellation.myChoice == "X" && gameConstellation.opponentsChoice == "B" {
		gameConstellation.state = "loose"
		gameConstellation.points = 0
	} else {
		gameConstellation.state = "win"
		gameConstellation.points = 6
	}

	switch gameConstellation.myChoice {
	case "X":
		gameConstellation.points += 1
	case "Y":
		gameConstellation.points += 2
	case "Z":
		gameConstellation.points += 3
	}

	fmt.Println(gameConstellation.state, gameConstellation.points, gameConstellation.opponentsChoice, gameConstellation.myChoice)
}

func readFilePart02() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		choices := strings.Split(line, " ")
		if len(choices) != 2 {
			return
		}
		gameConstellation := gameConstellation{opponentsChoice: choices[0], myChoice: choices[1]}
		calcPointsPart02(&gameConstellation)
		gameConstellations = append(gameConstellations, gameConstellation)
	}
}

func calcPointsPart02(gameConstellation *gameConstellation) {
	if gameConstellation.myChoice == "X" { // loose
		gameConstellation.points = 0
		switch gameConstellation.opponentsChoice {
		case "A":
			gameConstellation.points += 3
		case "B":
			gameConstellation.points += 1
		case "C":
			gameConstellation.points += 2
		}
	} else if gameConstellation.myChoice == "Y" { // draw
		gameConstellation.points = 3
		switch gameConstellation.opponentsChoice {
		case "A":
			gameConstellation.points += 1
		case "B":
			gameConstellation.points += 2
		case "C":
			gameConstellation.points += 3
		}
	} else {
		gameConstellation.points = 6
		switch gameConstellation.opponentsChoice { // win
		case "A":
			gameConstellation.points += 2
		case "B":
			gameConstellation.points += 3
		case "C":
			gameConstellation.points += 1
		}
	}

	fmt.Println(gameConstellation.state, gameConstellation.points, gameConstellation.opponentsChoice, gameConstellation.myChoice)
}

func main() {
	readFilePart02()

	sum := 0
	for _, constellation := range gameConstellations {
		sum += constellation.points
	}
	fmt.Println(sum)
}
