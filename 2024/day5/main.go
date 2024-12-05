package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type page int
type pageSet map[page]bool

type pageGraph struct {
	pagesMap map[page]pageSet
}

func createPageGraph(lines []string) *pageGraph {
	pg := &pageGraph{}
	pg.pagesMap = make(map[page]pageSet)
	for _, line := range lines {
		pg.addPage(line)
	}
	return pg
}

func (pg *pageGraph) addPage(line string) {
	split := strings.Split(strings.TrimSpace(line), "|")
	from, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	to, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}
	pg.connectPages(page(from), page(to))
}

func (pg *pageGraph) addPageNum(pageNum page) {
	if _, ok := pg.pagesMap[pageNum]; !ok {
		pg.pagesMap[pageNum] = make(pageSet)
	}
}

func (pg *pageGraph) connectPages(from, to page) {
	pg.addPageNum(from)
	pg.addPageNum(to)
	pg.pagesMap[from][to] = true
}

func main() {
	start := time.Now()

	linesGraph := readFile("day5/input_connections.txt")
	graph := createPageGraph(linesGraph)

	linesPrintQueue := readFile("day5/input_pages.txt")
	pagesList := createPagesList(linesPrintQueue)

	validPages := filterInvalidPages(graph, pagesList)
	fmt.Println("Part 1:", sumMiddlePage(validPages), len(validPages))

	validPages = getFixedInvalidPages(graph, pagesList)
	fmt.Println("Part 2:", sumMiddlePage(validPages), len(validPages))

	fmt.Println("Finished in", time.Since(start))
}

func sumMiddlePage(listPages [][]page) int {
	sum := 0
	for _, pages := range listPages {
		sum += int(pages[len(pages)/2])
	}
	return sum
}

func filterInvalidPages(graph *pageGraph, list [][]page) [][]page {
	validPages := [][]page{}

	for _, pages := range list {
		if isValidPageList(graph, pages) {
			validPages = append(validPages, pages)
		}
	}

	return validPages
}

func getFixedInvalidPages(graph *pageGraph, list [][]page) [][]page {
	fixedPages := [][]page{}

	for _, pages := range list {
		fixedList, fixApplied := getFixedInvalidPage(graph, pages)
		if fixApplied {
			fixedPages = append(fixedPages, *fixedList)
		}
	}

	return fixedPages
}

func isValidPageList(graph *pageGraph, pages []page) bool {
	for i := 0; i < len(pages); i++ {
		// skip page if it is not in graph -> can occur anywhere
		followingPages, ok := graph.pagesMap[pages[i]]
		if !ok {
			continue
		}
		for _, prevPage := range pages[:i] {
			if followingPages[prevPage] {
				return false
			}
		}
	}
	return true
}

func getFixedInvalidPage(graph *pageGraph, list []page) (*[]page, bool) {
	fixedList := append([]page{}, list...)
	fixApplied := false
	for i := 0; i < len(fixedList); i++ {
		followingPages, ok := graph.pagesMap[fixedList[i]]
		if !ok {
			continue
		}
		for j, prevPage := range fixedList[:i] {
			if followingPages[prevPage] {
				fixedList[i], fixedList[j] = fixedList[j], fixedList[i]
				i--
				fixApplied = true
				break
			}
		}
	}
	return &fixedList, fixApplied
}

func createPagesList(linesPrintQueue []string) [][]page {
	var pagesList [][]page
	for _, line := range linesPrintQueue {
		stringPages := strings.Split(strings.TrimSpace(line), ",")
		var pages []page
		for _, stringPage := range stringPages {
			pageNum, err := strconv.Atoi(stringPage)
			if err != nil {
				panic(err)
			}
			pages = append(pages, page(pageNum))
		}
		pagesList = append(pagesList, pages)
	}
	return pagesList
}

func readFile(file string) []string {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	return split
}
