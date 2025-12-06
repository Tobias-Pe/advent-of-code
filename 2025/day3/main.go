package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Bank struct {
	batteries []int
}

func NewBank(line string) *Bank {
	strBatteries := strings.Split(line, "")
	batteries := make([]int, len(strBatteries))
	for i, s := range strBatteries {
		batteries[i], _ = strconv.Atoi(s)
	}
	return &Bank{batteries}
}

func (b Bank) LargestPossibleJoltageP1() int {
	//fmt.Print(b.batteries)
	max1, max2 := 0, 0
	for i := 0; i < len(b.batteries)-1; i++ {
		if b.batteries[i] > max1 {
			max1, max2 = b.batteries[i], b.batteries[i+1]
		} else if b.batteries[i] > max2 {
			max2 = b.batteries[i]
		}
	}
	if b.batteries[len(b.batteries)-1] > max2 {
		max2 = b.batteries[len(b.batteries)-1]
	}
	res, _ := strconv.Atoi(
		strconv.FormatInt(int64(max1), 10) + strconv.FormatInt(int64(max2), 10),
	)
	//fmt.Println("|", max1, max2)
	return res
}

func (b Bank) LargestPossibleJoltageP2() int64 {
	//fmt.Print(b.batteries)
	nums := make([]int, 12)
	lastTakenIdx := -1
	for i := range nums {
		maxJ := len(b.batteries) - (12 - i)
		maxJoltage := 0
		for j := lastTakenIdx + 1; j <= maxJ; j++ {
			if b.batteries[j] > maxJoltage {
				maxJoltage = b.batteries[j]
				lastTakenIdx = j
			}
		}
		nums[i] = maxJoltage
	}

	strNum := ""
	for i := range nums {
		strNum += strconv.Itoa(nums[i])
	}
	res, _ := strconv.ParseInt(strNum, 10, 64)
	//fmt.Println("|", res)
	return res
}

func main() {
	start := time.Now()

	lines := readFile("day3/input.txt")
	banks := make([]*Bank, len(lines))
	for i, line := range lines {
		banks[i] = NewBank(line)
	}

	totalJoltage := 0
	for _, bank := range banks {
		totalJoltage += bank.LargestPossibleJoltageP1()
	}
	fmt.Println("Part01: Total Joltage is", totalJoltage)

	totalJoltageP2 := int64(0)
	for _, bank := range banks {
		totalJoltageP2 += bank.LargestPossibleJoltageP2()
	}
	fmt.Println("Part02: Total Joltage is", totalJoltageP2)

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
