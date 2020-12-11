package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type dir struct {
	drow int
	dcol int
}

var directions = []dir{dir{-1, -1}, dir{-1, 0}, dir{-1, 1}, dir{0, -1}, dir{0, 1}, dir{1, -1}, dir{1, 0}, dir{1, 1}}

func nextState1(layout [][]rune, nextLayout [][]rune, rows int, cols int, row int, col int) bool {
	changed := false
	if layout[row][col] != '.' {
		occupied := 0
		for _, d := range directions {
			r := row + d.drow
			c := col + d.dcol
			if r >= 0 && r < rows && c >= 0 && c < cols {
				if layout[r][c] == '#' {
					occupied++
				}
			}
		}
		if layout[row][col] == 'L' {
			if occupied == 0 {
				nextLayout[row][col] = '#'
				changed = true
			}
		} else if occupied > 3 {
			nextLayout[row][col] = 'L'
			changed = true
		}
	}
	return changed
}

func nextLayout(layout [][]rune, rows int, cols int, fn func([][]rune, [][]rune, int, int, int, int) bool) ([][]rune, bool) {
	next := make([][]rune, rows)
	for i := range layout {
		next[i] = make([]rune, cols)
		copy(next[i], layout[i])
	}
	changed := false
	for r := range layout {
		for c := range layout[r] {
			elemChanged := fn(layout, next, rows, cols, r, c)
			changed = changed || elemChanged
		}
	}
	return next, changed
}

func solve1(layout [][]rune, rows int, cols int) [][]rune {
	next, changed := nextLayout(layout, rows, cols, nextState1)
	for changed {
		next, changed = nextLayout(next, rows, cols, nextState1)
	}
	return next
}

func seeOccupied(layout [][]rune, rows int, cols int, row int, col int, d dir) int {
	occupied := 0
	found := false
	r := row + d.drow
	c := col + d.dcol
	for !found {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			break
		}
		if layout[r][c] != '.' {
			found = true
			if layout[r][c] == '#' {
				occupied = 1
			}
		}
		r = r + d.drow
		c = c + d.dcol
	}
	return occupied
}

func nextState2(layout [][]rune, nextLayout [][]rune, rows int, cols int, row int, col int) bool {
	changed := false
	if layout[row][col] != '.' {
		occupied := 0
		for _, d := range directions {
			occupied += seeOccupied(layout, rows, cols, row, col, d)
		}
		if layout[row][col] == 'L' {
			if occupied == 0 {
				nextLayout[row][col] = '#'
				changed = true
			}
		} else if occupied > 4 {
			nextLayout[row][col] = 'L'
			changed = true
		}
	}
	return changed
}

func solve2(layout [][]rune, rows int, cols int) [][]rune {
	next, changed := nextLayout(layout, rows, cols, nextState2)
	for changed {
		next, changed = nextLayout(next, rows, cols, nextState2)
	}
	return next
}

func countOccupied(layout [][]rune) int {
	occupied := 0
	for r := range layout {
		for c := range layout[r] {
			if layout[r][c] == '#' {
				occupied++
			}
		}
	}
	return occupied
}

func printLayout(layout [][]rune) {
	for _, row := range layout {
		fmt.Println(string(row))
	}
}

func main() {
	file, err := os.Open("day11_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var layout [][]rune
	for scanner.Scan() {
		layout = append(layout, []rune(strings.TrimSpace(scanner.Text())))
	}
	rows := len(layout)
	cols := len(layout[0])

	layout1 := solve1(layout, rows, cols)
	//printLayout(layout1)
	fmt.Println("Answer 1:", countOccupied(layout1))

	layout2 := solve2(layout, rows, cols)
	//printLayout(layout2)
	fmt.Println("Answer 2:", countOccupied(layout2))
}
