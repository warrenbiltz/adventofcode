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

type rule struct {
	patterns [][]int
	char     rune
}

func (r rule) String() string {
	if r.char != ' ' {
		return fmt.Sprintf("%s", string(r.char))
	}
	return fmt.Sprint(r.patterns)
}

var rules = make(map[int]rule)
var keys = make([]int, 0)

func parsePattern(p string) []int {
	var pattern []int
	for _, s := range strings.Fields(strings.TrimSpace(p)) {
		num, _ := strconv.Atoi(s)
		pattern = append(pattern, num)
	}
	return pattern
}

func parseRule(line string) {
	ruleSplits := strings.Split(line, ":")
	ruleNr, _ := strconv.Atoi(ruleSplits[0])
	keys = append(keys, ruleNr)
	pattern := strings.TrimSpace(ruleSplits[1])
	if pattern[0] == '"' {
		rules[ruleNr] = rule{nil, rune(pattern[1])}
	} else {
		var patterns [][]int
		for _, p := range strings.Split(pattern, "|") {
			patterns = append(patterns, parsePattern(p))
		}
		rules[ruleNr] = rule{patterns, ' '}
	}
}

func satisfiesRule(ruleNr int, str []rune, canTerminate bool) (bool, int) {
	if len(str) == 0 {
		return canTerminate, 0
	} else if rules[ruleNr].char != ' ' {
		if str[0] == rules[ruleNr].char {
			return true, 1
		}
		return false, 0
	}
	for _, subPattern := range rules[ruleNr].patterns {
		lengthMatched := 0
		matched := false
		for s, subRule := range subPattern {
			canTerminate := false
			if subRule == 31 {
				canTerminate = subPattern[s-1] == 11
			}
			subMatch, l := satisfiesRule(subRule, str[lengthMatched:], canTerminate)
			if subMatch {
				matched = subMatch
				lengthMatched += l
			} else {
				matched, lengthMatched = false, 0
				break
			}
		}
		if matched {
			return matched, lengthMatched
		}
	}
	return false, 0
}

func satisfies(line string) (bool, int) {
	matched, l := satisfiesRule(0, []rune(line), false)
	return matched && l == len(line), l
}

func modRules() {
	//rule 8
	var pattern8 [][]int
	pattern8 = append(pattern8, []int{42})
	pattern8 = append(pattern8, []int{42, 8})
	rules[8] = rule{pattern8, ' '}
	//rule 11
	var pattern11 [][]int
	pattern11 = append(pattern11, []int{42, 31})
	pattern11 = append(pattern11, []int{42, 11, 31})
	rules[11] = rule{pattern11, ' '}
}

func main() {
	file, err := os.Open("day19_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			break
		}
		parseRule(line)
	}
	modRules()
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Println(k, rules[k])
	}
	total := 0
	countZero := 0
	var matches []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matched, l := satisfies(line)
		fmt.Println("Matched?:", matched, l, line)
		if matched {
			matches = append(matches, line)
			countZero++
		}
		total++
	}
	fmt.Println("0 Rule Matches:", countZero, "out of", total)
}
