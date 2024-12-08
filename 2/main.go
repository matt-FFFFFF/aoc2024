package main

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func main() {
	setupLogger()
	ch := make(chan []int8)
	go readInput(ch)
	safe := 0
	for row := range ch {
		s, i := isSafe(row)
		slog.Info("Try0", "row", row, "safe", s, "fault", i)
		if s {
			safe++
		}
	}
	slog.Info("Safe", "count", safe)
}

func readInput(ch chan []int8) {
	inputFile, _ := os.Open("input.txt")
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		ch <- parseLine(line)
	}
	close(ch)
}

func parseLine(line string) []int8 {
	res := make([]int8, 0, 10)
	lineValues := strings.Split(line, " ")
	for _, value := range lineValues {
		i, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		res = append(res, int8(i))
	}
	return res
}

func isSafe(report []int8) (bool, int) {
	var increasing bool // true if the first value is less than the second
	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]
		// if the diff is zero or diverged by more than +-3 ...
		if diff == 0 || diff > 3 || diff < -3 {
			return false, i
		}
		// set increasing trend based on the first pair.
		if i == 1 {
			increasing = diff > 0
			continue
		}
		// if we are no longer following the increasing trend
		if increasing != (diff > 0) {
			return false, i
		}
	}
	return true, 0
}

func setupLogger() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(os.Stdout, logOpts)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
