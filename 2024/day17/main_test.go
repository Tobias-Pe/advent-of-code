package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_computer_executeInstructions(t *testing.T) {
	tests := []struct {
		name    string
		initC   computer
		resultC computer
	}{
		{name: "1", initC: computer{
			regA:             0,
			regB:             0,
			regC:             9,
			instructions:     []int{2, 6},
			outputs:          []int{},
			instructionIndex: 0,
		}, resultC: computer{
			regA:             0,
			regB:             1,
			regC:             9,
			instructions:     []int{2, 6},
			outputs:          []int{},
			instructionIndex: 2,
		}},
		{name: "2", initC: computer{
			regA:             10,
			regB:             0,
			regC:             0,
			instructions:     []int{5, 0, 5, 1, 5, 4},
			outputs:          []int{},
			instructionIndex: 0,
		}, resultC: computer{
			regA:             10,
			regB:             0,
			regC:             0,
			instructions:     []int{5, 0, 5, 1, 5, 4},
			outputs:          []int{0, 1, 2},
			instructionIndex: 6,
		}},
		{name: "3", initC: computer{
			regA:             2024,
			regB:             0,
			regC:             0,
			instructions:     []int{0, 1, 5, 4, 3, 0},
			outputs:          []int{},
			instructionIndex: 0,
		}, resultC: computer{
			regA:             0,
			regB:             0,
			regC:             0,
			instructions:     []int{0, 1, 5, 4, 3, 0},
			outputs:          []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0},
			instructionIndex: 6,
		}},
		{name: "4", initC: computer{
			regA:             0,
			regB:             29,
			regC:             0,
			instructions:     []int{1, 7},
			outputs:          []int{},
			instructionIndex: 0,
		}, resultC: computer{
			regA:             0,
			regB:             26,
			regC:             0,
			instructions:     []int{1, 7},
			outputs:          []int{},
			instructionIndex: 2,
		}},
		{name: "5", initC: computer{
			regA:             0,
			regB:             2024,
			regC:             43690,
			instructions:     []int{4, 0},
			outputs:          []int{},
			instructionIndex: 0,
		}, resultC: computer{
			regA:             0,
			regB:             44354,
			regC:             43690,
			instructions:     []int{4, 0},
			outputs:          []int{},
			instructionIndex: 2,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initC.executeInstructions()
			if !reflect.DeepEqual(tt.initC, tt.resultC) {
				t.Errorf("computer.executeInstructions() = %v, want %v", tt.initC, tt.resultC)
			}
		})
	}
}

func Test_readFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		want        string
	}{
		{name: "exmaple", fileContent: "Register A: 729\nRegister B: 0\nRegister C: 0\n\nProgram: 0,1,5,4,3,0", want: "4635635210"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.WriteFile("test.txt", []byte(tt.fileContent), 0644)
			computer := readFile("test.txt")
			os.Remove("test.txt")
			computer.executeInstructions()
			if computer.concatOutputs() != tt.want {
				t.Errorf("concatOutputs() = %v, want %v", computer.concatOutputs(), tt.want)
			}

		})
	}
}
