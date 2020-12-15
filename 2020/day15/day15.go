package main

import "fmt"

type history struct {
	last int
	prev int
}

func (h *history) update(seen int) {
	h.prev = h.last
	h.last = seen
}

func (h *history) age() int {
	return h.last - h.prev
}

func addOrUpdate(cache map[int]*history, num int, turn int) {
	if _, seen := cache[num]; seen {
		cache[num].update(turn)
	} else {
		newHistory := history{turn, turn}
		cache[num] = &newHistory
	}
}

func main() {
	var input = []int{0, 13, 1, 8, 6, 15}
	var nMax = 30000000

	cache := make(map[int]*history)
	for i, x := range input {
		addOrUpdate(cache, x, i)
	}

	n := len(input)
	lastNr := input[n-1]
	for n < nMax {
		nextNr := cache[lastNr].age()
		addOrUpdate(cache, nextNr, n)
		lastNr = nextNr
		n++
	}
	fmt.Println("Answer:", n, lastNr)
}
