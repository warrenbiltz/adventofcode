package main

import (
	"container/list"
	"fmt"
	"strings"
)

var cups *list.List
var cupMap = make([]*list.Element, 10)
var pickups []*list.Element
var availMap = make([]bool, 10)

func printList(start *list.Element) string {
	var sb strings.Builder
	count := 0
	for e := start; e != nil && count < cups.Len(); e = next(e) {
		sb.WriteString(fmt.Sprintf("%d", e.Value))
		count++
	}
	return sb.String()
}

func next(e *list.Element) *list.Element {
	n := e.Next()
	if n == nil {
		n = cups.Front()
	}
	return n
}

func playRound(start *list.Element) *list.Element {
	//pickup next 3 at current start
	pickups = make([]*list.Element, 3)
	var n = start
	for i := 0; i < 3; i++ {
		n = next(n)
		pickups[i] = n
	}
	for _, e := range pickups {
		availMap[e.Value.(int)] = false
	}
	//determine destination
	d := start.Value.(int) - 1
	for !availMap[d] && d > 0 {
		d--
	}
	if d == 0 {
		d = 9
		for !availMap[d] && d > 0 {
			d--
		}
	}
	dest := cupMap[d]
	for _, e := range pickups {
		cups.MoveAfter(e, dest)
		availMap[e.Value.(int)] = true
		dest = e
	}

	return next(start)
}

func play(rounds int) {
	curr := cups.Front()
	for i := 0; i < rounds; i++ {
		curr = playRound(curr)
		fmt.Printf("Round %d: %s\n", i, printList(curr))
	}
}

var testGame = []int{3, 8, 9, 1, 2, 5, 4, 6, 7}
var realGame = []int{3, 9, 8, 2, 5, 4, 7, 1, 6}

func main() {
	g := realGame
	cups = list.New()
	availMap[0] = false
	for _, n := range g {
		e := cups.PushBack(n)
		cupMap[n] = e
		availMap[n] = true
	}
	play(100)
	fmt.Printf("Final: %s\n", printList(cupMap[1]))
}
