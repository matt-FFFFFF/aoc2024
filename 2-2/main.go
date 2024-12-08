package main

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt" // the filename containing the input data
	tolerance = 3           // the acceptable diff between report numbers
)

func main() {
	setupLogger()
	ch := make(chan []int8)
	go readInput(ch)
	safe := 0
	for row := range ch {
		s := isSafe(row, -1)
		if s {
			slog.Info("safe", "row", row)
			safe++
		}
	}
	slog.Info("Complete", "safeCount", safe)
}

func setupLogger() {
	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logHandler := slog.NewTextHandler(os.Stdout, logOpts)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}

func readInput(ch chan []int8) {
	inputFile, _ := os.Open(inputFile)
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

func isSafe(report []int8, skip int) bool {
	slog.Info("isSafe",
		slog.Any("input", report),
		slog.Int("skip", skip),
	)
	if len(report) < 2 {
		slog.Info("safe as report has < 2 entries")
		return true
	}

	if skip >= 0 {
		newReport := make([]int8, 0, len(report)-1)
		for i, v := range report {
			if skip == i {
				continue
			}
			newReport = append(newReport, v)
		}
		report = newReport
		slog.Info("newReport",
			slog.Any("input", report),
		)
	}

	diffs := make([]int8, len(report)-1)
	for i := 1; i < len(report); i++ {
		diffs[i-1] = report[i] - report[i-1]
	}
	slog.Info("diffs", "val", diffs)
	var inc, dec, gt3, zero uint8
	for i, v := range diffs {
		switch {
		case v == 0:
			zero++
		case v > 3:
			gt3++
			inc++
		case v > 0 && v <= 3:
			inc++
		case v < -3:
			gt3++
			dec++
		case v < 0 && v >= -3:
			dec++
		}
		if inc > 0 && dec > 0 {
			if skip < 0 {
				for j := i + 1; j >= 0; j-- {
					if isSafe(report, j) {
						return true
					}
				}
			}
			return false
		}
		if gt3 > 0 || zero > 0 {
			if skip < 0 {
				for j := i + 1; j >= 0; j-- {
					if isSafe(report, j) {
						return true
					}
				}
			}
			return false
		}
	}
	return true
}
