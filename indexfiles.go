package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type fileIndex struct {
	first, last time.Time
}

// read first and lastline of file and parse date
func indexFiles(fi *map[string]fileIndex, globlist []string) {
	year := strconv.Itoa(time.Now().Year())
	fmt.Println("year:", year)
	for _, glob := range globlist {
		fmt.Println(glob)
		filelist, err := filepath.Glob(glob)
		if err == nil {
			for _, file := range filelist {
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
						break
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
						break
					}
					(*fi)[file] = fileIndex{firsttime, lasttime}
				}
			}
		}
	}
}
