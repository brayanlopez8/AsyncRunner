package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

//go run runner/ScriptRunner.go --curl='--location|--request|GET|https://gorest.co.in/public/v2/%s/%s'
//go run runner/ScriptRunner.go --curl='--location|--request|GET|https://gorest.co.in/public/v2/%s/%s' --osfile=AsyncRunner/test.csv
var OSfile = flag.String("osfile", "", "scan")
var CURL = flag.String("curl", "", "ConcurentRunner")

func main() {
	flag.Parse()
	csvLines, err := readFile()
	if err {
		return
	}
	var wg sync.WaitGroup
	for _, line := range csvLines {
		var params []string
		for i := 0; i < len(line); i++ {
			params = append(params, line[i])
		}
		wg.Add(1)
		go func(arg []string) {

			defer wg.Done()

			err := execute(arg)
			if err {
				return
			}
		}(params)

	}
	wg.Wait()
}

func readFile() ([][]string, bool) {
	csvFile, err := os.Open(*OSfile)
	if err != nil {
		fmt.Println(err)
		return nil, true
	}
	fmt.Println("Successfully Opened CSV file")

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	csvFile.Close()
	if err != nil {
		fmt.Println(err)
		return nil, true
	}
	return csvLines, false
}

func execute(arg []string) bool {
	fields := make([]interface{}, len(arg))
	for i, v := range arg {
		fields[i] = v
	}

	strcurl := fmt.Sprintf(*CURL, fields...)
	// fmt.Println(strcurl)
	params := strings.Split(strcurl, "|")
	curl := exec.Command("curl", params...)

	out, err := curl.Output()

	if err != nil {
		fmt.Println("error", err)
		return true
	}
	fmt.Println(string(out))
	return false
}
