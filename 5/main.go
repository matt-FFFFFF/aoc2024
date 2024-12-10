package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"
)

type rule [2]uint8
type book []uint8
type page2pos map[uint8]uint8

func main() {
	middles := 0
	rulesCh := make(chan string)
	bookCh := make(chan string)
	go readInput(rulesCh, bookCh)

	rulesOut := make(chan []rule, 1)
	bookOut := make(chan []book, 1)
	go processRules(rulesCh, rulesOut)
	go processBooks(bookCh, bookOut)
	rules, books := gatherData(rulesOut, bookOut)
bookLoop:
	for _, book := range books {
		slog.Info("book", slog.String("values", book.String()))
		p2p := bookPage2PosMap(book)
	ruleLoop:
		for _, rule := range rules {
			leftPos, lOk := p2p[rule[0]]
			rightPos, rOk := p2p[rule[1]]
			// if the pages cannot be looked up then this rule doesn't apply
			if !lOk || !rOk {
				continue ruleLoop
			}
			if leftPos > rightPos {
				continue bookLoop
			}
		}
		m := int(book[(len(book)-1)/2])
		middles += m
		slog.Info("Adding middle", slog.Int("middle", m), slog.Int("newTotal", middles))
	}
	fmt.Println(middles)
}

func bookPage2PosMap(b book) page2pos {
	res := make(page2pos)
	for i, v := range b {
		res[v] = uint8(i)
	}
	return res
}

func gatherData(rulesOut <-chan []rule, booksOut <-chan []book) ([]rule, []book) {
	var rules []rule
	var books []book
	for {
		select {
		case r, ok := <-rulesOut:
			if !ok {
				continue
			}
			rules = r
		case b, ok := <-booksOut:
			if !ok {
				continue
			}
			books = b
		}
		if rules != nil && books != nil {
			break
		}
	}
	return rules, books
}

// processRules returns a map of the
func processRules(in <-chan string, out chan<- []rule) {
	defer close(out)
	res := make([]rule, 0, len(in))
	for s := range in {
		left, err := strconv.Atoi(string(s[0:2]))
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(s[3:5])
		if err != nil {
			panic(err)
		}
		res = append(res, rule{uint8(left), uint8(right)})
	}
	out <- res
}

func processBooks(in <-chan string, out chan<- []book) {
	defer close(out)
	res := make([]book, 0, len(in))
	for line := range in {
		bk := make(book, 0, 10)
		ds := strings.Split(line, ",")
		for _, is := range ds {
			i, err := strconv.Atoi(is)
			if err != nil {
				panic(err)
			}
			bk = append(bk, uint8(i))
		}
		res = append(res, bk)
	}
	out <- res
}

func readInput(rulesCh, bookCh chan<- string) {
	inputFile, _ := os.Open(inputFile)
	scanner := bufio.NewScanner(inputFile)
	ch := rulesCh
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			close(ch)
			ch = bookCh
			continue
		}
		ch <- line
	}
	close(ch)
}

func (b book) String() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for _, v := range b {
		sb.WriteString(strconv.Itoa(int(v)))
		sb.WriteString(", ")
	}
	sb.WriteRune(']')
	return sb.String()
}
