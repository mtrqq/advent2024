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

type rulesMap map[int64][]int64

func readInstructions(filePath string) (rulesMap, [][]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open file: %v", err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read file: %v", err)
	}

	lastRuleIndex := -1
	instructionsSplit := bytes.Split(buffer, []byte(aoc.EndOfLine))

	rules := rulesMap{}
	for index, ruleStr := range instructionsSplit {
		// expect to receive empty string for separator indication
		if len(ruleStr) == 0 {
			lastRuleIndex = index
			break
		}

		left, right, found := bytes.Cut(ruleStr, []byte(aoc.Pipe))
		if !found {
			return nil, nil, fmt.Errorf("unable to parse rule at line#%d: not found separator", index)
		}

		dependency, err := strconv.ParseInt(string(left), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to parse number %q at line#%d: %v", left, index, err)
		}

		dependent, err := strconv.ParseInt(string(right), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to parse number %q at line#%d: %v", left, index, err)
		}

		rules[dependent] = append(rules[dependent], dependency)
	}

	if len(instructionsSplit) <= lastRuleIndex+1 {
		return nil, nil, fmt.Errorf("unable to parse instructions: no print sequences")
	}

	printSequences := make([][]int64, 0) // too lazy to compute expected sizing
	for index, sequenceStr := range instructionsSplit[lastRuleIndex+1:] {
		pageNumbersStr := bytes.Split(sequenceStr, []byte(aoc.Comma))
		sequence := make([]int64, len(pageNumbersStr))
		for numberIndex, pageNumberStr := range pageNumbersStr {
			number, err := strconv.ParseInt(string(pageNumberStr), aoc.Decimal, aoc.BitSize64)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to parse number %q at line#%d: %v", pageNumberStr, index+lastRuleIndex+1, err)
			}

			sequence[numberIndex] = number
		}

		printSequences = append(printSequences, sequence)
	}

	return rules, printSequences, nil
}

func expandRules(rules rulesMap) rulesMap {

}

func isValidSequence(rules rulesMap, printSeq []int64) bool {

}

func calculateMagicNumberPart1(rules rulesMap, printSequences [][]int64) int64 {
	result := int64(0)
	for _, seq := range printSequences {
		if isValidSequence(rules, seq) {
			result += seq[len(seq)/2]
		}
	}

	return result
}

func main() {
	flag.Parse()

	rules, printSequences, err := readInstructions(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read instructions file: %v", err)
	}

	resultPart1 := calculateMagicNumberPart1(rules, printSequences)
	log.Printf("[PrintSeq#Part1] Result: %d", resultPart1)
}
