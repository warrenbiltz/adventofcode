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
	"time"
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

func (t *tile) hasFitsOnSide(side int) bool {
	if fits, found := outlineToTiles[t.getOutline(side)]; found {
		return len(fits) > 1
	}
	return false
}

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
	var leftTile, topTile int
	var leftOutline, topOutline uint
	candidates := make([]idOutline, 0)
	if col > 0 {
		leftTile = grid[row][col-1]
		leftOutline = tiles[leftTile].getOutline(right)
		candidates = append(candidates, outlineToTiles[leftOutline]...)
	}
	if row > 0 {
		topTile = grid[row-1][col]
		topOutline = tiles[topTile].getOutline(down)
		if col == 0 {
			candidates = append(candidates, outlineToTiles[topOutline]...)
		}
	}

	for _, c := range candidates {
		if _, avail := available[c.id]; avail {
			delete(available, c.id)
			if col > 0 {
				tiles[c.id].orientWithOutlineOnSide(c.outline, left)
			} else if row > 0 {
				tiles[c.id].orientWithOutlineOnSide(c.outline, up)
			}
			leftFits := true
			if col > 0 {
				leftFits = leftOutline == tiles[c.id].getOutline(left)
			}
			topFits := true
			if row > 0 {
				topFits = topOutline == tiles[c.id].getOutline(up)
			}
			if leftFits && topFits {
				grid[row][col] = c.id
				solved, grid, available = solveHelper(grid, available)
				if solved {
					break
				}
			}
			available[c.id] = true
		}
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
			if tiles[k].hasFitsOnSide(down) && tiles[k].hasFitsOnSide(right) {
				grid[0][0] = k
				solved, grid, available = solveHelper(grid, available)
				if solved {
					break
				}
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
	t.img = align(t.img, string(t.orientation))
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

type monster struct {
	positions []gridPos
	maxPos    gridPos
	size      int
}

func getMonster() monster {
	file, err := os.Open("monster.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var rowMax = 0
	var colMax = 0
	var positions []gridPos
	row := 0
	for scanner.Scan() {
		for col, r := range scanner.Text() {
			if r == '#' {
				positions = append(positions, gridPos{row, col})
				if row > rowMax {
					rowMax = row
				}
				if col > colMax {
					colMax = col
				}
			}
		}
		row++
	}
	return monster{positions, gridPos{rowMax, colMax}, len(positions)}
}

func rotateMonster(m monster) monster {
	const monsterImgSize = 20
	var rowMax = 0
	var colMax = 0
	var rowMin = monsterImgSize
	var colMin = monsterImgSize

	var newPositions []gridPos
	for _, p := range m.positions {
		r := p.c
		c := monsterImgSize - p.r - 1
		if r < rowMin {
			rowMin = r
		}
		if r > rowMax {
			rowMax = r
		}
		if c < colMin {
			colMin = c
		}
		if c > colMax {
			colMax = c
		}
		newPositions = append(newPositions, gridPos{r, c})
	}

	for i := range newPositions {
		newPositions[i] = gridPos{newPositions[i].r - rowMin, newPositions[i].c - colMin}
	}
	return monster{newPositions, gridPos{rowMax - rowMin, colMax - colMin}, m.size}
}

func flipMonster(m monster) monster {
	var newPositions []gridPos
	for _, p := range m.positions {
		c := m.maxPos.c - p.c
		newPositions = append(newPositions, gridPos{p.r, c})
	}
	return monster{newPositions, m.maxPos, m.size}
}

func findMonsterAt(img [][]rune, m monster, pos gridPos) bool {
	for _, m := range m.positions {
		if img[pos.r+m.r][pos.c+m.c] != '#' {
			return false
		}
	}
	return true
}
func findMonsters(img [][]rune, m monster) []gridPos {
	var monsterPositions []gridPos

	for r := 0; r+m.maxPos.r < len(img); r++ {
		for c := 0; c+m.maxPos.c < len(img); c++ {
			pos := gridPos{r, c}
			if findMonsterAt(img, m, pos) {
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

	start := time.Now()
	grid := solveGrid()
	gridSolveTime := time.Since(start)
	fullImg := createFullImg(grid)
	alive := countAlive(fullImg)

	monster := getMonster()
	var monsterPositions []gridPos
	for i := 0; i < 4; i++ {
		monsterPositions = findMonsters(fullImg, monster)
		if len(monsterPositions) > 0 {
			break
		} else {
			monster = flipMonster(monster)
			monsterPositions = findMonsters(fullImg, monster)
			if len(monsterPositions) > 0 {
				break
			} else {
				monster = flipMonster(monster)
			}
		}
		monster = rotateMonster(monster)
	}
	// fmt.Println("FINAL IMG:")
	// for _, row := range fullImg {
	// 	fmt.Println(string(row))
	// }
	numMonsters := len(monsterPositions)
	fmt.Println("Roughness: ", numMonsters, alive-int64(numMonsters*monster.size))

	totalDuration := time.Since(start)
	fmt.Println("Solved in:", gridSolveTime, totalDuration)
}
