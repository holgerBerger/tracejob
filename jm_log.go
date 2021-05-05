package main

import (
	"bufio"
	"fmt"
	//"github.com/ulikunitz/xz"
	"github.com/xi2/xz"
	"os"
	"strings"
	"sync"
)

func read_jm_log(filename string, jobs []string, alllogs *allLogs, wg *sync.WaitGroup, archive bool) {
	var reader *bufio.Reader
	loglines := make([]string, 0, 10000)

	file, err := os.Open(filename)
	if err == nil {
		fmt.Println("reading", filename, "...")
		if archive {
			breader := bufio.NewReaderSize(file, 16*1024*1024)
			xzreader, err := xz.NewReader(breader, 0)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader = bufio.NewReaderSize(xzreader, 16*1024*1024)
		} else {
			reader = bufio.NewReaderSize(file, 16*1024*1024)
		}
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
	if len(loglines) > 0 {
		alllogs.Append(loglines, filename)
	}
	wg.Done()
}
