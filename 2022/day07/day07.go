package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type data interface {
	getSize() int
	getName() string
}

type file struct {
	size int
	name string
}

func (file file) getSize() int {
	return file.size
}

func (file file) getName() string {
	return file.name
}

type directory struct {
	parent  *directory
	name    string
	content []data
}

func (directory directory) listContent() string {
	names := ""
	for _, data := range directory.content {
		names += data.getName()
		names += ", "
	}
	return names
}

func newDirectory(parent *directory, name string) directory {
	return directory{
		parent:  parent,
		name:    name,
		content: []data{},
	}
}

func (directory directory) getSize() int {
	sum := 0
	for _, data := range directory.content {
		sum += data.getSize()
	}
	return sum
}

func (directory directory) getName() string {
	return directory.name
}

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	topLevel = newDirectory(&topLevel, "/")
	currentDir := &topLevel
	for _, line := range strings.Split(strings.ReplaceAll(lines, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "$ cd") {
			targetForDir := strings.TrimPrefix(line, "$ cd ")
			currentDir = changeDirectory(currentDir, targetForDir)
		} else if !strings.HasPrefix(line, "$") && len(line) > 1 {
			args := strings.Split(line, " ")
			populateDirectory(currentDir, args[0], args[1])
		}
	}
}

func populateDirectory(current *directory, argument string, name string) {
	if argument == "dir" {
		newDir := newDirectory(current, name)
		current.content = append(current.content, &newDir)
	} else {
		convertedSize, err := strconv.Atoi(argument)
		if err != nil {
			fmt.Println(err)
			return
		}
		current.content = append(current.content, &file{
			size: convertedSize,
			name: name,
		})
	}
}

func changeDirectory(current *directory, target string) *directory {
	switch target {
	case "/":
		return &topLevel
	case "..":
		return current.parent
	default:
		for _, data := range current.content {
			if directory, isDirectory := data.(*directory); isDirectory && target == directory.name {
				return directory
			}
		}
		fmt.Println("No Directory ", target, " found :(")
	}
	return nil
}

func calcSumDirectoriesSizeWithAtMost100000() int {
	var sumOfAllDirectories = 0
	sumDirectoriesSizeWithAtMost(&topLevel, 100000, &sumOfAllDirectories)
	return sumOfAllDirectories
}

func sumDirectoriesSizeWithAtMost(dir *directory, max int, sumOfAllDirectories *int) int {
	sum := 0
	for _, data := range dir.content {
		if directory, isDirectory := data.(*directory); isDirectory {
			sum += sumDirectoriesSizeWithAtMost(directory, max, sumOfAllDirectories)
		} else {
			sum += data.getSize()
		}
	}
	if sum <= max {
		*sumOfAllDirectories += sum

	}
	return sum
}

type pair struct {
	name string
	size int
}

func calcSizeOfAllDirs() []pair {
	var sizes []pair
	fmt.Println("Toplevel size: ", getSizeOfAllChildDirs(&topLevel, &sizes))
	return sizes
}

func getSizeOfAllChildDirs(dir *directory, sizes *[]pair) int {
	sum := 0
	for _, data := range dir.content {
		if directory, isDirectory := data.(*directory); isDirectory {
			sum += getSizeOfAllChildDirs(directory, sizes)
		} else {
			sum += data.getSize()
		}
	}
	*sizes = append(*sizes, pair{
		name: dir.name,
		size: sum,
	})

	return sum
}

var topLevel directory

func main() {
	readFile()
	fmt.Println("Part01: ", calcSumDirectoriesSizeWithAtMost100000())
	sizesOfAllDirs := calcSizeOfAllDirs()
	sort.Slice(sizesOfAllDirs, func(i, j int) bool {
		return sizesOfAllDirs[i].size > sizesOfAllDirs[j].size
	})
	fmt.Println(sizesOfAllDirs)
}
