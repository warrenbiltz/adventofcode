package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type tile struct {
	id          int
	img         [][]rune
	outline     map[rune]uint //N, E, S, W, and n, e, s, w for flipped
	orientation []rune        //outline for up, right, down, left
}

func (t tile) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TILE %d:\n", t.id))
	sb.WriteString("Outlines:")
	for _, r := range "NnEeSsWw" {
		sb.WriteString(fmt.Sprintf(" %s: %d", string(r), t.outline[r]))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Orientation: %s\n", string(t.orientation)))
	for _, row := range t.img {
		sb.WriteString(string(row))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *tile) orient(orientation string) {
	t.orientation = []rune(orientation)
}

func (t *tile) orientWithOutlineOnSide(outline rune, side int) {
	for _, o := range possibleOrientations {
		if outline == rune(o[side]) {
			t.orient(o)
			return
		}
	}
}

const up = 0
const right = 1
const down = 2
const left = 3

func (t *tile) getOutline(side int) uint {
	return t.outline[t.orientation[side]]
}

func getMatchingSide(side int) int {
	return (side + 2) % 4
}

type idOutline struct {
	id      int
	outline rune
}

func (io idOutline) String() string {
	return fmt.Sprintf("%d|%s", io.id, string(io.outline))
}

var tiles = make(map[int]*tile)
var keys = make([]int, 0)
var possibleOrientations = []string{"NESW", "nWsE", "EsWn", "enws", "swne", "SeNw", "wNeS", "WSEN"}
var outlineToTiles = make(map[uint][]idOutline)
var gridSize = 0
var numTiles = 0

func registerNewTile(id int, img [][]rune) {
	outline := make(map[rune]uint)
	//north borders
	border := uint(0)
	borderFlip := uint(0)
	for i, r := range img[0] {
		val := 0
		if r == '#' {
			val = 1
		}
		border = border<<1 | uint(val)
		borderFlip = borderFlip | uint(val)<<i
	}
	outline['N'] = border
	outline['n'] = borderFlip
	outlineToTiles[border] = append(outlineToTiles[border], idOutline{id, 'N'})
	outlineToTiles[borderFlip] = append(outlineToTiles[borderFlip], idOutline{id, 'n'})
	//east borders
	border = uint(0)
	borderFlip = uint(0)
	last := len(img) - 1
	for i := 0; i <= last; i++ {
		val := 0
		if img[i][last] == '#' {
			val = 1
		}
		border = border<<1 | uint(val)
		borderFlip = borderFlip | uint(val)<<i
	}
	outline['E'] = border
	outline['e'] = borderFlip
	outlineToTiles[border] = append(outlineToTiles[border], idOutline{id, 'E'})
	outlineToTiles[borderFlip] = append(outlineToTiles[borderFlip], idOutline{id, 'e'})
	//south borders
	border = uint(0)
	borderFlip = uint(0)
	for i, r := range img[last] {
		val := 0
		if r == '#' {
			val = 1
		}
		border = border<<1 | uint(val)
		borderFlip = borderFlip | uint(val)<<i
	}
	outline['S'] = border
	outline['s'] = borderFlip
	outlineToTiles[border] = append(outlineToTiles[border], idOutline{id, 'S'})
	outlineToTiles[borderFlip] = append(outlineToTiles[borderFlip], idOutline{id, 's'})
	//west borders
	border = uint(0)
	borderFlip = uint(0)
	for i := 0; i <= last; i++ {
		val := 0
		if img[i][0] == '#' {
			val = 1
		}
		border = border<<1 | uint(val)
		borderFlip = borderFlip | uint(val)<<i
	}
	outline['W'] = border
	outline['w'] = borderFlip
	outlineToTiles[border] = append(outlineToTiles[border], idOutline{id, 'W'})
	outlineToTiles[borderFlip] = append(outlineToTiles[borderFlip], idOutline{id, 'w'})

	newTile := tile{id, img, outline, []rune("NESW")}
	tiles[id] = &newTile
	keys = append(keys, id)
}

func parseID(line string) int {
	idField := strings.Fields(line)[1]
	id, _ := strconv.Atoi(idField[0 : len(idField)-1])
	return id
}

func newGrid() [][]int {
	grid := make([][]int, gridSize)
	for i := range grid {
		grid[i] = make([]int, gridSize)
	}
	return grid
}

func printGrid(grid [][]int) {
	for _, row := range grid {
		fmt.Println(row)
	}
}

func solveHelper(grid [][]int, available map[int]bool) (bool, [][]int, map[int]bool) {
	if len(available) == 0 {
		return true, grid, available
	}
	tileNr := numTiles - len(available)
	row := tileNr / gridSize
	col := tileNr % gridSize
	solved := false
	for id := range available {
		delete(available, id)
		for _, o := range possibleOrientations {
			tiles[id].orient(o)
			leftFits := true
			if col > 0 {
				leftTile := grid[row][col-1]
				leftFits = tiles[leftTile].getOutline(right) == tiles[id].getOutline(left)
			}
			topFits := true
			if row > 0 {
				topTile := grid[row-1][col]
				topFits = tiles[topTile].getOutline(down) == tiles[id].getOutline(up)
			}
			if leftFits && topFits {
				grid[row][col] = id
				solved, grid, available = solveHelper(grid, available)
				if solved {
					break
				}
			}
		}
		if solved {
			break
		}
		available[id] = true
	}
	return solved, grid, available
}

func solveGrid() [][]int {
	solved := false
	grid := newGrid()
	available := make(map[int]bool)
	for _, k := range keys {
		available[k] = true
	}

	for _, k := range keys {
		delete(available, k)
		for _, o := range possibleOrientations {
			tiles[k].orient(o)
			grid[0][0] = k
			solved, grid, available = solveHelper(grid, available)
			if solved {
				break
			}
		}
		if solved {
			break
		}
		available[k] = true
	}

	fmt.Println("Solved:", solved)
	printGrid(grid)
	prod := int64(1)
	prod *= int64(grid[0][0])
	prod *= int64(grid[0][gridSize-1])
	prod *= int64(grid[gridSize-1][gridSize-1])
	prod *= int64(grid[gridSize-1][0])

	fmt.Println("Corners Prod:", prod)
	return grid
}

func printImg(img [][]rune) {
	for _, row := range img {
		fmt.Println(string(row))
	}
}

func (t *tile) stripImgBorder() {
	stripped := t.img[1 : len(t.img)-1]
	for r := range stripped {
		stripped[r] = stripped[r][1 : len(stripped[r])-1]
	}
	t.img = stripped
}

func rotateLeft(img [][]rune) [][]rune {
	imgSize := len(img)
	rotated := make([][]rune, imgSize)
	for i := range rotated {
		rotated[i] = make([]rune, imgSize)
	}
	for i := 0; i < imgSize; i++ {
		for j := 0; j < imgSize; j++ {
			rotated[i][j] = img[j][imgSize-i-1]
		}
	}
	return rotated
}

func rotateLeftTimes(img [][]rune, times int) [][]rune {
	for i := 0; i < times; i++ {
		img = rotateLeft(img)
	}
	return img
}

func flipRow(row []rune) []rune {
	for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
		row[i], row[j] = row[j], row[i]
	}
	return row
}

