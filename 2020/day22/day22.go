package main

import (
	"fmt"
	"strings"
)

type hand []int
type game struct {
	p1   hand
	p2   hand
	seen map[string]bool
}

func newGame(p1 hand, p2 hand) game {
	//copy hand
	p1cp := make(hand, len(p1))
	p2cp := make(hand, len(p2))
	copy(p1cp, p1)
	copy(p2cp, p2)
	return game{p1cp, p2cp, make(map[string]bool)}
}

func (h hand) hashStr() string {
	var sb strings.Builder
	for _, n := range h {
		sb.WriteString(fmt.Sprintf(":%d", n))
	}
	return sb.String()
}

func (g *game) hashStr() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("P1:%s|", g.p1.hashStr()))
	sb.WriteString(fmt.Sprintf("P2:%s", g.p2.hashStr()))
	return sb.String()
}

func (g *game) over() bool {
	if len(g.p1) == 0 || len(g.p2) == 0 {
		return true
	}
	key := g.hashStr()
	_, found := g.seen[key]
	if !found {
		g.seen[key] = true
	}
	return found
}

func getScore(h hand) int {
	factor := len(h)
	score := 0
	for _, n := range h {
		score += factor * n
		factor--
	}
	return score
}

func (g *game) getWinner() (int, int, hand) {
	if len(g.p1) == 0 {
		return 2, getScore(g.p2), g.p2
	}
	return 1, getScore(g.p1), g.p1
}

func (g *game) playRound() {
	p1c := g.p1[0]
	g.p1 = g.p1[1:]
	p2c := g.p2[0]
	g.p2 = g.p2[1:]
	if p1c > p2c {
		g.p1 = append(g.p1, p1c, p2c)
	} else {
		g.p2 = append(g.p2, p2c, p1c)
	}
}

func (g *game) play() {
	for !g.over() {
		g.playRound()
	}
}

func (g *game) playRoundRec() {
	p1c := g.p1[0]
	g.p1 = g.p1[1:]
	p2c := g.p2[0]
	g.p2 = g.p2[1:]
	if p1c <= len(g.p1) && p2c <= len(g.p2) {
		//play new subgame
		subGame := newGame(g.p1[:p1c], g.p2[:p2c])
		subGame.playRec()
		winner, _, _ := subGame.getWinner()
		if winner == 1 {
			g.p1 = append(g.p1, p1c, p2c)
		} else {
			g.p2 = append(g.p2, p2c, p1c)
		}
	} else {
		if p1c > p2c {
			g.p1 = append(g.p1, p1c, p2c)
		} else {
			g.p2 = append(g.p2, p2c, p1c)
		}
	}
}

func (g *game) playRec() {
	for !g.over() {
		g.playRoundRec()
	}
}

var infiniteGame = newGame(hand{43, 19}, hand{2, 29, 14})
var testGame = newGame(hand{9, 2, 6, 3, 1}, hand{5, 8, 4, 7, 10})
var realGame = newGame([]int{5, 20, 28, 30, 48, 7, 41, 24, 29, 8, 37, 32, 16, 17, 34, 27, 46, 43, 14, 49, 35, 11, 6, 38, 1},
	[]int{22, 18, 50, 31, 12, 13, 33, 39, 45, 21, 19, 26, 44, 10, 42, 3, 4, 15, 36, 2, 40, 47, 9, 23, 25})

func main() {
	g := newGame(realGame.p1, realGame.p2)
	g.play()
	w, s, h := g.getWinner()
	fmt.Println("Winner", w, s, h)

	recGame := newGame(realGame.p1, realGame.p2)
	recGame.playRec()
	w, s, h = recGame.getWinner()
	fmt.Println("Winner Rec", w, s, h)
}
