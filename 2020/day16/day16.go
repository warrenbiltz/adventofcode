package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type interval struct {
	min int
	max int
}

type rule struct {
	name      string
	intervals []interval
}

func (in *interval) includes(val int) bool {
	if val < in.min || val > in.max {
		return false
	}
	return true
}

func parseInterval(intervalStr string) interval {
	intervalSplits := strings.Split(intervalStr, "-")
	min, _ := strconv.Atoi(intervalSplits[0])
	max, _ := strconv.Atoi(intervalSplits[1])
	return interval{min, max}
}

func parseRule(line string) rule {
	ruleSplits := strings.Split(line, ":")
	name := ruleSplits[0]
	fields := strings.Fields(strings.TrimSpace(ruleSplits[1]))
	return rule{
		name,
		[]interval{
			parseInterval(fields[0]),
			parseInterval(fields[2]),
		},
	}
}

func parseTicket(line string) []int {
	ticketSplits := strings.Split(line, ",")
	var ticket []int
	for _, v := range ticketSplits {
		val, _ := strconv.Atoi(v)
		ticket = append(ticket, val)
	}
	return ticket
}

func mergeIntervals(rules []rule) []interval {
	var allIntervals []interval
	for _, r := range rules {
		allIntervals = append(allIntervals, r.intervals...)
	}
	sort.SliceStable(allIntervals, func(i, j int) bool {
		if allIntervals[i].min == allIntervals[j].min {
			return allIntervals[i].max < allIntervals[j].max
		}
		return allIntervals[i].min < allIntervals[j].min
	})
	var merged []interval
	curInterval := allIntervals[0]
	for _, i := range allIntervals[1:] {
		if i.min > curInterval.max+1 {
			merged = append(merged, curInterval)
			curInterval = i
		} else if i.max > curInterval.max {
			curInterval.max = i.max
		}
	}
	merged = append(merged, curInterval)
	return merged
}

func isValidVal(intervals []interval, val int) bool {
	for _, in := range intervals {
		if in.includes(val) {
			return true
		}
	}
	return false
}

func isValidTicket(intervals []interval, ticket []int) bool {
	for _, v := range ticket {
		if !isValidVal(intervals, v) {
			return false
		}
	}
	return true
}

func valueSatisfiesRule(r rule, val int) bool {
	for _, in := range r.intervals {
		if in.includes(val) {
			return true
		}
	}
	return false
}

func isPossibleField(r rule, tickets [][]int, field int) bool {
	for _, t := range tickets {
		if !valueSatisfiesRule(r, t[field]) {
			return false
		}
	}
	return true
}

func getRemainingField(possibleFields []map[int]bool, pos int) int {
	for k := range possibleFields[pos] {
		return k
	}
	return -1
}

func determineFields(rules []rule, tickets [][]int) []int {
	var queue []int
	var possibleFields []map[int]bool
	for f, r := range rules {
		fieldMap := make(map[int]bool)
		for i := range rules {
			if isPossibleField(r, tickets, i) {
				fieldMap[i] = true
			}
		}
		possibleFields = append(possibleFields, fieldMap)
		if len(fieldMap) == 1 {
			queue = append(queue, f)
		}
	}

	for len(queue) > 0 {
		next := queue[0]
		field := getRemainingField(possibleFields, next)
		queue = queue[1:]
		for f, p := range possibleFields {
			if len(p) > 1 {
				delete(p, field)
				if len(p) == 1 {
					queue = append(queue, f)
				}
			}
		}
	}

	fields := make([]int, len(rules))
	for f := range possibleFields {
		fields[f] = getRemainingField(possibleFields, f)
	}
	return fields
}

func main() {
	file, err := os.Open("day16_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var rules []rule
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			break
		}
		rules = append(rules, parseRule(line))
	}
	scanner.Scan() //skip your ticket:
	scanner.Scan()
	myTicket := parseTicket(strings.TrimSpace(scanner.Text()))
	scanner.Scan() //skip next empty
	scanner.Scan() //skip nearby tickets:

	var tickets [][]int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			break
		}
		tickets = append(tickets, parseTicket(line))
	}

	mergedIntervals := mergeIntervals(rules)
	invalidSum := 0
	for _, t := range tickets {
		for _, v := range t {
			if !isValidVal(mergedIntervals, v) {
				invalidSum += v
			}
		}
	}
	var validTickets [][]int
	for _, t := range tickets {
		if isValidTicket(mergedIntervals, t) {
			validTickets = append(validTickets, t)
		}
	}

	fields := determineFields(rules, validTickets)
	fmt.Println("Fields:", fields)
	prod := 1
	for f, r := range rules {
		if strings.Contains(r.name, "departure") {
			prod *= myTicket[fields[f]]
		}
	}

	fmt.Println("Invalid Sum:", invalidSum)
	fmt.Println("Answer 2:", prod)
}
