package main

import (
	"bufio"
	"fmt"
	//"github.com/ulikunitz/xz"
	"github.com/xi2/xz"
	"os"
	"regexp"
	"strings"
	"sync"
)

// read a nqs logfile and attach lines to <alllogs>
func read_nqs_log(filename string, jobs []string, alllogs *allLogs, wg *sync.WaitGroup, archive bool) {
	var reader *bufio.Reader
	loglines := make([]string, 0, 10000)
	rids := make(map[string]bool)
	jsvs := make(map[string]bool)
	r_rids := regexp.MustCompile(`.*gmr_attach_rcb: (.*): .*request. \(rid (.*) qid`)
	r_jsvs1 := regexp.MustCompile(`.*jcb_alloc: (.*),.*,(.*): JCB was attached.`)
	r_jsvs2 := regexp.MustCompile(`.*jcb_free: (.*),.*,(.*): JCB was detatched.`)

	file, err := os.Open(filename)
	if err == nil {
		fmt.Println("reading", filename, "...")
		if archive {
			breader := bufio.NewReaderSize(file, 16*1024*1024)
			xzreader, _ := xz.NewReader(breader, 0)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader = bufio.NewReaderSize(xzreader, 16*1024*1024)
		} else {
			reader = bufio.NewReaderSize(file, 16*1024*1024)
		}
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
						m := r_rids.FindStringSubmatch(string(line))
						if m != nil {
							rids[m[1]] = true
						}
					}
				}
			} // job
			if !match {
				// check if line might contain a r..... match
				for rid, _ := range rids {
					if strings.Contains(line, rid) {
						match = true
						loglines = append(loglines, string(line))
						//fmt.Print("-append ", line)
					}
				}
				if opts.JSV {
					if strings.Contains(line, "jcb_alloc") {
						m := r_jsvs1.FindStringSubmatch(string(line))
						if m != nil {
							_, ok := rids[m[2]]
							if ok {
								jsvs[m[1]] = true
								if opts.Verbose {
									fmt.Println("using jsv", m[1], "host:", JSVmap[m[1][1:]])
								}
								go Clientlogs.Fetch(JSVmap[m[1][1:]])
								loglines = append(loglines, string(line))
							}
						}
					} else if strings.Contains(line, "jcb_free") {
						m := r_jsvs2.FindStringSubmatch(string(line))
						if m != nil {
							_, ok := rids[m[2]]
							if ok {
								delete(jsvs, m[1])
								loglines = append(loglines, string(line))
							}
						}
					}
					for jsv, _ := range jsvs {
						if strings.Contains(line, jsv+":") || strings.Contains(line, jsv+",") {
							match = true
							loglines = append(loglines, string(line))
							//fmt.Print("-append ", line)
						}
					}
				} // JSV
			}
		}
		file.Close()
	}
	//fmt.Println(loglines)
	if len(loglines) > 0 {
		alllogs.Append(loglines, filename)
	}
	wg.Done()
}
