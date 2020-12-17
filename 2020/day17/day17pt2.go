package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type pos struct {
	x int
	y int
	z int
	w int
}

type cuboid struct {
	min pos
	max pos
}

func (c *cuboid) update(p pos) {
	if p.x < c.min.x {
		c.min.x = p.x
	}
	if p.y < c.min.y {
		c.min.y = p.y
	}
	if p.z < c.min.z {
		c.min.z = p.z
	}
	if p.w < c.min.w {
		c.min.w = p.w
	}
	if p.x > c.max.x {
		c.max.x = p.x
	}
	if p.y > c.max.y {
		c.max.y = p.y
	}
	if p.z > c.max.z {
		c.max.z = p.z
	}
	if p.w > c.max.w {
		c.max.w = p.w
	}
}

func isCellActive(cells map[pos]bool, p pos) bool {
	_, found := cells[p]
	return found
}

func countNeighbors(cells map[pos]bool, p pos) int {
	neighbors := 0
	for dw := -1; dw < 2; dw++ {
		for dz := -1; dz < 2; dz++ {
			for dy := -1; dy < 2; dy++ {
				for dx := -1; dx < 2; dx++ {
					np := pos{p.x + dx, p.y + dy, p.z + dz, p.w + dw}
					if isCellActive(cells, np) {
						neighbors++
					}
				}
			}
		}
	}
	if isCellActive(cells, p) {
		neighbors--
	}
	return neighbors
}

func nextGen(cells map[pos]bool, c cuboid) (map[pos]bool, cuboid) {
	cube := cuboid{pos{0, 0, 0, 0}, pos{0, 0, 0, 0}}
	newCells := make(map[pos]bool)
	for w := c.min.w - 1; w <= c.max.w+1; w++ {
		for z := c.min.z - 1; z <= c.max.z+1; z++ {
			for y := c.min.y - 1; y <= c.max.y+1; y++ {
				for x := c.min.x - 1; x <= c.max.x+1; x++ {
					p := pos{x, y, z, w}
					neighbors := countNeighbors(cells, p)
					if isCellActive(cells, p) {
						if neighbors == 2 || neighbors == 3 {
							newCells[p] = true
							cube.update(p)
						}
					} else {
						if neighbors == 3 {
							newCells[p] = true
							cube.update(p)
						}
					}
				}
			}
		}
	}
	return newCells, cube
}

func prettyPrint(c cuboid, cells map[pos]bool) {
	fmt.Println("CUBE:", c)
	for w := c.min.w; w <= c.max.w; w++ {
		for z := c.min.z; z <= c.max.z; z++ {
			fmt.Println("z:", z, "w:", w)
			for y := c.min.y; y <= c.max.y; y++ {
				var row bytes.Buffer
				for x := c.min.x; x <= c.max.x; x++ {
					p := pos{x, y, z, w}
					if isCellActive(cells, p) {
						row.WriteString("#")
					} else {
						row.WriteString(".")
					}
				}
				fmt.Println(row.String())
			}
			fmt.Println("")
		}
	}
	fmt.Println("ALIVE:", len(cells))
}

func main() {
	file, err := os.Open("day17_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	cube := cuboid{pos{0, 0, 0, 0}, pos{0, 0, 0, 0}}
	cells := make(map[pos]bool)
	y := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for x, c := range line {
			if c == '#' {
				p := pos{x, y, 0, 0}
				cells[p] = true
				cube.update(p)
			}
		}
		y++
	}
	prettyPrint(cube, cells)

	generations := 6
	for g := 0; g < generations; g++ {
		cells, cube = nextGen(cells, cube)
	}
	prettyPrint(cube, cells)
}
