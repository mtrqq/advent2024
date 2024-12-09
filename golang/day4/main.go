package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mtrqq/advent2024/golang"
)

var inputFile = flag.String("input", "input.txt", "Path to input file")

var (
	targetWord          = []byte("XMAS")
	validDiagonalShapes = map[string]struct{}{
		"SSAMM": {},
		"SMASM": {},
		"MMASS": {},
		"MSAMS": {},
	}
)

type byteGrid struct {
	array   [][]byte
	rows    int
	columns int
}

func newGrid(array [][]byte, rows int, columns int) byteGrid {
	return byteGrid{
		array:   array,
		rows:    rows,
		columns: columns,
	}
}

func (g byteGrid) get(row int, column int) byte {
	return g.array[row][column]
}

func readByteGrid(filePath string) (byteGrid, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return byteGrid{}, fmt.Errorf("unable to open file: %v", err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		return byteGrid{}, fmt.Errorf("unable to read file: %v", err)
	}

	array := bytes.Split(buffer, []byte(golang.EndOfLine))

	if len(array) == 0 {

		return newGrid(array, 0, 0), nil
	}

	rows := len(array)
	columns := len(array[0])

	for rowIndex, row := range array {
		if len(row) != columns {
			return byteGrid{}, fmt.Errorf("unmatching column count at row#%d", rowIndex)
		}
	}

	return newGrid(array, rows, columns), nil
}

func matchesLineInDirection(grid byteGrid, fromRow, fromColumn, directionRow, directionColumn int, targetLine []byte) bool {
	line := make([]byte, len(targetLine))
	row, column := fromRow, fromColumn
	for i := 0; i < len(targetLine); i++ {
		if row >= grid.rows || column >= grid.columns || row < 0 || column < 0 {
			return false
		}

		line[i] = grid.get(row, column)
		row += directionRow
		column += directionColumn
	}

	return bytes.Equal(line, targetLine)
}

func countWordsPart1(grid byteGrid) int {
	count := 0
	for row := 0; row < grid.rows; row++ {
		for column := 0; column < grid.columns; column++ {
			if targetWord[0] != grid.get(row, column) {
				continue
			}

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if matchesLineInDirection(grid, row, column, dx, dy, []byte(targetWord)) {
						count++
					}
				}
			}
		}
	}

	return count
}

func contains[K comparable, V any](mapping map[K]V, key K) bool {
	if _, exists := mapping[key]; exists {
		return true
	}

	return false
}

func collectXShape(grid byteGrid, row, column int) []byte {
	if row-1 < 0 || column-1 < 0 {
		return nil
	}

	if row+1 >= grid.rows || column+1 >= grid.columns {
		return nil
	}

	return []byte{
		grid.get(row-1, column-1),
		grid.get(row-1, column+1),
		grid.get(row, column),
		grid.get(row+1, column-1),
		grid.get(row+1, column+1),
	}
}

func countWordsPart2(grid byteGrid) int {
	count := 0

	for row := 0; row < grid.rows; row++ {
		for column := 0; column < grid.columns; column++ {
			if grid.get(row, column) != 'A' {
				continue
			}

			shape := string(collectXShape(grid, row, column))
			if contains(validDiagonalShapes, shape) {
				count++
			}
		}
	}

	return count
}

func main() {
	flag.Parse()

	grid, err := readByteGrid(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read instructions file: %v", err)
	}

	countPart1 := countWordsPart1(grid)
	log.Printf("[XmasWord#Part1] Count: %d", countPart1)

	countPart2 := countWordsPart2(grid)
	log.Printf("[XmasWord#Part2] Count: %d", countPart2)
}
