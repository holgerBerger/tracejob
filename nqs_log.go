package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

func read_nqs_log(filename string, jobs []string, alllogs *allLogs, wg *sync.WaitGroup) {
	loglines := make([]string, 0, 10000)
	rids := make([]string, 0, 10)
	r := regexp.MustCompile(`.*gmr_attach_rcb: (.*): .*request. \(rid (.*) qid`)

	file, err := os.Open(filename)
	if err == nil {
		fmt.Println("reading", filename, "...")
		reader := bufio.NewReaderSize(file, 16*1024*1024)
		for {
			match := false
			line, err := reader.ReadString('\n')
			// line, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			for _, job := range jobs {
				pos := strings.Index(line, job)
				// filter out some substring matches
				if pos > 0 && strings.ContainsRune(" :(", rune(line[pos-1])) {
					match = true
					loglines = append(loglines, string(line))
					//fmt.Print("append ", line)
					// find r... numer of jobs and add to list
					if strings.Contains(line, "gmr_attach_rcb") {
						m := r.FindStringSubmatch(string(line))
						if m != nil {
							rids = append(rids, m[1])
						}
					}
				}
				if !match {
					// check if line might contain a r..... match
					for _, rid := range rids {
						if strings.Contains(line, rid) {
							match = true
							loglines = append(loglines, string(line))
							//fmt.Print("-append ", line)
						}
					}
				}
			} // job
		}
		file.Close()
	}
	//fmt.Println(loglines)
	alllogs.Append(loglines)
	wg.Done()
}