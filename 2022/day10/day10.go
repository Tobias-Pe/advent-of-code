package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type command struct {
	isNoop     bool
	amount     int
	cyclesLeft int
}

var registerX int

var statesX []int

var commands []*command

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	splittedLines := strings.Split(strings.ReplaceAll(lines, "\r\n", "\n"), "\n")
	for _, line := range splittedLines {
		stringCommand := strings.Split(line, " ")
		if len(stringCommand) < 2 {
			commands = append(commands, &command{
				isNoop:     true,
				amount:     0,
				cyclesLeft: 1,
			})
		} else {
			amount, err := strconv.Atoi(stringCommand[1])
			if err != nil {
				return
			}

			commands = append(commands, &command{
				isNoop:     false,
				amount:     amount,
				cyclesLeft: 2,
			})
		}
	}
}

func executeCommands() {
	registerX = 1

	for i := 0; i < len(commands); {
		//up
		commands[i].cyclesLeft--

		//during
		statesX = append(statesX, registerX)

		//down
		if commands[i].cyclesLeft == 0 {
			registerX += commands[i].amount
			i++
		}
	}
}

func calcSignalStrength() int {
	sum := 0
	nextCycle := 20
	for i, x := range statesX {
		if i+1 == nextCycle {
			sum += x * nextCycle
			nextCycle += 40
		}
	}
	return sum
}

func printLogic() {
	output := "#"
	for i, posSprite := range statesX {
		if i%40 == 0 {
			output += "\n#"
		} else {
			if posSprite == i%40 || posSprite-1 == i%40 || posSprite+1 == i%40 {
				output += "#"
			} else {
				output += "."
			}
		}
	}
	fmt.Println(output)
}

func main() {
	readFile("./input.txt")
	executeCommands()
	fmt.Println("Part01: ", calcSignalStrength())
	printLogic()
}
