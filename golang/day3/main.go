package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"

	aoc "github.com/mtrqq/advent2024/golang"
)

var inputFile = flag.String("input", "input.txt", "Path to input file")

var (
	mulPattern               = regexp.MustCompile(`mul\((\d+?),(\d+?)\)`)
	validDigits              = []byte("1234567890")
	enableOperationsCommand  = []byte("do()")
	disableOperationsCommand = []byte("don't()")
	mulCommandPrefix         = []byte("mul(")
)

func readInstructions(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read instructions file: %v", err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read instructions file: %v", err)
	}

	return buffer, nil
}

func executeInstructionsPart1(instructions []byte) int64 {
	matches := mulPattern.FindAllSubmatch(instructions, -1)
	if len(matches) == 0 {
		return 0
	}

	result := int64(0)
	for _, match := range matches {
		firstOperand, err := strconv.ParseInt(string(match[1]), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			log.Printf("Matched against malformed number %s: %v", match[1], err)
			continue
		}

		secondOperand, err := strconv.ParseInt(string(match[2]), aoc.Decimal, aoc.BitSize64)
		if err != nil {
			log.Printf("Matched against malformed number %s: %v", match[2], err)
			continue
		}

		result += firstOperand * secondOperand
	}

	return result
}

func isDigit(symbol byte) bool {
	return slices.Contains(validDigits, symbol)
}

func locateNumber(buffer []byte, offset int) (int, int, bool) {
	numberStart := offset
	if len(buffer) <= numberStart || !isDigit(buffer[numberStart]) {
		return 0, 0, false
	}

	numberEnd := numberStart + 1
	for len(buffer) > numberEnd && isDigit(buffer[numberEnd]) {
		numberEnd++
	}

	return numberStart, numberEnd, true
}

// parseMulCommand parses first mul command that it finds
// in a provided instructions buffer, in case no mul command
// found - returns
func parseMulCommand(instructions []byte) (int64, int64, bool) {
	if !bytes.HasPrefix(instructions, mulCommandPrefix) {
		return 0, 0, false
	}

	opStart, opEnd, found := locateNumber(instructions, len(mulCommandPrefix))
	if !found {
		return 0, 0, false
	}

	firstOperandStr := string(instructions[opStart:opEnd])
	firstOperand, err := strconv.ParseInt(firstOperandStr, aoc.Decimal, aoc.BitSize64)
	if err != nil {
		log.Printf("Error: unable to parse number %s: %v", firstOperandStr, err)
	}

	if len(instructions) < opEnd || instructions[opEnd] != ',' {
		return 0, 0, false
	}

	opStart, opEnd, found = locateNumber(instructions, opEnd+1)
	if !found {
		return 0, 0, false
	}

	secondOperandStr := string(instructions[opStart:opEnd])
	secondOperand, err := strconv.ParseInt(secondOperandStr, aoc.Decimal, aoc.BitSize64)
	if err != nil {
		log.Printf("Error: unable to parse number %s: %v", secondOperandStr, err)
	}

	if len(instructions) < opEnd || instructions[opEnd] != ')' {
		return 0, 0, false
	}

	return firstOperand, secondOperand, true
}

// too boring to do that with regexp, let's try something new :~)

func executeInstructionsPart2(instructions []byte) int64 {
	operationsDisabled := false
	result := int64(0)
	for index := range instructions {
		if bytes.HasPrefix(instructions[index:], enableOperationsCommand) {
			operationsDisabled = false
		}

		if bytes.HasPrefix(instructions[index:], disableOperationsCommand) {
			operationsDisabled = true
		}

		if !operationsDisabled {
			operand1, operand2, validCmd := parseMulCommand(instructions[index:])
			if validCmd {
				result += operand1 * operand2
			}
		}
	}

	return result
}

func main() {
	flag.Parse()

	instructions, err := readInstructions(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read instructions file: %v", err)
	}

	part1Output := executeInstructionsPart1(instructions)
	log.Printf("[MulRegexp#Part1] Answer: %d", part1Output)

	part2Output := executeInstructionsPart2(instructions)
	log.Printf("[Mul#Part2] Answer: %d", part2Output)
}
