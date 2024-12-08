package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

const (
	numRows = 1000
)

func main() {
	leftColumn, rightColumn := readInput()
	slices.Sort(leftColumn)
	slices.Sort(rightColumn)
	totalDistance := 0
	similarityMap := make(map[int]int)
	for i := range numRows {
		// part 2 make a map of the right column values count
		if _, ok := similarityMap[rightColumn[i]]; !ok {
			similarityMap[rightColumn[i]] = 0
		}
		similarityMap[rightColumn[i]]++

		distance := leftColumn[i] - rightColumn[i]

		if distance < 0 {
			distance = -distance
		}
		totalDistance += distance
		fmt.Println(i, leftColumn[i], rightColumn[i], distance)
	}
	fmt.Println("Total distance:", totalDistance)

	// Part 2 bits
	var similarity int
	for _, v := range leftColumn {
		if _, ok := similarityMap[v]; !ok {
			continue
		}
		similarity += v * similarityMap[v]
	}
	fmt.Println("Similarity:", similarity)
}

func readInput() ([]int, []int) {
	const lineLength = 14

	leftColumn := make([]int, numRows)
	rightColumn := make([]int, numRows)

	inputFile, _ := os.Open("input.txt")
	i := 0
	line := make([]byte, lineLength)
	for {
		_, err := inputFile.Read(line)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		leftColumn[i] = parseBytesAsInt(line[0:5])
		rightColumn[i] = parseBytesAsInt(line[8:13])
		i++
	}
	return leftColumn, rightColumn
}

func parseBytesAsInt(b []byte) int {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		panic(err)
	}
	return i
}
