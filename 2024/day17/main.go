package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type computer struct {
	regA, regB, regC int
	instructions     []int
	outputs          []int
	instructionIndex int
}

func (c *computer) concatOutputs() string {
	sb := strings.Builder{}
	for _, o := range c.outputs {
		sb.WriteString(strconv.Itoa(o))
		sb.WriteRune(',')
	}
	return sb.String()[:len(sb.String())-1]
}

func (c *computer) executeInstructions() {
	for len(c.instructions) > c.instructionIndex {
		c.executeInstruction()
	}
}

func (c *computer) executeInstructionsWithAbort() bool {
	for len(c.instructions) > c.instructionIndex {
		c.executeInstruction()
		if len(c.outputs) > 0 && c.outputs[len(c.outputs)-1] != c.instructions[len(c.outputs)-1] {
			return false
		}
	}
	return reflect.DeepEqual(c.outputs, c.instructions)
}

func (c *computer) getComboOperandValue() int {
	switch c.instructions[c.instructionIndex+1] {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		return c.instructions[c.instructionIndex+1]
	case 4:
		return c.regA
	case 5:
		return c.regB
	case 6:
		return c.regC
	case 7:
		panic("Combo operand 7 is reserved and will not appear in valid programs.")
	}
	panic("Invalid combo operand")
}

func (c *computer) reduceTo3Bits(x int) int {
	return x & 7
}

func (c *computer) executeInstruction() {
	currInstruction := c.instructions[c.instructionIndex]
	switch currInstruction {
	case 0:
		c.adv()
	case 1:
		c.bxl()
	case 2:
		c.bst()
	case 3:
		c.jnz()
	case 4:
		c.bxc()
	case 5:
		c.out()
	case 6:
		c.bdv()
	case 7:
		c.cdv()
	default:
		panic("Invalid instruction")
	}
	c.instructionIndex += 2
}

func (c *computer) dv() int {
	numerator := float64(c.regA)
	denominator := math.Pow(2, float64(c.getComboOperandValue()))
	result := int(math.Trunc(numerator / denominator))
	return result
}

func (c *computer) adv() {
	result := c.dv()
	c.regA = result
}

func (c *computer) bdv() {
	result := c.dv()
	c.regB = result
}

func (c *computer) cdv() {
	result := c.dv()
	c.regC = result
}

func (c *computer) bxl() {
	literalOperand := c.instructions[c.instructionIndex+1]
	result := c.regB ^ literalOperand // TODO this could be a bug as literalOperand is supposed to be only 3 bits
	c.regB = result
}

func (c *computer) bst() {
	val := c.getComboOperandValue()
	result := c.reduceTo3Bits(val % 8)
	c.regB = result
}

func (c *computer) jnz() {
	if c.regA == 0 {
		return
	}
	literalOperand := c.instructions[c.instructionIndex+1]
	c.instructionIndex = literalOperand
	c.instructionIndex -= 2 // prevent general jump 2 further
}

func (c *computer) bxc() {
	c.regB = c.regB ^ c.regC
}

func (c *computer) out() {
	val := c.getComboOperandValue()
	result := val % 8
	c.outputs = append(c.outputs, result)
}

func runComputerForUser(file string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("A: ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == "abort" {
			return
		}
		atoi, err := strconv.Atoi(strings.TrimSpace(strings.ReplaceAll(text, ".", "")))
		if err != nil {
			fmt.Println(err)
			continue
		}
		c := readFile(file)
		c.regA = atoi
		c.executeInstructions()
		fmt.Println(c.concatOutputs(), len(c.outputs))
	}
}

func runComputerFor2s(file string) {
	for i := 0; i < 200; i++ {
		pow := math.Pow(2, float64(i))
		c := readFile(file)
		c.regA = int(pow)
		c.executeInstructions()
		if len(c.outputs) > len(c.instructions) {
			return
		}
		fmt.Printf("A: %.0f | ", pow)
		fmt.Println(c.concatOutputs(), len(c.outputs))
	}
}

func runComputerFor8s(file string) {
	for i := 0; i < 20; i++ {
		pow := math.Pow(8, float64(i))
		c := readFile(file)
		c.regA = int(pow)
		c.executeInstructions()
		if len(c.outputs) > len(c.instructions) {
			break
		}
		fmt.Printf("A: %.0f | ", pow)
		fmt.Println(c.concatOutputs(), len(c.outputs))
	}
	init := 0
	for i := 0; i < 200000; i++ {
		c := readFile(file)
		c.regA = init
		c.executeInstructions()
		if len(c.outputs) > len(c.instructions) {
			break
		}
		fmt.Printf("A: %d | ", init)
		fmt.Println(c.concatOutputs(), len(c.outputs))
		init += 8
	}
}

func runComputerFor1s(file string) {
	for i := 0; i < 128; i++ {
		c := readFile(file)
		c.regA = i
		c.executeInstructions()
		if len(c.outputs) > len(c.instructions) {
			return
		}
		fmt.Printf("A: %d | ", i)
		fmt.Println(c.concatOutputs(), len(c.outputs))
	}
}

func bruteForceP2(file string, solA, indx, targetIndx int) int {
	if indx == targetIndx {
		if solA < int(math.Pow(8, float64(targetIndx-1))) {
			return -1
		}
		comptr := readFile(file)
		comptr.regA = solA
		success := comptr.executeInstructionsWithAbort()
		if success {
			return solA
		}
		return -1
	}

	if indx == 0 {
		solA = 0
	} else {
		solA = solA << 3
	}

	for num := range 8 {
		a := solA | num
		comptr := readFile(file)
		comptr.regA = a
		comptr.executeInstructions()
		outLen := indx + 1
		if len(comptr.outputs) < outLen {
			continue
		}
		outRelevant := comptr.outputs[len(comptr.outputs)-outLen:]
		instructRelevant := comptr.instructions[len(comptr.instructions)-outLen:]
		if reflect.DeepEqual(outRelevant, instructRelevant) {
			sol := bruteForceP2(file, a, outLen, targetIndx)
			if sol != -1 {
				return sol
			}
		}
	}

	return -1
}

func main() {
	start := time.Now()

	file := "day17/input.txt"
	comptr := readFile(file)
	comptr.executeInstructions()
	fmt.Println("Part 1:", comptr.concatOutputs())
	//runComputerForUser(file)
	fmt.Println("Part 2: ", bruteForceP2(file, 0, 0, len(comptr.instructions)))

	fmt.Println("Finished in", time.Since(start))
}

func readFile(file string) computer {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	lines = strings.ReplaceAll(lines, "\r\n", "\n")
	lines = strings.TrimSpace(lines)
	split := strings.Split(lines, "\n")
	computer := computer{}
	regA, _ := strconv.Atoi(strings.TrimLeft(split[0], "Register A: "))
	regB, _ := strconv.Atoi(strings.TrimLeft(split[1], "Register B: "))
	regC, _ := strconv.Atoi(strings.TrimLeft(split[2], "Register C: "))
	computer.regA = regA
	computer.regB = regB
	computer.regC = regC
	computer.instructions = []int{}
	for _, instrStr := range strings.Split(strings.TrimLeft(split[4], "Program: "), ",") {
		atoi, _ := strconv.Atoi(instrStr)
		computer.instructions = append(computer.instructions, atoi)
	}
	computer.outputs = []int{}
	computer.instructionIndex = 0
	return computer
}
