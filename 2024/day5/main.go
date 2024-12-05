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
	startingNodes pageSet
	pagesMap      map[page]pageSet
}

func createPageGraph(lines []string) *pageGraph {
	pg := &pageGraph{}
	pg.startingNodes = make(pageSet)
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

func (pg *pageGraph) addPageNum(pageNum page, isFrom bool) {
	if _, ok := pg.pagesMap[pageNum]; !ok {
		pg.pagesMap[pageNum] = make(pageSet)
		if isFrom {
			pg.startingNodes[pageNum] = true
		}
	}
	if !isFrom {
		delete(pg.startingNodes, pageNum)
	}
}

func (pg *pageGraph) connectPages(from, to page) {
	pg.addPageNum(from, true)
	pg.addPageNum(to, false)
	pg.pagesMap[from][to] = true
}

func (pg *pageGraph) dfs(current, end page) []page {
	if current == end {
		return []page{}
	}
	nextPages, ok := pg.pagesMap[current]
	if !ok {
		return nil
	}
	for nextPage := range nextPages {
		history := pg.dfs(nextPage, end)
		if history != nil {
			history = append(history, nextPage)
			return history
		}
	}
	return nil
}

func main() {
	start := time.Now()

	linesGraph := readFile("day5/input_connections.txt")
	graph := createPageGraph(linesGraph)

	linesPrintQueue := readFile("day5/input_pages.txt")
	pagesList := createPagesList(linesPrintQueue)

	filteredPagesList := filterInvalidPages(graph, pagesList)

	fmt.Println("Finished in", time.Since(start))
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

func isValidPageList(graph *pageGraph, pages []page) bool {
	invalidPages := make(pageSet)
	var lastRelevantPage page = -1
	for _, page := range pages {
		// invalid page occurance
		if _, match := invalidPages[page]; match {
			return false
		}
		// skip page if it is not in graph -> can occure anywhere
		if _, ok := graph.pagesMap[page]; !ok {
			continue
		}
		if _, ok := graph.startingNodes[page]; ok {
			lastRelevantPage = page
		}
		if lastRelevantPage == -1 {
			lastRelevantPage = page
		} else {
			history := graph.dfs(lastRelevantPage, page)
			if history == nil {
				return false
			}

		}
		invalidPages[page] = true
	}
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