func flipVertical(img [][]rune) [][]rune {
	for r := range img {
		img[r] = flipRow(img[r])
	}
	return img
}

func align(img [][]rune, o string) [][]rune {
	//from NESW
	if o == "nWsE" {
		img = flipVertical(img)
	} else if o == "EsWn" {
		//rotateLeft
		img = rotateLeft(img)
	} else if o == "enws" {
		//rotateLeft, flipVertical
		img = rotateLeft(img)
		img = flipVertical(img)
	} else if o == "swne" {
		//rotateLeft x2
		img = rotateLeftTimes(img, 2)
	} else if o == "SeNw" {
		//rotateLeft x2, flipVertical
		img = rotateLeftTimes(img, 2)
		img = flipVertical(img)
	} else if o == "wNeS" {
		//rotateLeft x3
		img = rotateLeftTimes(img, 3)
	} else if o == "WSEN" {
		//rotateLeft x3 or rotateRight, flipVertical
		img = rotateLeftTimes(img, 3)
		img = flipVertical(img)
	}
	return img
}

func (t *tile) align() {
	fmt.Println("Orig:", t.id, string(t.orientation))
	printImg(t.img)
	t.img = align(t.img, string(t.orientation))
	fmt.Println("Align:", t.id, string(t.orientation))
	printImg(t.img)
}

