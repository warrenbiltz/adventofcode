package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type position struct {
	x int
	y int
	o int
}

func (pos position) String() string {
	return fmt.Sprintf("%d | %d | %s | %d", pos.x, pos.y, string(pos.orientationStr()), pos.distance())
}

func (pos *position) orientationStr() rune {
	if pos.o == 0 {
		return 'N'
	} else if pos.o == 1 {
		return 'E'
	} else if pos.o == 2 {
		return 'S'
	} else {
		return 'W'
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (pos *position) distance() int {
	return abs(pos.x) + abs(pos.y)
}

func (pos *position) turn(o rune, d int) {
	nd := (d / 90) % 4
	if o == 'R' {
		pos.o = (pos.o + nd) % 4
	} else {
		pos.o = (pos.o - nd + 4) % 4
	}
}

func (pos *position) move(o rune, d int) {
	if o == 'N' {
		pos.y += d
	} else if o == 'S' {
		pos.y -= d
	} else if o == 'E' {
		pos.x += d
	} else if o == 'W' {
		pos.x -= d
	} else if o == 'F' {
		pos.move(pos.orientationStr(), d)
	} else {
		pos.turn(o, d)
	}
}

func (pos *position) moveStr(dir string) {
	o := rune(dir[0])
	d, _ := strconv.Atoi(dir[1:])
	pos.move(o, d)
}

func main() {
	file, err := os.Open("day12_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	pos := position{0, 0, 1}
	for scanner.Scan() {
		pos.moveStr(strings.TrimSpace(scanner.Text()))
	}
	fmt.Println("Ship:", pos)
}
