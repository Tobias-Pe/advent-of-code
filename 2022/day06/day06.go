package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Queue []string

func (q *Queue) push(char string) {
	*q = append(*q, char)
}

func (q *Queue) pop() (string, bool) {
	if len(*q) == 0 {
		return "", false
	}
	char := (*q)[0]
	*q = (*q)[1:]
	return char, true
}

var characters []string

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	allLines := string(content)
	lines := strings.Split(allLines, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			characters = strings.Split(line, "")
		}
	}
}

func findFirstMarker() int {
	markerCharacterMap := make(map[string]bool)
	var marker Queue
	for i, character := range characters {
		if hasCharacter := markerCharacterMap[character]; hasCharacter {
			for char, wasPopped := marker.pop(); wasPopped && char != character; char, wasPopped = marker.pop() {
				delete(markerCharacterMap, char)
			}
		}
		markerCharacterMap[character] = true
		marker.push(character)

		if len(marker) == 4 {
			return i + 1
		}
	}
	return -1
}

func findFirstMessage() int {
	markerCharacterMap := make(map[string]bool)
	var marker Queue
	for i, character := range characters {
		if hasCharacter := markerCharacterMap[character]; hasCharacter {
			for char, wasPopped := marker.pop(); wasPopped && char != character; char, wasPopped = marker.pop() {
				delete(markerCharacterMap, char)
			}
		}
		markerCharacterMap[character] = true
		marker.push(character)

		if len(marker) == 14 {
			return i + 1
		}
	}
	return -1
}

func main() {
	readFile()
	fmt.Println(findFirstMarker())
	fmt.Println(findFirstMessage())
}
