package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
)

var worryLevelDivisor = 1

type monkey struct {
	itemsLevel                 []*big.Int
	operationTask              string
	operationAmount            int
	testTask                   string
	testAmount                 int
	testTrueTargetMonkeyIndex  int
	testFalseTargetMonkeyIndex int
	counterInspections         int
}

func (monkey *monkey) playRound() {
	for _, worryLevel := range monkey.itemsLevel {
		monkey.inspect(worryLevel)
		monkey.getBored2(worryLevel)
		isTestSuccess := monkey.test(worryLevel)
		monkey.passItem(isTestSuccess, worryLevel)
	}
	monkey.itemsLevel = []*big.Int{}
}

func (monkey *monkey) inspect(worryLevel *big.Int) {
	monkey.counterInspections++
	var operationAmount *big.Int
	if monkey.operationAmount == -1 {
		operationAmount = worryLevel
	} else {
		operationAmount = big.NewInt(int64(monkey.operationAmount))
	}
	switch monkey.operationTask {
	case "+":
		worryLevel.Add(worryLevel, operationAmount)
	case "-":
		worryLevel.Sub(worryLevel, operationAmount)
	case "/":
		worryLevel.Div(worryLevel, operationAmount)
	case "*":
		worryLevel.Mul(worryLevel, operationAmount)
	default:
		fmt.Println(monkey.operationTask, " unknown operation task!")
	}
}

func (monkey *monkey) getBored(worryLevel *big.Int) {
	worryLevel.Div(worryLevel, big.NewInt(3))
}
func (monkey *monkey) getBored2(worryLevel *big.Int) {
	worryLevel.Mod(worryLevel, big.NewInt(int64(worryLevelDivisor)))
}

func (monkey *monkey) test(worryLevel *big.Int) bool {
	switch monkey.testTask {
	case "divisible":
		var mod = big.NewInt(0)
		mod = mod.Mod(worryLevel, big.NewInt(int64(monkey.testAmount)))
		return mod.Uint64() == 0
	default:
		fmt.Println(monkey.testTask, " unknown test task!")

	}
	return false
}

func (monkey *monkey) passItem(isTestSuccess bool, worryLevel *big.Int) {
	if isTestSuccess {
		monkeys[monkey.testTrueTargetMonkeyIndex].itemsLevel = append(monkeys[monkey.testTrueTargetMonkeyIndex].itemsLevel, worryLevel)
	} else {
		monkeys[monkey.testFalseTargetMonkeyIndex].itemsLevel = append(monkeys[monkey.testFalseTargetMonkeyIndex].itemsLevel, worryLevel)
	}
}

var monkeys []*monkey

func readFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	splittedLines := strings.Split(strings.ReplaceAll(lines, "\r\n", "\n"), "\n\n")
	for _, monkeyLines := range splittedLines {
		newMonkey := monkey{
			itemsLevel:                 []*big.Int{},
			operationTask:              "",
			operationAmount:            0,
			testTask:                   "",
			testAmount:                 0,
			testTrueTargetMonkeyIndex:  0,
			testFalseTargetMonkeyIndex: 0,
		}
		for i, line := range strings.Split(monkeyLines, "\n") {
			switch i {
			case 1:
				itemsString := strings.TrimPrefix(line, "  Starting items: ")
				for _, item := range strings.Split(itemsString, ", ") {
					worryLevel, _ := strconv.Atoi(item)
					newMonkey.itemsLevel = append(newMonkey.itemsLevel, big.NewInt(int64(worryLevel)))
				}
			case 2:
				operationString := strings.TrimPrefix(line, "  Operation: new = old ")
				operationStringArgs := strings.Split(operationString, " ")
				newMonkey.operationTask = operationStringArgs[0]
				if operationStringArgs[1] == "old" {
					newMonkey.operationAmount = -1
				} else {
					operationAmount, _ := strconv.Atoi(operationStringArgs[1])
					newMonkey.operationAmount = operationAmount
				}
			case 3:
				testString := strings.TrimPrefix(line, "  Test: ")
				testArgs := strings.Split(testString, " by ")
				newMonkey.testTask = testArgs[0]
				testAmount, _ := strconv.Atoi(testArgs[1])
				worryLevelDivisor *= testAmount
				newMonkey.testAmount = testAmount
			case 4:
				valueString := strings.TrimPrefix(line, "    If true: throw to monkey ")
				value, _ := strconv.Atoi(valueString)
				newMonkey.testTrueTargetMonkeyIndex = value
			case 5:
				valueString := strings.TrimPrefix(line, "    If false: throw to monkey ")
				value, _ := strconv.Atoi(valueString)
				newMonkey.testFalseTargetMonkeyIndex = value
			}
		}
		monkeys = append(monkeys, &newMonkey)
	}
}

func playRound() {
	for _, monkey := range monkeys {
		monkey.playRound()
	}
}

func playRounds(count int) {
	for i := 0; i < count; i++ {
		playRound()
		fmt.Println(i)
	}
}

func printMonkeyBusiness() {
	for i, monkey := range monkeys {
		fmt.Println("Monkey ", i, " counter: ", monkey.counterInspections)
	}
}

func main() {
	readFile("./input.txt")
	playRounds(10000)
	printMonkeyBusiness()

}
