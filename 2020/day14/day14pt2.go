package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseMem(memStr string) int64 {
	var mem int64
	fmt.Sscanf(memStr, "mem[%d]", &mem)
	return mem
}

func ithBit(num int64, bitPos int64) int64 {
	if (num & bitPos) > 0 {
		return 1
	}
	return 0
}

func getMasks(maskStr string) (int64, []int64) {
	var clearMask int64 = 0
	var valMask int64 = 0
	var xPositions []int64
	twoPowX := int64(1)
	for i, c := range maskStr {
		clearBit := int64(1)
		valBit := int(0)
		if c == 'X' {
			clearBit = 0
			xPositions = append(xPositions, int64(35-i))
			twoPowX = twoPowX * 2
		} else {
			valBit, _ = strconv.Atoi(string(c))
		}
		clearMask = (clearMask << 1) | clearBit
		valMask = (valMask << 1) | int64(valBit)
	}
	var addressMasks []int64
	for n := int64(0); n < twoPowX; n++ {
		bitPos := int64(1)
		addressMask := int64(0)
		for _, x := range xPositions {
			bit := ithBit(n, bitPos)
			addressMask = addressMask | (bit << x)
			bitPos = bitPos * 2
		}
		addressMasks = append(addressMasks, addressMask|valMask)
	}

	return clearMask, addressMasks
}

func getAddresses(clearMask int64, addressMasks []int64, main int64) []int64 {
	var addresses []int64
	base := main & clearMask
	for _, m := range addressMasks {
		addresses = append(addresses, base|m)
	}
	return addresses
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
	var clearMask int64
	var addressMasks []int64
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		cmd := strings.TrimSpace(line[0])
		arg := strings.TrimSpace(line[1])
		if cmd == "mask" {
			clearMask, addressMasks = getMasks(arg)
		} else {
			main := parseMem(cmd)
			value, _ := strconv.Atoi(arg)
			addresses := getAddresses(clearMask, addressMasks, main)
			for _, m := range addresses {
				mem[m] = int64(value)
			}
		}
	}
	var sum int64
	for _, v := range mem {
		sum += v
	}
	fmt.Println("MEM SUM:", sum)
}
