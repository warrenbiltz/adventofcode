package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type opstack []rune

func (s opstack) Push(v rune) opstack {
	return append(s, v)
}

func (s opstack) Pop() (opstack, rune) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s opstack) Last() rune {
	return s[len(s)-1]
}

type numstack []int

func (s numstack) Push(v int) numstack {
	return append(s, v)
}

func (s numstack) Pop() (numstack, int) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s numstack) Last() int {
	return s[len(s)-1]
}

func opOrder(op rune) int {
	////pt1
	//return 1
	//pt2
	if op == '+' {
		return 2
	} else if op == '*' {
		return 1
	}
	return 0
}

func applyOp(op rune, left int, right int) int {
	if op == '+' {
		return left + right
	} else if op == '*' {
		return left * right
	}
	return 0
}

func runeToInt(r rune) (int, bool) {
	num := int(r - '0')
	return num, num >= 0 && num < 10
}

func evalTopExp(nums numstack, ops opstack) (numstack, opstack) {
	var left, right int
	var op rune
	nums, right = nums.Pop()
	nums, left = nums.Pop()
	ops, op = ops.Pop()
	num := applyOp(op, left, right)
	nums = nums.Push(num)
	return nums, ops
}

func evalLine(line string) int {
	nums := make(numstack, 0)
	ops := make(opstack, 0)
	for _, r := range line {
		if r == ' ' {
			continue
		} else if r == '(' {
			ops = ops.Push(r)
		} else if r == ')' {
			for len(ops) > 0 && ops.Last() != '(' {
				nums, ops = evalTopExp(nums, ops)
			}
			ops, _ = ops.Pop()
		} else if num, isNum := runeToInt(r); isNum {
			nums = nums.Push(num)
		} else {
			nextOp := rune(r)
			for len(ops) > 0 && opOrder(ops.Last()) >= opOrder(nextOp) {
				nums, ops = evalTopExp(nums, ops)
			}
			ops = ops.Push(nextOp)
		}
	}
	for len(ops) > 0 {
		nums, ops = evalTopExp(nums, ops)
	}
	return nums.Last()
}

func main() {
	file, err := os.Open("day18_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	sum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		res := evalLine(line)
		fmt.Println(res, "=", line)
		sum += res
	}
	fmt.Println("Sum:", sum)
}
