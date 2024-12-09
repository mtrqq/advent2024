package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"

	aoc "github.com/mtrqq/advent2024/golang"
	"golang.org/x/exp/constraints"
)

var inputFile = flag.String("input", "input.txt", "Path to input file")

func abs[T constraints.Signed](value T) T {
	if value < 0 {
		return -value
	}

	return value
}

func calculateDistancePart1(left, right []int64) uint64 {
	slices.Sort(left)
	slices.Sort(right)

	distanceTotal := uint64(0)
	for index := range left {
		leftNumber := left[index]
		rightNumber := right[index]
		distance := abs(leftNumber - rightNumber)
		distanceTotal += uint64(distance)
	}

	return distanceTotal
}

func calculateDistancePart2(left, right []int64) uint64 {
	rightCounts := map[int64]int{}
	for _, rightNumber := range right {
		rightCounts[rightNumber]++
	}

	distanceTotal := uint64(0)
	for _, leftNumber := range left {
		distance := rightCounts[leftNumber] * int(leftNumber)
		distanceTotal += uint64(distance)
	}

	return distanceTotal
}

func parseListsFile(filePath string) ([]int64, []int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []int64{}, []int64{}, fmt.Errorf("unable to open file: %v", err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		return []int64{}, []int64{}, fmt.Errorf("unable to read file: %v", err)
	}

	lines := bytes.Split(buffer, []byte(aoc.EndOfLine))

	left := make([]int64, len(lines))
	right := make([]int64, len(lines))

	for index, line := range lines {
		// If it's an empty line - skip processing
		if bytes.Equal(line, []byte("")) {
			log.Printf("Skip processing of line#%d: empty line", index)
			continue
		}

		before, after, found := bytes.Cut(line, []byte(aoc.Space))
		if !found {
			return []int64{}, []int64{}, fmt.Errorf("unable to find 2 numbers at line#%d", index)
		}

		after = bytes.TrimLeft(after, aoc.Space)

		beforeNumber, err := strconv.ParseInt(string(before), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			return []int64{}, []int64{}, fmt.Errorf("unable to parse number at line#%d (%s): %v", index, before, err)
		}

		afterNumber, err := strconv.ParseInt(string(after), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			return []int64{}, []int64{}, fmt.Errorf("unable to parse number at line#%d (%s): %v", index, before, err)
		}

		left[index] = beforeNumber
		right[index] = afterNumber
	}

	return left, right, nil

}

func main() {
	flag.Parse()

	log.Printf("Processing file at: %s\n", *inputFile)
	left, right, err := parseListsFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed parsing file: %v", err)
	}

	answer := calculateDistancePart1(left, right)
	log.Printf("[Part1] Distance between lists: %v", answer)

	answer = calculateDistancePart2(left, right)
	log.Printf("[Part2] Distance between lists: %v", answer)
}
