package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type mulOperation struct {
	left, right int
}

type condOperation struct {
	do bool
}

func main() {
	start := time.Now()

	lines := readFile("day3/input.txt")

	fmt.Println("Part 1:", sumAllMulOperations(lines))
	fmt.Println("Part 2:", sumAllMulOperationsWithConditionalStatements(lines))

	fmt.Println("Finished in", time.Since(start))
}

func sumAllMulOperations(lines []string) int {
	singleLineBuilder := strings.Builder{}
	for _, line := range lines {
		singleLineBuilder.WriteString(line)
	}
	singleLine := singleLineBuilder.String()
	sum := 0
	for i := 0; i < len(singleLine); {
		mulOp, increment := findMulOperation(singleLine, i)
		i += increment
		if mulOp != nil {
			sum += mulOp.left * mulOp.right
		}
	}
	return sum
}

func sumAllMulOperationsWithConditionalStatements(lines []string) int {
	singleLineBuilder := strings.Builder{}
	for _, line := range lines {
		singleLineBuilder.WriteString(line)
	}
	singleLine := singleLineBuilder.String()
	sum := 0
	enabled := true
	for i := 0; i < len(singleLine); {
		condOp, increment := findCondOperation(singleLine, i)
		i += increment
		if condOp != nil {
			enabled = condOp.do
		}
		mulOp, increment := findMulOperation(singleLine, i)
		i += increment
		if mulOp != nil && enabled {
			sum += mulOp.left * mulOp.right
		}
	}
	return sum
}

func findCondOperation(input string, startOffset int) (foundOperation *condOperation, skipAmount int) {
	doString := "do()"
	doNotString := "don't()"

	if input[startOffset:startOffset+len(doString)] == doString {
		return &condOperation{true}, len(doString)
	}

	if input[startOffset:startOffset+len(doNotString)] == doNotString {
		return &condOperation{false}, len(doNotString)
	}

	return nil, 0
}

func findMulOperation(input string, startOffset int) (foundOperation *mulOperation, skipAmount int) {
	offsetAfterMul := startOffset + 4
	if offsetAfterMul > len(input) {
		return nil, offsetAfterMul - startOffset
	}
	nextFourChars := input[startOffset:offsetAfterMul]
	if nextFourChars != "mul(" {
		return nil, 1
	}

	firstNumberEndIndex := offsetAfterMul + strings.Index(input[offsetAfterMul:], ",")
	if firstNumberEndIndex < offsetAfterMul {
		return nil, offsetAfterMul - startOffset
	}
	firstNumberString := input[offsetAfterMul:firstNumberEndIndex]
	firstNumber, err := strconv.Atoi(firstNumberString)
	if err != nil {
		return nil, offsetAfterMul - startOffset
	}

	offsetAfterFirstNumberAndComma := offsetAfterMul + len(firstNumberString) + 1
	if offsetAfterFirstNumberAndComma > len(input) {
		return nil, offsetAfterMul - startOffset
	}
	secondNumberEndIndex := offsetAfterFirstNumberAndComma + strings.Index(input[offsetAfterFirstNumberAndComma:], ")")
	if secondNumberEndIndex < offsetAfterFirstNumberAndComma {
		return nil, offsetAfterMul - startOffset
	}
	secondNumberString := input[offsetAfterFirstNumberAndComma:secondNumberEndIndex]
	secondNumber, err := strconv.Atoi(secondNumberString)
	if err != nil {
		return nil, offsetAfterFirstNumberAndComma - startOffset
	}

	return &mulOperation{
		left:  firstNumber,
		right: secondNumber,
	}, 4 + len(firstNumberString) + 1 + len(secondNumberString) + 1
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
