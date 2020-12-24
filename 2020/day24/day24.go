package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type cubePos struct {
	x int
	y int
	z int
}

func (p cubePos) String() string {
	return fmt.Sprintf("(%d|%d|%d)", p.x, p.y, p.z)
}

var neighbors = []string{"ne", "e", "se", "sw", "w", "nw"}

type tileMap map[cubePos]bool

func flipTile(t cubePos, blackTiles tileMap) {
	if _, found := blackTiles[t]; found {
		delete(blackTiles, t)
	} else {
		blackTiles[t] = true
	}
}

func parseStep(step string, from cubePos) cubePos {
	var dx, dy, dz int
	if step == "ne" {
		dx = 1
		dz = -1
	} else if step == "e" {
		dx = 1
		dy = -1
	} else if step == "se" {
		dy = -1
		dz = 1
	} else if step == "sw" {
		dx = -1
		dz = 1
	} else if step == "w" {
		dx = -1
		dy = 1
	} else { //nw
		dy = 1
		dz = -1
	}
	return cubePos{from.x + dx, from.y + dy, from.z + dz}
}

func parseDirections(line string) cubePos {
	currPos := cubePos{0, 0, 0}
	for i := 0; i < len(line); i++ {
		if line[i] == 'n' || line[i] == 's' {
			currPos = parseStep(line[i:i+2], currPos)
			i++
		} else {
			currPos = parseStep(string(line[i]), currPos)
		}
	}
	return currPos
}

func solve1(directions []string) tileMap {
	blackTiles := make(tileMap)
	for _, d := range directions {
		t := parseDirections(d)
		flipTile(t, blackTiles)
	}
	fmt.Println("Day 0:", len(blackTiles))
	return blackTiles
}

func countBlack(t cubePos, blackTiles tileMap) int {
	if _, found := blackTiles[t]; found {
		return 1
	}
	return 0
}

func countBlackNeighbors(t cubePos, blackTiles tileMap) int {
	b := 0
	for _, s := range neighbors {
		n := parseStep(s, t)
		b += countBlack(n, blackTiles)
		if b > 2 {
			break
		}
	}
	return b
}

func getNeighbors(t cubePos, neighborMap tileMap) {
	neighborMap[t] = true
	for _, s := range neighbors {
		n := parseStep(s, t)
		neighborMap[n] = true
	}
}

func nextTiles(prevGen tileMap) tileMap {
	candidates := make(tileMap)
	for t := range prevGen {
		getNeighbors(t, candidates)
	}

	newTiles := make(tileMap)
	for t := range candidates {
		b := countBlackNeighbors(t, prevGen)
		if countBlack(t, prevGen) == 0 {
			if b == 2 {
				newTiles[t] = true
			}
		} else if b == 1 || b == 2 {
			newTiles[t] = true
		}
	}
	return newTiles
}

func solve2(blackTiles tileMap) {
	const days = 100
	for d := 0; d < days; d++ {
		blackTiles = nextTiles(blackTiles)
	}
	fmt.Printf("Day %d: %d Tiles\n", days, len(blackTiles))
}

func main() {
	file, err := os.Open("day24_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var directions []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		directions = append(directions, line)
	}
	blackTiles := solve1(directions)
	solve2(blackTiles)
}
