package main

import "fmt"

const t1 = 5764801
const t2 = 17807724

const p1 = 10943862
const p2 = 12721030

const subject = 7
const base = 20201227

type modExpMap map[int]int //map from mod to loop size

func determineLoopSize(target int) int {
	acc := 1
	l := 0
	for acc != target {
		acc = (acc * subject) % base
		l++
	}
	return l
}
func solve1(one int, two int) (int, int, int) {

	l1 := determineLoopSize(one)
	l2 := determineLoopSize(two)
	acc := 1
	l := 0
	for l < l2 {
		acc = (acc * one) % base
		l++
	}
	return l1, l2, acc
}

func main() {
	p, k, enc := solve1(p1, p2)
	fmt.Println(p, k, enc)
}
