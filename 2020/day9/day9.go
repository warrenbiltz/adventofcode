package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"math"
)

func add(num_count map[int64]int, num int64) {
	if count, found := num_count[num]; found {
		num_count[num] = count + 1;
		fmt.Println("Duplicate:", num)
	} else {
		num_count[num] = 1;
	}
}

func remove(num_count map[int64]int, num int64) {
	if count, found := num_count[num]; found {
		if count == 1 {
			delete(num_count, num)
		} else {
			num_count[num] = count - 1;
		}
	}
}

func is_valid(num_count map[int64]int, num int64) bool {
	for n1, v1 := range num_count {
		n2 := num - n1
		if n1 == n2 {
			if v1 > 1 {
				return true
			} 
		} else {
			if _, found := num_count[n2]; found {
				return true
			}
		}
	}
	return false
}

func solve1(numbers []int64) int64 {
	num_count := make(map[int64]int)
	//preamble
	for i:= 0; i < 25; i++ {
		add(num_count, numbers[i])
	}

	for i := 25; i < len(numbers); i++ {
		if !is_valid(num_count, numbers[i]) {
			return numbers[i]
		}
		remove(num_count, numbers[i-25])
		add(num_count, numbers[i])
	}
	return -1
}

func find_sequence_sum(numbers []int64, num int64) (int, int) {
	start := 0
	end := 1
	sliding_sum := numbers[start] + numbers[end]
	for end < len(numbers)-1 {
		if sliding_sum == num {
			break
		} else if sliding_sum < num {
			end++
			sliding_sum += numbers[end]
		} else {
			sliding_sum -= numbers[start]
			start++
			if end == start {
				end++
				sliding_sum += numbers[end]
			}
		}
	}
	return start, end
}

func get_min_max_sum(numbers []int64, start int, end int) (int64, int64, int64) {
	min := int64(math.MaxInt64)
	max := int64(0)
	for i := start; i <= end; i++ {
		if min > numbers[i] {
			min = numbers[i]
		}
		if max < numbers[i] {
			max = numbers[i]
		}
	}
	return min, max, min + max
}

func main() {
	file, err := os.Open("day9_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var numbers []int64
	for scanner.Scan() {
		num, _ := strconv.ParseInt(strings.TrimSpace(scanner.Text()), 10, 64)
		numbers = append(numbers, num)
	}
	file.Close()

	ans1 := solve1(numbers)
	start, end := find_sequence_sum(numbers, ans1)
	min, max, sum := get_min_max_sum(numbers, start, end)

	var check_sum int64
	for i := start; i <= end; i++ {
		check_sum += numbers[i]
	}
	fmt.Println("Answer 1:", ans1)
	fmt.Println("find_sequence_sum:", start, end, check_sum)
	fmt.Println("Answer 2:", min, max, sum)
}