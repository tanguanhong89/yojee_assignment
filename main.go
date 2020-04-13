package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() { // example yojee cleaned.csv 11.552931 104.933636 4
	args := os.Args[1:]
	_ = args

	x, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		panic("Please use a float number for your 2nd argument (X coordinate of starting point)")
	}
	y, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		panic("Please use a float number for your 3rd argument (Y coordinate of starting point)")
	}
	workerCount, err := strconv.Atoi(args[3])
	if err != nil {
		panic("Please use a int number for your 4th argument (Number of workers)")
	}
	start := []float64{x, y}
	dest := parse2DStringAs2DFloats(readCSV(args[0]))
	findPaths(start, dest, workerCount)
}

func readCSV(path string) [][]string {
	file, err := os.Open(path)
	res := [][]string{}
	h := func(x, y string) []string { return strings.Split(x, y) }
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, h(scanner.Text(), ","))
	}
	return res
}

func parse2DStringAs2DFloats(s [][]string) [][]float64 {
	res := [][]float64{}
	for d := range s {
		_t := []float64{}
		for dd := range s[d] {
			v, _ := strconv.ParseFloat(s[d][dd], 32)
			_t = append(_t, v)
		}
		res = append(res, _t)
	}
	return res
}
