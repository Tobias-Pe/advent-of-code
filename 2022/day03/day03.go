package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type rucksack struct {
	items          []string
	compartments   [2][]string
	wrongItemsSet  map[string]bool
	badge          string
	prioBadge      int
	prioWrongItems int
}

var rucksacks []*rucksack

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	for _, line := range strings.Split(lines, "\n") {
		characters := strings.Split(line, "")
		if len(characters) > 0 {
			var compartments [2][]string
			compartments[0] = characters[:len(characters)/2]
			compartments[1] = characters[len(characters)/2:]
			rucksack := rucksack{compartments: compartments, wrongItemsSet: map[string]bool{}, items: characters}
			rucksacks = append(rucksacks, &rucksack)
		}
	}
}

func assignWrongItems() {
	for _, rucksack := range rucksacks {
		itemSet := make(map[string]bool)
		for _, item := range rucksack.compartments[0] {
			itemSet[item] = true
		}

		for _, item := range rucksack.compartments[1] {
			if _, ok := itemSet[item]; ok {
				rucksack.wrongItemsSet[item] = true
			}
		}
	}
}

func calcPrioWrongItems() {
	for _, rucksack := range rucksacks {
		for key, _ := range rucksack.wrongItemsSet {
			prio := 0
			if strings.ToUpper(key) == key {
				key = strings.ToLower(key)
				prio += 26
			}
			prio += int(key[0]) - 96
			rucksack.prioWrongItems += prio
		}
	}
}

func assignBadge() int {
	prio := 0
	for firstOfThreeBackpacks := 0; firstOfThreeBackpacks < len(rucksacks); firstOfThreeBackpacks += 3 {
		itemSet := make(map[string]bool)
		for _, item := range rucksacks[firstOfThreeBackpacks].items {
			itemSet[item] = false
		}
		for _, item := range rucksacks[firstOfThreeBackpacks+1].items {
			if _, ok := itemSet[item]; ok {
				itemSet[item] = true
			}
		}
		for _, item := range rucksacks[firstOfThreeBackpacks+2].items {
			if _, ok := itemSet[item]; ok && itemSet[item] {
				if strings.ToUpper(item) == item {
					item = strings.ToLower(item)
					prio += 26
				}
				prio += int(item[0]) - 96
				break
			}
		}
	}
	return prio
}

func main() {
	readFile()
	assignWrongItems()
	calcPrioWrongItems()
	sum := 0
	for _, r := range rucksacks {
		sum += r.prioWrongItems
	}
	fmt.Println("Part 01: ", sum)
	fmt.Println("Part 02: ", assignBadge())
}
