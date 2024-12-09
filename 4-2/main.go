package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const (
	inputFile  = "input.txt"
	searchTerm = "XMAS"
)

type matrix struct {
	rows        []row
	rowCount    int
	columnCount int
	cur         cursor
}

type cursor struct {
	x int
	y int
}

type row []rune

func init() {
	setupLogger()
}

func main() {
	ch := make(chan []rune)
	go readInput(ch)
	m := new(matrix)
	m.rows = make([]row, 0)
	for r := range ch {
		m.addRow(r)
	}
	fmt.Println(m)
	//fmt.Println(searchTerm, "found", m.wordsearch(searchTerm), "times")
	fmt.Println("X-MAS found", m.xSearch([3]rune{'M', 'A', 'S'}), "times")
}

func setupLogger() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(os.Stdout, logOpts)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func readInput(ch chan []rune) {
	inputFile, _ := os.Open(inputFile)
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		ch <- []rune(line)
	}
	close(ch)
}

func (m *matrix) addRow(r row) {
	m.rows = append(m.rows, r)
	m.rowCount = len(m.rows)
	if m.columnCount == 0 {
		m.columnCount = len(r)
		return
	}
	if m.columnCount != len(r) {
		panic("inconsistent number of columns")
	}
}

func (m *matrix) xSearch(rs [3]rune) int {
	res := 0
	for rowNum := 1; rowNum < (m.rowCount - 1); rowNum++ {
		for columnNum := 1; columnNum < (m.columnCount - 1); columnNum++ {
			if m.rows[rowNum][columnNum] != rs[1] {
				continue
			}
			slog.Debug("Found center rune", slog.Int("x", columnNum), slog.Int("y", rowNum))
			m.cur.x = columnNum
			m.cur.y = rowNum
			if !m.isX(rs[0], rs[2]) {
				continue
			}
			res++
		}
	}
	return res
}

func (m matrix) String() string {
	sb := strings.Builder{}
	lm := len(m.rows)
	for i, row := range m.rows {
		sb.WriteString(row.String())
		if i == lm-1 {
			continue
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (rw row) String() string {
	sb := strings.Builder{}
	for _, r := range rw {
		sb.WriteRune(r)
	}
	return sb.String()
}

func (c cursor) String() string {
	return fmt.Sprintf("x: %d, y: %d", c.x, c.y)
}

// cursorVal retrieves the value at the cursor,
// optionally after moving the supplied vectors.
// Note this does not update the cursor.
func (m *matrix) cursorVal(vs ...vector) (rune, cursor) {
	tmpCursor := m.cur
	for _, v := range vs {
		tmpCursor.x += v.horiz
		tmpCursor.y += v.vert
	}
	return m.rows[tmpCursor.y][tmpCursor.x], tmpCursor
}

// isX fetches the NE, NW, SE & SW values from the cursor.
// It makes sure they are oneof the supplied values and that it makes a coherent X.
func (m *matrix) isX(left, right rune) bool {
	count := 0
	mx := [2][2]rune{}
	for vec, coord := range map[vector][2]int{
		vectorNE: {1, 0},
		vectorNW: {0, 0},
		vectorSE: {1, 1},
		vectorSW: {1, 0},
	} {
		v, _ := m.cursorVal(vec)
		if v != left && v != right {
			return false
		}
		slog.Debug(vec.String(), slog.String("val", string(v)))
		if v == right {
			count++
		}
		mx[coord[0]][coord[1]] = v
	}
	// check we have the right combination of letters and the diagonals are not equal
	return count == 2 && mx[1][0] != mx[0][1] && mx[0][0] != mx[1][1]
}
