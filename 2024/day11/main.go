package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type stoneArrangement map[stone]int

func createStoneArrangement(input string) stoneArrangement {
	stoneArrangement := make(stoneArrangement)
	for _, s := range strings.Fields(input) {
		stoneArrangement[createStone(s)] += 1
	}
	return stoneArrangement
}

func (stA *stoneArrangement) blink() {
	newSta := stoneArrangement{}
	for st, mul := range *stA {
		newStone, isPresent := st.blink()
		newSta[st] += mul
		if isPresent {
			newSta[newStone] += mul
		}
	}
	*stA = newSta
}

func (stA *stoneArrangement) blinkTimes(times int) {
	for i := 0; i < times; i++ {
		stA.blink()
		//fmt.Println(i, stA)
	}
}

func (stA stoneArrangement) count() int {
	sum := 0
	for _, count := range stA {
		sum += count
	}
	return sum
}

func (stA stoneArrangement) String() string {
	sB := strings.Builder{}
	for st, mul := range stA {
		sB.WriteString(fmt.Sprintf("|%v * %v", st, mul))
	}
	return sB.String()
}

type stone struct {
	num int
}

func createStone(s string) stone {
	st := stone{}
	atoi, _ := strconv.Atoi(s)
	st.num = atoi
	return st
}

func (s *stone) blink() (stone, bool) {
	if s.num == 0 {
		s.num = 1
		return stone{}, false
	}
	if len(strconv.Itoa(s.num))%2 == 0 {
		itoa := strconv.Itoa(s.num)
		s1Num, err := strconv.Atoi(itoa[:len(itoa)/2])
		if err != nil {
			panic(err)
		}
		s2Num, err := strconv.Atoi(itoa[len(itoa)/2:])
		if err != nil {
			panic(err)
		}
		s.num = s1Num
		return stone{num: s2Num}, true
	}
	s.num *= 2024
	return stone{}, false
}

func main() {
	start := time.Now()

	input := readFile("day11/input.txt")
	arrangement := createStoneArrangement(input[0])
	fmt.Println(arrangement.String())
	arrangement.blinkTimes(25)
	fmt.Println("Part 1:", arrangement.count())

	arrangement = createStoneArrangement(input[0])
	fmt.Println(arrangement.String())
	arrangement.blinkTimes(75)
	fmt.Println("Part 2:", arrangement.count())

	fmt.Println("Finished in", time.Since(start))
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
