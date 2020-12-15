package main

import (
	"fmt"
	"strconv"
	"strings"
)

type bus struct {
	id     int64
	offset int64
}

var test = "7,13,x,x,59,x,31,19"
var input = "19,x,x,x,x,x,x,x,x,41,x,x,x,37,x,x,x,x,x,367,x,x,x,x,x,x,x,x,x,x,x,x,13,x,x,x,17,x,x,x,x,x,x,x,x,x,x,x,29,x,373,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,x,23"

func isValid(schedule []bus, start int64) bool {
	for _, b := range schedule {
		if (start+b.offset)%b.id != 0 {
			return false
		}
	}
	return true
}

func getGroups(schedule []bus) [][]bus {
	edges := make([][]int, len(schedule))
	for i := range schedule {
		base := schedule[i]
		for j := i + 1; j < len(schedule); j++ {
			b := schedule[j]
			if (b.offset-base.offset)%base.id == 0 {
				edges[i] = append(edges[i], j)
				edges[j] = append(edges[j], i)
			}
		}
	}
	visited := make([]bool, len(schedule))
	var groups [][]bus
	for i := range schedule {
		if !visited[i] {
			visited[i] = true
			group := addToGroup(schedule, edges, visited, []bus{}, i)
			groups = append(groups, group)
		}
	}
	return groups
}

func addToGroup(schedule []bus, edges [][]int, visited []bool, group []bus, i int) []bus {
	group = append(group, schedule[i])
	for _, j := range edges[i] {
		if !visited[j] {
			visited[j] = true
			group = addToGroup(schedule, edges, visited, group, j)
		}
	}
	return group
}

func consolidateSchedule(groups [][]bus) []bus {
	var consolidated []bus
	for _, g := range groups {
		base := int64(1)
		offset := int64(0)
		for _, b := range g {
			base *= b.id
			if offset < b.offset {
				offset = b.offset
			}
		}
		consolidated = append(consolidated, bus{base, offset})
	}
	return consolidated
}

func solve1(schedule []bus) int64 {
	groups := getGroups(schedule)
	fmt.Println("Groups:")
	for _, g := range groups {
		fmt.Println(g)
	}
	consolidated := consolidateSchedule(groups)
	fmt.Println("Consolidated:")
	for _, c := range consolidated {
		fmt.Println(c)
	}

	maxBus := consolidated[0]
	for _, c := range consolidated {
		if c.id > maxBus.id {
			maxBus = c
		}
	}
	fmt.Println("Max:", maxBus)

	maxID := int64(maxBus.id)
	time := int64(100000000000000/maxID)*maxID - maxBus.offset
	found := false
	for !found {
		time += maxID
		found = isValid(consolidated, time)
	}
	return time
}

func solve2(schedule []bus) int64 {
	acc := schedule[0]
	for _, b := range schedule[1:] {
		sync := acc.offset
		for (sync+b.offset)%b.id != 0 {
			sync += acc.id
		}
		acc.id = acc.id * b.id
		acc.offset = sync
	}
	return acc.offset
}

func main() {
	scheduleStr := strings.Split(input, ",")
	var schedule []bus

	for i, s := range scheduleStr {
		if s != "x" {
			id, _ := strconv.Atoi(s)
			schedule = append(schedule, bus{int64(id), int64(i)})
		}
	}
	fmt.Println("Solve 1:", solve1(schedule))
	fmt.Println("Solve 2:", solve2(schedule))
}
