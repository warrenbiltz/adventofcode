package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

type Instruction struct {
	cmd string
	arg int
}

func (ins *Instruction) Invertable() bool {
	if ins.cmd == "acc" || ins.cmd == "nop" && ins.arg == 0 {
		return false
	}
	return true
}

func (ins *Instruction) Invert() {
	if ins.cmd == "jmp" {
		ins.cmd = "nop"
	} else if ins.cmd == "nop" {
		ins.cmd = "jmp"
	}
}

func parse(line string) Instruction {
	fields := strings.Fields(line)
	num, _ := strconv.Atoi(fields[1])
	return Instruction{fields[0], num}
}

func execute(instruction Instruction, acc int, line int) (int, int) {
	if instruction.cmd == "jmp" {
		return acc, line + instruction.arg
	} else if instruction.cmd == "acc" {
		return acc + instruction.arg, line + 1
	} else {
		return acc, line + 1
	}
}

func trace_cycle(execution_path map[int]int, last int, start int) ([]int, int) {
	var steps []int
	min := start
	if _, ok := execution_path[last]; !ok {
		steps = append(steps, start)
		next := execution_path[start]
		for next != start {
			if next < min {
				min = next
			}
			steps = append(steps, next)
			next = execution_path[next]
		}
	}
	return steps, min
}

func run(instructions []Instruction, execution_path map[int]int, start int) (int, int) {
	acc := 0
	line := start
	for line < len(instructions) {
		if _, ok := execution_path[line]; ok {
			break
		}
		prev := line
		acc, line = execute(instructions[line], acc, line)
		execution_path[prev] = line
	}
	return acc, line
}

func solve1(instructions []Instruction, execution_path map[int]int) (int, int) {
	return run(instructions, execution_path, 0)
}

func main() {
	file, err := os.Open("day8_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var instructions []Instruction
	for scanner.Scan() {
		instructions = append(instructions, parse(strings.TrimSpace(scanner.Text())))
	}
	file.Close()
	execution_path := make(map[int]int)
	ans1, last := solve1(instructions, execution_path)
	steps, cycle_min := trace_cycle(execution_path, len(instructions), last)
	fmt.Println("Answer 1: ", ans1)
	fmt.Println("Steps:", len(steps), "Cycle Min:", cycle_min)

	ans2 := 0
	inverted := 0
	for _, line := range(steps) {
		if instructions[line].Invertable() {
			instructions[line].Invert()
			execution_path := make(map[int]int)
			ans, _ := solve1(instructions, execution_path)
			if _, ok := execution_path[len(instructions)-1]; ok {
				ans2 = ans
				inverted = line
				break
			}
			instructions[line].Invert()
		}
	}
	fmt.Println("Answer 2:", ans2, "Inverted:", inverted, "To:", instructions[inverted])
}
