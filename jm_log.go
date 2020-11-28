package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func read_jm_log(filename string, jobs []string, alllogs *allLogs, wg *sync.WaitGroup) {
	loglines := make([]string, 0, 10000)

	file, err := os.Open(filename)
	if err == nil {
		fmt.Println("reading", filename, "...")
		reader := bufio.NewReaderSize(file, 16*1024*1024)
		for {
			line, err := reader.ReadString('\n')
			// line, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			for _, job := range jobs {
				pos := strings.Index(line, job)
				// filter out some substring matches
				if pos > 0 && strings.ContainsRune(":(", rune(line[pos-1])) {
					loglines = append(loglines, string(line))
				}
			} // job
		}
		file.Close()
	}
	//fmt.Println(loglines)
	alllogs.Append(loglines)
	wg.Done()
}
