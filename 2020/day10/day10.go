package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"sort"
)

func solve2(numbers []int) int64 {
	ways := make([]int64, len(numbers))
	ways[0] = 1
	for i := 1; i < len(numbers); i++ {
		sum := int64(0);
		for j:= i-3; j < i; j++ {
			if j >= 0 && numbers[i] - numbers[j] <= 3 {
				sum += ways[j]
			}
		}
		if numbers[i] < 4 {
			sum += 1
		}
		ways[i] = sum
	}
	fmt.Println("ways:", ways)
	return ways[len(numbers)-1]
}

func solve1(numbers []int) (int, int, int) {
	ones := 1
	threes := 1
	for i:=1; i < len(numbers); i++ {
		diff := numbers[i] - numbers[i-1];
		if diff == 1 {
			ones++
		} else if diff == 3 {
			threes++
		}
	}
	return ones, threes, ones * threes
}
func main() {
	file, err := os.Open("day10_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var numbers []int
	for scanner.Scan() {
		num, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		numbers = append(numbers, num)
	}
	file.Close()
	sort.Ints(numbers)
	fmt.Println(numbers)

	ones, threes, ans1 := solve1(numbers)
	ans2 := solve2(numbers)

	fmt.Println("Answer 1:", ones, threes, ans1)
	fmt.Println("Answer 2:", ans2)
}