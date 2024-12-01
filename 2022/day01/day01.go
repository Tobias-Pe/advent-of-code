package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type elf struct {
	calories []int
}

var elfs []elf

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)

	currentElf := elf{}
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			elfs = append(elfs, currentElf)
			currentElf = elf{}
		} else {
			calories, err := strconv.Atoi(line)

			if err != nil {
				fmt.Println("Error during conversion")
				return
			}

			currentElf.calories = append(currentElf.calories, calories)
		}
	}
	elfs = append(elfs, currentElf)

	fmt.Println(elfs)
}

func main() {
	readFile()
}
