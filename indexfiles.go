package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type fileIndex struct {
	first, last time.Time
}

// read first and lastline of files and parse date in parallel
func indexFiles(fi *map[string]fileIndex, globlist []string) {
	year := strconv.Itoa(time.Now().Year())
	wg := sync.WaitGroup{}
	wg.Add(1)
	m := sync.RWMutex{}
	for _, glob := range globlist {
		filelist, err := filepath.Glob(glob)
		if err == nil {
			for _, file := range filelist {
				wg.Add(1)
				go func(file string) {
					defer wg.Done()
					if opts.Verbose {
						fmt.Println(" indexing file", file)
					}
					f, err := os.Open(file)
					if err == nil {
						r := bufio.NewReader(f)
						defer f.Close()
						line, _ := r.ReadBytes('\n')
						firstline := string(line)
						// fmt.Println("firstline:", firstline)
						f.Seek(-1024, 2)
						lastline := ""
						for {
							line, err := r.ReadBytes('\n')
							if err != nil {
								break
							}
							lastline = string(line)
						}
						//fmt.Println("lastline:", lastline)

						var firsttime, lasttime time.Time

						if firstline != "" {
							// 09/02 12:27:30
							parsed, err := time.Parse("2006 01/02 15:04:05", year+" "+firstline[:14])
							if err == nil {
								firsttime = parsed
								//	fmt.Println(firsttime.String())
							} else {
								fmt.Println("error in Parsing fistline of ", file, ":", err.Error())
							}
						} else {
							return
						}
						if lastline != "" {
							// 09/02 12:27:30
							parsed, err := time.Parse("2006 01/02 15:04:05", year+" "+lastline[:14])
							if err == nil {
								lasttime = parsed
								//	fmt.Println(lasttime.String())
							} else {
								fmt.Println("error in Parsing fistline of ", file, ":", err.Error())
							}
						} else {
							return
						}
						m.Lock()
						(*fi)[file] = fileIndex{firsttime, lasttime}
						m.Unlock()
					}
				}(file)
			} // files
		}
	}
	wg.Done()
	wg.Wait()
}