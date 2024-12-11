package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var dP map[stone]stoneArrangement

type stoneArrangement []stone

func createStoneArrangement(input string) stoneArrangement {
	stoneArrangement := make(stoneArrangement, len(strings.Fields(input)))
	for i, s := range strings.Fields(input) {
		stoneArrangement[i] = createStone(s)
	}
	return stoneArrangement
}

func (stA *stoneArrangement) blink() {
	for i := 0; i < len(*stA); i++ {
		newStone, isPresent := (*stA)[i].blink()
		if isPresent {
			tmpStALeft := append([]stone{}, (*stA)[:i+1]...)
			tmpStALeft = append(tmpStALeft, newStone)
			*stA = append(tmpStALeft, (*stA)[i+1:]...)
			i++
		}
	}
}

func (stA *stoneArrangement) blinkTimes(times int) {
	dP = make(map[stone]stoneArrangement)
	for i := 0; i < times; i++ {
		stA.blink()
		fmt.Println(i, stA)
	}
}

func (stA stoneArrangement) getBlinkAt(age int) stoneArrangement {
	stAAtAge := stoneArrangement{}
	for _, s := range stA {
		stAAtAge = append(stAAtAge, s.getBlinkAt(age)...)
	}
	return stAAtAge
}

func (s stone) getBlinkAt(age int) stoneArrangement {

	if age <= 0 {
		return []stone{s}
	}
	lookedForStone := stone{
		num: s.num,
		age: age,
	}
	if dP[lookedForStone] != nil {
		return dP[lookedForStone]
	}
	stA := append(stoneArrangement{}, s.getBlinkAt(age-1)...)
	stA.blink()
	dP[lookedForStone] = stA
	return dP[lookedForStone]
}

func (stA stoneArrangement) String() string {
	sB := strings.Builder{}
	for i := 0; i < len(stA); i++ {
		sB.WriteString(fmt.Sprintf("%v ", (stA)[i].num))
	}
	return sB.String()
}

type stone struct {
	num int
	age int
}

func createStone(s string) stone {
	st := stone{age: 0}
	atoi, _ := strconv.Atoi(s)
	st.num = atoi
	return st
}

func (s *stone) blink() (stone, bool) {
	s.age = s.age + 1
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
		return stone{num: s2Num, age: 0}, true
	}
	s.num *= 2024
	return stone{}, false
}

func main() {
	start := time.Now()

	input := readFile("day11/input.txt")
	dP = make(map[stone]stoneArrangement)
	arrangement := createStoneArrangement(input[0])
	fmt.Println(arrangement.String())
	arrangement.blinkTimes(6)
	fmt.Println("Part 1:", len(arrangement))
	fmt.Println("Finished in", time.Since(start))
	start = time.Now()
	arrangement = createStoneArrangement(input[0])
	fmt.Println(arrangement.String())
	arrangement.getBlinkAt(6)
	fmt.Println("Part 1 fast:", len(arrangement))

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
