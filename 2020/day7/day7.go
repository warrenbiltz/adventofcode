package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

func parse_bags(bag_rule string) map[string]int {
	bag_map := make(map[string]int)
	if !strings.Contains(bag_rule, "no other") {
		bags := strings.Split(bag_rule, ",")
		for _, bag := range bags {
			bag_fields := strings.Fields(bag)
			bag_color := bag_fields[1] + " " + bag_fields[2]
			bag_count, _ := strconv.Atoi(bag_fields[0])
			bag_map[bag_color] = bag_count
		}
	}
	return bag_map
}

func parse_rule(rules map[string]map[string]int, rule string) {
	rule_parts := strings.Split(rule, "contain")
	rule_id_fields := strings.Fields(rule_parts[0])
	rule_id := rule_id_fields[0] + " " + rule_id_fields[1]
	rules[rule_id] = parse_bags(rule_parts[1])
}

func search(rules map[string]map[string]int, start string, cache map[string]bool) bool {
	path_found := false
	if val, ok := cache[start]; ok {
		path_found = val
	} else {
		for k, _ := range rules[start] {
			if k == "shiny gold" || search(rules, k, cache) {
				path_found = true
				break
			}
		}
		cache[start] = path_found
	}
	return path_found
}

func solve_1(rules map[string]map[string]int) int {
	cache := make(map[string]bool)
	ans := 0
	for k, _ := range rules {
		if search(rules, k, cache) {
			ans += 1
		}
	}
	return ans
}


func count(rules map[string]map[string]int, start string, cache map[string]int) int {
	num_bags := 0
	if val, ok := cache[start]; ok {
		num_bags = val
	} else {
		for k, v := range rules[start] {
			num_bags += v * (count(rules, k, cache) + 1)
		}
		cache[start] = num_bags
	}
	return num_bags
}

func solve_2(rules map[string]map[string]int) int {
	cache := make(map[string]int)
	return count(rules, "shiny gold", cache)
}

func main() {
	file, err := os.Open("day7_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	rules := make(map[string]map[string]int)

	for scanner.Scan() {
		rule := strings.TrimSpace(scanner.Text())
		parse_rule(rules, rule)
	}

	fmt.Println("Answer 1: ", solve_1(rules))
	fmt.Println("Answer 2: ", solve_2(rules))
	file.Close()
}
