package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type tree struct {
	isEdge      bool
	height      int
	visible     bool
	scenicScore int
	i           int
	j           int
}

var input [][]int
var trees []*tree

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	splittedLines := strings.Split(strings.ReplaceAll(lines, "\r\n", "\n"), "\n")
	input = make([][]int, len(splittedLines)-1)
	for i, line := range splittedLines {
		treeHeights := strings.Split(line, "")
		if len(treeHeights) < 2 {
			return
		}
		input[i] = make([]int, len(treeHeights))
		for j, strHeight := range treeHeights {
			height, err := strconv.Atoi(strHeight)
			if err != nil {
				return
			}
			input[i][j] = height
		}
	}
}

func populateTrees() {
	for i, heights := range input {
		for j, height := range heights {
			tree := newTree(i, j, height)
			trees = append(trees, tree)
		}
	}
}

func newTree(i int, j int, height int) *tree {
	tree := tree{
		isEdge:      i == 0 || j == 0 || i == len(input)-1 || j == len(input[0])-1,
		height:      height,
		visible:     false,
		i:           i,
		j:           j,
		scenicScore: 1,
	}
	tree.checkVisibility()
	tree.calculateScenicScore()
	return &tree
}

func (tree *tree) checkVisibility() {
	if tree.isEdge {
		tree.visible = true
		return
	}

	tree.visible = false

	var hasHigher bool

	hasHigher = false
	for i := 0; i < tree.i; i++ {
		if input[i][tree.j] >= tree.height {
			hasHigher = true
			break
		}
	}
	if !hasHigher {
		tree.visible = true
		return
	}

	hasHigher = false
	for i := tree.i + 1; i < len(input); i++ {
		if input[i][tree.j] >= tree.height {
			hasHigher = true
			break
		}
	}
	if !hasHigher {
		tree.visible = true
		return
	}

	hasHigher = false
	for j := 0; j < tree.j; j++ {
		if input[tree.i][j] >= tree.height {
			hasHigher = true
			break
		}
	}
	if !hasHigher {
		tree.visible = true
		return
	}

	hasHigher = false
	for j := tree.j + 1; j < len(input[0]); j++ {
		if input[tree.i][j] >= tree.height {
			hasHigher = true
			break
		}
	}
	if !hasHigher {
		tree.visible = true
		return
	}
}

func countVisibleTrees() int {
	sum := 0
	for _, tree := range trees {
		if tree.visible {
			sum++
		}
	}
	return sum
}

func (tree *tree) calculateScenicScore() {
	if tree.isEdge {
		tree.scenicScore = 0
		return
	}

	var momentaryScore int

	momentaryScore = 0
	for i := tree.i + 1; i < len(input); i++ {
		momentaryScore++
		if input[i][tree.j] >= tree.height {
			break
		}
	}
	tree.scenicScore *= momentaryScore

	momentaryScore = 0
	for i := tree.i - 1; i >= 0; i-- {
		momentaryScore++
		if input[i][tree.j] >= tree.height {
			break
		}
	}
	tree.scenicScore *= momentaryScore

	momentaryScore = 0
	for j := tree.j + 1; j < len(input); j++ {
		momentaryScore++
		if input[tree.i][j] >= tree.height {
			break
		}
	}
	tree.scenicScore *= momentaryScore

	momentaryScore = 0
	for j := tree.j - 1; j >= 0; j-- {
		momentaryScore++
		if input[tree.i][j] >= tree.height {
			break
		}
	}
	tree.scenicScore *= momentaryScore
}

func findHighestScenicScoreTree() int {
	max := 0
	for _, tree := range trees {
		if tree.scenicScore > max {
			max = tree.scenicScore
		}
	}
	return max
}

func main() {
	readFile()
	populateTrees()
	fmt.Println("Part01: ", countVisibleTrees())
	fmt.Println("Part02: ", findHighestScenicScoreTree())
}
