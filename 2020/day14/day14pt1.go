package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseValue(valStr string, clearMask int64, valMask int64) int64 {
	value, _ := strconv.Atoi(valStr)
	return (int64(value) & clearMask) | valMask
}

func parseMem(memStr string) int64 {
	var mem int64
	fmt.Sscanf(memStr, "mem[%d]", &mem)
	return mem
}

func parseMask(maskStr string) (int64, int64) {
	var clearMask int64 = 0
	var valMask int64 = 0

	for _, c := range maskStr {
		clearBit := int64(0)
		valBit := int(0)
		if c == 'X' {
			clearBit = 1
		} else {
			valBit, _ = strconv.Atoi(string(c))
		}
		clearMask = (clearMask << 1) | clearBit
		valMask = (valMask << 1) | int64(valBit)
	}
	return clearMask, valMask
}

func main() {
	file, err := os.Open("day14_input.txt")
	if err != nil {
		log.Fatalf("failed to open")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	mem := make(map[int64]int64)
	var clearMask, valMask int64
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		cmd := strings.TrimSpace(line[0])
		arg := strings.TrimSpace(line[1])
		if cmd == "mask" {
			clearMask, valMask = parseMask(arg)
		} else {
			mem[parseMem(cmd)] = parseValue(arg, clearMask, valMask)
		}
	}
	var sum int64
	for _, v := range mem {
		sum += v
	}
	fmt.Println("MEM SUM:", sum)
}
