package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type pairEntity interface {
	compareTo(entity pairEntity) int
}

type number struct {
	value int
}

func (left number) compareTo(right pairEntity) int {
	rightNumber, ok := right.(*number)
	if ok {
		fmt.Println("Compare", left.value, rightNumber.value)
		if rightNumber.value == left.value {
			return 0
		} else if rightNumber.value > left.value {
			fmt.Println("Left side is smaller, so inputs are in the right order")
			return 1
		} else {
			fmt.Println("Right side is smaller, so inputs are not in the right order")
			return -1
		}
	} else {
		rightList := right.(*list)
		return list{pairEntities: []pairEntity{left}}.compareTo(rightList)
	}
}

type list struct {
	pairEntities []pairEntity
	parent       *list
	isDivider    bool
}

func (left list) compareTo(right pairEntity) int {
	rightNumber, ok := right.(*number)
	if ok {
		return left.compareTo(&list{pairEntities: []pairEntity{rightNumber}})
	} else {
		rightList := right.(*list)
		fmt.Println("Compare", left.pairEntities, rightList.pairEntities)
		for i, pairEntity := range left.pairEntities {
			if i >= len(rightList.pairEntities) {
				fmt.Println("Right side ran out of items, so inputs are not in the right order")
				return -1
			}
			compareTo := pairEntity.compareTo(rightList.pairEntities[i])
			if compareTo != 0 {
				return compareTo
			}
		}
		if len(left.pairEntities) == len(rightList.pairEntities) {
			return 0
		}
		fmt.Println("Left side ran out of items, so inputs are in the right order")
		return 1
	}
}

type pair struct {
	left  *list
	right *list
}

var pairs []*pair
var lists []*list

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	splittedPairLines := strings.Split(lines, "\n\n")
	for _, line := range splittedPairLines {
		strPairs := strings.Split(line, "\n")
		p := &pair{
			left: &list{
				pairEntities: []pairEntity{},
				parent:       nil,
			},
			right: &list{
				pairEntities: []pairEntity{},
				parent:       nil,
			},
		}

		var currentList *list
		for i, pair := range strPairs {
			if i == 0 {
				currentList = p.left
			} else {
				currentList = p.right
			}
			currentLine := strings.Split(pair, "")
			currentLine = append(currentLine[:0], currentLine[1:]...)
			currentLine = currentLine[:len(currentLine)-1]
			indexLetter := 0
			for indexLetter < len(currentLine) {
				letter := currentLine[indexLetter]
				if letter == "[" {
					newList := &list{
						pairEntities: []pairEntity{},
						parent:       currentList,
					}
					currentList.pairEntities = append(currentList.pairEntities, newList)
					currentList = newList
				} else if letter == "]" {
					currentList = currentList.parent
				} else if unicode.IsNumber(rune(letter[0])) {
					strNumber := ""
					for index := indexLetter; index < len(currentLine) && unicode.IsNumber(rune(currentLine[index][0])); index++ {
						strNumber += currentLine[index]
					}
					indexLetter += len(strNumber) - 1
					atoi, err := strconv.Atoi(strNumber)
					if err != nil {
						fmt.Println(err)
						return
					}
					currentList.pairEntities = append(currentList.pairEntities, &number{value: atoi})
				}
				indexLetter++
			}
		}
		pairs = append(pairs, p)
		lists = append(lists, p.left)
		lists = append(lists, p.right)
	}
	divider1 := &list{
		pairEntities: []pairEntity{},
		parent:       nil,
		isDivider:    true,
	}
	divider1.pairEntities = append(divider1.pairEntities, &list{
		pairEntities: []pairEntity{&number{value: 2}},
		parent:       divider1,
	})
	lists = append(lists, divider1)
	divider2 := &list{
		pairEntities: []pairEntity{},
		parent:       nil,
		isDivider:    true,
	}
	divider2.pairEntities = append(divider2.pairEntities, &list{
		pairEntities: []pairEntity{&number{value: 6}},
		parent:       divider2,
	})
	lists = append(lists, divider2)
}

func sortLists() {
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].compareTo(lists[j]) == 1
	})
}

func main() {
	readFile("input.txt")
	sum := 0
	for i, p := range pairs {
		fmt.Println("\n\n=== Pair", i+1, " ===")
		if p.left.compareTo(p.right) == 1 {
			fmt.Println("Right Order at: ", i+1)
			sum += i + 1
		}
	}
	fmt.Println(sum)

	sortLists()

	mul := 1
	for i, list := range lists {
		if list.isDivider {
			fmt.Println(i)
			mul *= i + 1
		}
	}
	fmt.Println(mul)
}
