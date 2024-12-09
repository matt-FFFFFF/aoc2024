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
	fmt.Println(searchTerm, "found", m.wordsearch(searchTerm), "times")
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

func (m *matrix) wordsearch(s string) int {
	res := 0
	sr := []rune(s)
	for rowNum, row := range m.rows {
		for columnNum, val := range row {
			if val != sr[0] {
				continue
			}
			slog.Debug("Start rune located", slog.Int("x", columnNum), slog.Int("y", rowNum))
			// find available vectors for given length
			vectors := m.availableVectors(len(s), columnNum, rowNum)
			// for each valid vector, scan from current location for 2nd rune onwards
		vectorLoop:
			for _, vec := range vectors {
				slog.Debug("Starting vector", slog.String("vector", vec.String()))
				m.cur = cursor{
					x: columnNum,
					y: rowNum,
				}
				slog.Debug("Cursor pos", slog.Int("x", m.cur.x), slog.Int("y", m.cur.y))
				for _, r := range sr[1:] {
					v := m.moveVector(vec)
					slog.Debug("Cursor val", slog.String("val", string(m.cursorVal())))
					if v != r {
						slog.Debug("No match", slog.String("want", string(r)), slog.String("got", string(m.cursorVal())))
						continue vectorLoop
					}
					slog.Debug("Match", slog.String("want", string(r)), slog.String("got", string(m.cursorVal())))
				}
				// increment count if match
				res++
				slog.Debug("Found word", slog.Int("x", columnNum), slog.Int("y", rowNum), slog.String("vector", vec.String()))
			}
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

func (m matrix) availableVectors(length int, x, y int) []vector {
	nOk := y >= length-1
	eOk := x <= m.rowCount-length
	sOk := y <= m.columnCount-length
	wOk := x >= length-1

	res := make([]vector, 0, 8)
	if nOk {
		slog.Debug("Available vector N")
		res = append(res, vectorN)
	}
	if nOk && eOk {
		slog.Debug("Available vector NE")
		res = append(res, vectorNE)
	}
	if eOk {
		slog.Debug("Available vector E")
		res = append(res, vectorE)
	}
	if sOk && eOk {
		slog.Debug("Available vector SE")
		res = append(res, vectorSE)
	}
	if sOk {
		slog.Debug("Available vector S")
		res = append(res, vectorS)
	}
	if sOk && wOk {
		slog.Debug("Available vector SW")
		res = append(res, vectorSW)
	}
	if wOk {
		slog.Debug("Available vector W")
		res = append(res, vectorW)
	}
	if nOk && wOk {
		slog.Debug("Available vector NW")
		res = append(res, vectorNW)
	}
	return res
}

func (m *matrix) moveVector(v vector) rune {
	m.cur.x += v.horiz
	m.cur.y += v.vert
	slog.Debug("Cursor pos", slog.Int("x", m.cur.x), slog.Int("y", m.cur.y))
	return m.cursorVal()
}

func (m *matrix) cursorVal() rune {
	return m.rows[m.cur.y][m.cur.x]
}
