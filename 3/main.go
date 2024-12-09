package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	multRe := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	input := readInput()
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	scanner.Split(splitDoDont)
	res := 0
	for scanner.Scan() {
		res += parseMult(scanner.Bytes(), multRe)
	}
	fmt.Println(res)
}

func readInput() []byte {
	inputFile, _ := os.ReadFile("input.txt")
	return inputFile
}

// parseMult will use the supplied regular expression to search for the
// values to multiply.
// it will sum all of the multiplications and return the value.
func parseMult(b []byte, re *regexp.Regexp) int {
	res := 0
	a := re.FindAllSubmatch(b, -1)
	for _, match := range a {
		res += multByteSlice(match[1], match[2])
	}
	return res
}

// multByteSlice will multiply the string values represented as byte slices.
// It assumes that the inputs are integers stored as strings.
// It will panic if the conversion to integer fails.
func multByteSlice(bss ...[]byte) int {
	res := 0
	for _, bs := range bss {
		i, err := strconv.Atoi(string(bs))
		if err != nil {
			panic(err)
		}
		if res == 0 {
			res = i
			continue
		}
		res = res * i
	}
	return res
}

// splitDoDont is a bufio.SplitFunc that will provide tokens based on the
// bytes between `do()` and `don't()` in the supplied data.
// It will assume that the beginning of the data to be read is in the do() state
// and will advance the scanner to the start of the next do state.
func splitDoDont(data []byte, atEOF bool) (int, []byte, error) {
	var token []byte

	dontIndex := indexString(data, "don't()")
	if dontIndex == -1 {
		if !atEOF {
			return 0, nil, nil
		}
		return len(data), data, bufio.ErrFinalToken
	}
	token = data[:dontIndex]

	doIndex := indexString(data[dontIndex:], "do()")
	if doIndex == -1 {
		if !atEOF {
			return 0, nil, nil
		}
		return len(data), token, bufio.ErrFinalToken
	}

	return doIndex + dontIndex, token, nil
}

// indexString find the index of the given string within the byte slice.
// It will return -1 if the string is not found.
func indexString(in []byte, s string) int {
	sb := []byte(s)
outer:
	for index := 0; index < len(in); index++ {
		for i := 0; i < len(sb); i++ {
			if in[index+i] != sb[i] {
				index += i
				continue outer
			}
		}
		return index
	}
	return -1
}