func createFullImg(grid [][]int) [][]rune {
	fullImg := make([][]rune, 0)
	for _, row := range grid {
		for _, id := range row {
			tiles[id].stripImgBorder()
			tiles[id].align()
		}
	}
	tileImgSize := len(tiles[grid[0][0]].img)
	for _, row := range grid {
		for tileImgRow := 0; tileImgRow < tileImgSize; tileImgRow++ {
			fullImgRow := make([]rune, 0)
			for _, id := range row {
				fullImgRow = append(fullImgRow, tiles[id].img[tileImgRow]...)
			}
			fullImg = append(fullImg, fullImgRow)
		}
	}
	return fullImg
}

func countAlive(img [][]rune) int64 {
	alive := int64(0)
	for _, row := range img {
		for _, col := range row {
			if col == '#' {
				alive++
			}
		}
	}
	return alive
}

type gridPos struct {
	r int
	c int
}

func getMonsterPattern() []gridPos {
	file, err := os.Open("monster.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var monster []gridPos
	row := 0
	for scanner.Scan() {
		for col, r := range scanner.Text() {
			if r == '#' {
				monster = append(monster, gridPos{row, col})
			}
		}
		row++
	}
	fmt.Println(monster)
	return monster
}

func findMonsterAt(img [][]rune, monsterPattern []gridPos, pos gridPos) bool {
	for _, m := range monsterPattern {
		if img[pos.r+m.r][pos.c+m.c] != '#' {
			return false
		}
	}
	return true
}
func findMonsters(img [][]rune, monsterPattern []gridPos) []gridPos {
	var monsterPositions []gridPos

	for r := 0; r < len(img)-3; r++ {
		for c := 0; c < len(img)-20; c++ {
			pos := gridPos{r, c}
			if findMonsterAt(img, monsterPattern, pos) {
				monsterPositions = append(monsterPositions, pos)
			}
		}
	}
	return monsterPositions
}

func main() {
	file, err := os.Open("day20_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var id int = -1
	var img [][]rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			registerNewTile(id, img)
			id = -1
			img = make([][]rune, 0)
		} else if strings.Contains(line, "Tile") {
			id = parseID(line)
		} else {
			img = append(img, []rune(line))
		}
	}
	registerNewTile(id, img)
	sort.Ints(keys)

	numTiles = len(tiles)
	gridSize = int(math.Sqrt(float64(numTiles)))

	fmt.Printf("Got %d (%d) Tiles:\n\n", numTiles, gridSize)
	for _, k := range keys {
		fmt.Println(k)
	}
	// for _, t := range tiles {
	// 	fmt.Println(t)
	// }
	// for k, v := range outlineToTiles {
	// 	fmt.Println(k, ":", v)
	// }
	grid := solveGrid()
	fullImg := createFullImg(grid)
	alive := countAlive(fullImg)

	monsterPattern := getMonsterPattern()
	monsterSize := len(monsterPattern)
	var monsterPositions []gridPos

	for i := 0; i < 4; i++ {
		monsterPositions = findMonsters(fullImg, monsterPattern)
		if len(monsterPositions) > 0 {
			fmt.Println("FOUND", len(monsterPositions))
			break
		} else {
			fullImg = flipVertical(fullImg)
			monsterPositions = findMonsters(fullImg, monsterPattern)
			if len(monsterPositions) > 0 {
				fmt.Println("FOUND", len(monsterPositions))
				break
			} else {
				fullImg = flipVertical(fullImg)
			}
		}
		fullImg = rotateLeft(fullImg)
	}
	fmt.Println("FINAL IMG:")
	for _, row := range fullImg {
		fmt.Println(string(row))
	}
	numMonsters := len(monsterPositions)
	fmt.Println("Roughness: ", numMonsters, alive-int64(numMonsters*monsterSize))
}
