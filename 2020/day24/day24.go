package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type hexaPos struct {
	q int
	r int
}

func (p hexaPos) String() string {
	return fmt.Sprintf("(%d|%d)", p.q, p.r)
}

func (p hexaPos) add(o hexaPos) hexaPos {
	return hexaPos{p.q + o.q, p.r + o.r}
}

var neighbors = []hexaPos{hexaPos{1, -1}, hexaPos{1, 0}, hexaPos{0, 1}, hexaPos{-1, 1}, hexaPos{-1, 0}, hexaPos{0, -1}}
var stepMap = map[string]hexaPos{
	"ne": hexaPos{1, -1},
	"e":  hexaPos{1, 0},
	"se": hexaPos{0, 1},
	"sw": hexaPos{-1, 1},
	"w":  hexaPos{-1, 0},
	"nw": hexaPos{0, -1},
}

type tileMap map[hexaPos]bool

func flipTile(t hexaPos, blackTiles tileMap) {
	if _, found := blackTiles[t]; found {
		delete(blackTiles, t)
	} else {
		blackTiles[t] = true
	}
}

func parseDirections(line string) hexaPos {
	currPos := hexaPos{0, 0}
	for i := 0; i < len(line); i++ {
		if line[i] == 'n' || line[i] == 's' {
			currPos = currPos.add(stepMap[line[i:i+2]])
			i++
		} else {
			currPos = currPos.add(stepMap[string(line[i])])
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

func countBlack(t hexaPos, blackTiles tileMap) int {
	if _, found := blackTiles[t]; found {
		return 1
	}
	return 0
}

func countBlackNeighbors(t hexaPos, blackTiles tileMap) int {
	b := 0
	for _, s := range neighbors {
		n := t.add(s)
		b += countBlack(n, blackTiles)
		if b > 2 {
			break
		}
	}
	return b
}

func getNeighbors(t hexaPos, neighborMap tileMap) {
	neighborMap[t] = true
	for _, s := range neighbors {
		n := t.add(s)
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

	start := time.Now()
	blackTiles := solve1(directions)
	solve2(blackTiles)
	totalDuration := time.Since(start)
	fmt.Println("Solved in:", totalDuration)
}
