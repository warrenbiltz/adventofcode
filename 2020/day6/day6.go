package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func count_chars(count map[rune]int, line string) {
	for _, c := range line {
		count[rune(c)] += 1
	}
}

func count_unanimous(count map[rune]int, n int) int {
	res := 0
	for _, v := range count {
		if v == n {
			res++
		}
	}
	return res
}

func main() {
	file, err := os.Open("day6_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	count := make(map[rune]int)
	group_count := 0
	sum1 := 0
	sum2 := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			sum1 += len(count)
			sum2 += count_unanimous(count, group_count)
			count = make(map[rune]int)
			group_count = 0
		} else {
			group_count++
		}
		count_chars(count, line)
	}
	sum1 += len(count)
	sum2 += count_unanimous(count, group_count)
	fmt.Println("Answer 1: ", sum1)
	fmt.Println("Answer 2: ", sum2)

	file.Close()
}
