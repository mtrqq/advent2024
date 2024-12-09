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

const (
	targetWord = "XMAS"
)

type byteRoutes map[byte][]gridPoint

type byteGrid struct {
	array   [][]byte
	rows    int
	columns int

	routeCache map[gridPoint]byteRoutes
}

func newGrid(array [][]byte, rows int, columns int) byteGrid {
	return byteGrid{
		array:      array,
		rows:       rows,
		columns:    columns,
		routeCache: make(map[gridPoint]byteRoutes),
	}
}

func (g byteGrid) get(row int, column int) byte {
	return g.array[row][column]
}

func (g byteGrid) routesFrom(row, column int) byteRoutes {
	from := gridPoint{row: row, column: column}
	if routes, exists := g.routeCache[from]; exists {
		return routes
	}

	routes := byteRoutes{}
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			adjRow, adjCol := from.row+i, from.column+j
			if adjRow >= g.rows || adjCol >= g.columns || adjRow < 0 || adjCol < 0 {
				continue
			}

			adjPoint := gridPoint{row: adjRow, column: adjCol}
			cellByte := g.array[adjRow][adjCol]
			points := append(routes[cellByte], adjPoint)
			routes[cellByte] = points
		}
	}

	g.routeCache[from] = routes
	return routes
}

func (g byteGrid) routesFromTo(row int, column int, b byte) []gridPoint {
	return g.routesFrom(row, column)[b]
}

type gridPoint struct {
	row    int
	column int
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

type countEntry struct {
	row    int
	column int
	offset int
}

func countWordsPart1(grid byteGrid) int {
	stack := []countEntry{}

	for i := 0; i < grid.rows; i++ {
		for j := 0; j < grid.columns; j++ {
			stack = append(stack, countEntry{
				row:    i,
				column: j,
				offset: 0,
			})
		}
	}

	count := 0
	for len(stack) > 0 {
		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if targetWord[entry.offset] != grid.get(entry.row, entry.column) {
			continue
		}

		if entry.offset >= len(targetWord)-1 {
			count++
			continue
		}

		nextLetter := targetWord[entry.offset+1]
		nextPoints := grid.routesFromTo(entry.row, entry.column, nextLetter)
		for _, point := range nextPoints {
			stack = append(stack, countEntry{
				row:    point.row,
				column: point.column,
				offset: entry.offset + 1,
			})
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

	count := countWordsPart1(grid)
	log.Printf("[XmasWord#Part1] Count: %d", count)
}
