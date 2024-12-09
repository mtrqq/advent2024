package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	aoc "github.com/mtrqq/advent2024/golang"
)

var inputFile = flag.String("input", "input.txt", "Path to input file")

const (
	minLevels = 2
)

func parseReportsFile(filePath string) ([][]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open reports file: %v", err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read reports file: %v", err)
	}

	lines := bytes.Split(buffer, []byte(aoc.EndOfLine))
	reports := make([][]int64, len(lines))
	for index, line := range lines {
		levelsStr := bytes.Split(line, []byte(aoc.Space))
		if len(levelsStr) < minLevels {
			return nil, fmt.Errorf("unable to parse report at line#%d: count of levels less than %d", index, minLevels)
		}

		reports[index] = make([]int64, len(levelsStr))

		for levelIndex, levelStr := range levelsStr {
			level, err := strconv.ParseInt(string(levelStr), aoc.Decimal, aoc.BitSize64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse level at line#%d seq#%d: %v", index, levelIndex, err)
			}

			reports[index][levelIndex] = level
		}
	}

	return reports, nil
}

func testSequence(seq []int64, minDelta int64, maxDelta int64) bool {
	currentItem := 1
	prevItem := 0

	for currentItem < len(seq) {
		delta := seq[currentItem] - seq[prevItem]
		if delta < minDelta || delta > maxDelta {
			return false
		}

		currentItem++
		prevItem++
	}

	return true
}

func isSafeReport(report []int64) bool {
	if len(report) < 2 {
		log.Printf("Warning: found report with less than 2 levels, considering unsafe")
		return false
	}

	if report[1] > report[0] {
		return testSequence(report, 1, 3)
	}

	if report[1] < report[0] {
		return testSequence(report, -3, -1)
	}

	return false
}

func countSafeReportsPart1(reports [][]int64) int {
	safeCount := 0

	for _, report := range reports {
		if isSafeReport(report) {
			safeCount++
		}
	}

	return safeCount
}

func isSafeReportDampener(report []int64) bool {
	if isSafeReport(report) {
		return true
	}

	for index := range report {
		skipped := aoc.ArrayWithoutItem(report, index)
		if isSafeReport(skipped) {
			return true
		}
	}

	return false
}

func countSafeReportsPart2(reports [][]int64) int {
	safeCount := 0

	for _, report := range reports {
		if isSafeReportDampener(report) {
			safeCount++
		}
	}

	return safeCount
}

func main() {
	flag.Parse()

	reports, err := parseReportsFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to parse reports file: %v", err)
	}

	count := countSafeReportsPart1(reports)
	log.Printf("[Part#1] Count of safe reports: %d", count)

	count = countSafeReportsPart2(reports)
	log.Printf("[Part#2] Count of safe reports: %d", count)
}
