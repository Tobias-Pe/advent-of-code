package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type diskPart struct {
	id     int
	isFree bool
	size   int
}

func (dP diskPart) String() string {
	toBeRepeated := "."
	if !dP.isFree {
		toBeRepeated = strconv.Itoa(dP.id)
	}
	return strings.Repeat(toBeRepeated, dP.size)
}

func (dP diskPart) calcFilesysChecksum(pos int) int {
	if dP.isFree {
		return 0
	}

	sum := 0

	for i := 0; i < dP.size; i++ {
		sum += dP.id * (pos + i)
	}

	return sum
}

type diskMap struct {
	diskParts []diskPart
}

func (dM diskMap) String() string {
	sB := strings.Builder{}
	for _, dP := range dM.diskParts {
		sB.WriteString(dP.String())
	}
	return sB.String()
}

func (dM diskMap) defragmentSplit() []diskPart {
	var defragmented []diskPart

	j := len(dM.diskParts) - 1
	for i := 0; i <= j; i++ {
		defragmented = append(defragmented, dM.diskParts[i])
		for defragmented[len(defragmented)-1].isFree && i <= j {
			oldFreeSize := defragmented[len(defragmented)-1].size

			lastDP := &dM.diskParts[j]
			if lastDP.isFree {
				j--
				continue
			}
			oldLastDPSize := lastDP.size
			lastDP.size = max(lastDP.size-oldFreeSize, 0)
			if lastDP.size == 0 {
				lastDP.isFree = true
				j--
			}

			defragmented[len(defragmented)-1].isFree = false
			defragmented[len(defragmented)-1].id = lastDP.id
			defragmented[len(defragmented)-1].size = min(oldLastDPSize, oldFreeSize)
			if defragmented[len(defragmented)-1].size < oldFreeSize {
				defragmented = append(defragmented, diskPart{
					id:     0,
					isFree: true,
					size:   oldFreeSize - defragmented[len(defragmented)-1].size,
				})
			}
		}

		for defragmented[len(defragmented)-1].isFree {
			defragmented = defragmented[:len(defragmented)-1]
		}
	}

	return defragmented
}

func (dM diskMap) defragmentWhole() []diskPart {
	defragmented := append([]diskPart{}, dM.diskParts...)
	alreadyTouched := make(map[int]bool)
	for lastElement := len(defragmented) - 1; lastElement >= 0; lastElement-- {
		toBeMoved := defragmented[lastElement]
		if toBeMoved.isFree || alreadyTouched[toBeMoved.id] {
			continue
		}
		alreadyTouched[toBeMoved.id] = true
		for i := 0; i < lastElement; i++ {
			if !defragmented[i].isFree || defragmented[i].size < toBeMoved.size {
				continue
			}

			if defragmented[i].size == toBeMoved.size {
				defragmented[i] = toBeMoved
				defragmented[lastElement].isFree = true
				break
			}

			defragmentedLeft := append([]diskPart{}, defragmented[:i+1]...)
			defragmentedLeft = append(defragmentedLeft, diskPart{defragmented[i].id, defragmented[i].isFree, defragmented[i].size - toBeMoved.size})
			defragmented = append(defragmentedLeft, defragmented[i+1:]...)
			lastElement++
			defragmented[i] = toBeMoved
			defragmented[lastElement].isFree = true
			break
		}
	}
	for defragmented[len(defragmented)-1].isFree {
		defragmented = defragmented[:len(defragmented)-1]
	}

	return defragmented
}

func createDiskMap(line string) *diskMap {
	dM := diskMap{diskParts: make([]diskPart, len(line))}
	isFree := false
	fileCounter := 0
	for i := 0; i < len(line); i++ {
		size, _ := strconv.Atoi(string(line[i]))
		dM.diskParts[i] = diskPart{
			id:     fileCounter,
			isFree: isFree,
			size:   size,
		}
		if isFree {
			fileCounter++
		}
		isFree = !isFree
	}
	return &dM
}

func calcFilesysChecksum(diskParts []diskPart) int64 {
	sum := int64(0)

	pos := 0
	for _, part := range diskParts {
		sum += int64(part.calcFilesysChecksum(pos))
		pos += part.size
	}

	return sum
}

func main() {
	start := time.Now()

	input := readFile("day9/input.txt")
	if len(input) > 1 {
		panic("input expected one liner")
	}
	diskMap := createDiskMap(input[0])
	defragment := diskMap.defragmentSplit()
	fmt.Println("Part 1:", calcFilesysChecksum(defragment))
	defragmentWhole := createDiskMap(input[0]).defragmentWhole()
	fmt.Println("Part 2:", calcFilesysChecksum(defragmentWhole))

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
