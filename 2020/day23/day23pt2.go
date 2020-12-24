package main

import (
	"container/list"
	"fmt"
	"strings"
	"time"
)

var numCups = 1000000
var cups *list.List
var cupMap = make([]*list.Element, numCups+1)
var pickups []*list.Element
var pickupMap = make(map[int]bool)

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
		pickupMap[e.Value.(int)] = true
	}
	//determine destination
	d := start.Value.(int) - 1
	_, found := pickupMap[d]
	for found && d > 0 {
		d--
		_, found = pickupMap[d]
	}
	if d == 0 {
		d = numCups
		_, found = pickupMap[d]
		for found && d > 0 {
			d--
			_, found = pickupMap[d]
		}
	}
	dest := cupMap[d]
	for _, e := range pickups {
		cups.MoveAfter(e, dest)
		delete(pickupMap, e.Value.(int))
		dest = e
	}
	return next(start)
}

func play(rounds int) {
	curr := cups.Front()
	for i := 0; i < rounds; i++ {
		curr = playRound(curr)
	}
}

var testGame = []int{3, 8, 9, 1, 2, 5, 4, 6, 7}
var realGame = []int{3, 9, 8, 2, 5, 4, 7, 1, 6}

func main() {
	g := realGame
	cups = list.New()
	for _, n := range g {
		e := cups.PushBack(n)
		cupMap[n] = e
	}
	for n := cups.Len() + 1; n <= numCups; n++ {
		e := cups.PushBack(n)
		cupMap[n] = e
	}
	start := time.Now()
	play(10000000)
	e1 := next(cupMap[1])
	e2 := next(e1)
	e1v := e1.Value.(int)
	e2v := e2.Value.(int)
	//fmt.Printf("Answer 1: %s\n", printList(cupMap[1]))
	fmt.Printf("Answer 2: %d | %d | %d\n", e1v, e2v, e1v*e2v)
	totalDuration := time.Since(start)
	fmt.Println("Solved in:", totalDuration)
}
