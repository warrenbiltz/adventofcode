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
}

func (pos position) String() string {
	return fmt.Sprintf("%d | %d | %d", pos.x, pos.y, pos.distance())
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

func (pos *position) turnLeft() {
	x := pos.x
	pos.x = -pos.y
	pos.y = x
}

func (pos *position) turnRight() {
	x := pos.x
	pos.x = pos.y
	pos.y = -x
}

func (pos *position) turn(o rune, d int) {
	nd := (d / 90) % 4
	if nd%2 == 0 {
		pos.x = -pos.x
		pos.y = -pos.y
	} else if o == 'R' {
		if nd == 1 {
			pos.turnRight()
		} else {
			pos.turnLeft()
		}
	} else if o == 'L' {
		if nd == 1 {
			pos.turnLeft()
		} else {
			pos.turnRight()
		}
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
	} else {
		pos.turn(o, d)
	}
}

func (pos *position) follow(wp position, n int) {
	pos.x += n * wp.x
	pos.y += n * wp.y
}

func parse(dir string) (rune, int) {
	o := rune(dir[0])
	d, _ := strconv.Atoi(dir[1:])
	return o, d
}

func main() {
	file, err := os.Open("day12_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	waypoint := position{10, 1}
	ship := position{0, 0}
	for scanner.Scan() {
		o, d := parse(strings.TrimSpace(scanner.Text()))
		if o != 'F' {
			waypoint.move(o, d)
		} else {
			ship.follow(waypoint, d)
		}
	}
	fmt.Println("Waypoint:", waypoint)
	fmt.Println("Ship:", ship)
}
