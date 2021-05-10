package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

const clientlogpathprefix = "/tmp/tracejob_temp/"

type ClientLog struct {
	clients map[string]bool
	mutex   sync.RWMutex
	wg      sync.WaitGroup
}

func NewClientlog() ClientLog {
	var cl ClientLog
	cl.clients = make(map[string]bool)
	return cl
}

// fetch log data if not yet done
func (cl *ClientLog) Fetch(client string) {
	cl.mutex.Lock()
	_, ok := cl.clients[client]
	if ok {
		cl.mutex.Unlock()
		return
	}
	cl.clients[client] = true
	cl.mutex.Unlock()

	if opts.Verbose {
		fmt.Println("fetching logs from", client)
	}
	cl.wg.Add(1)
	out, err := exec.Command("/usr/bin/ssh", "-lroot", "-oStrictHostKeyChecking=no", client, "/usr/bin/journalctl").Output()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("/tmp/tracejob_temp", 0700)

	file, err := os.Create(clientlogpathprefix + client)
	if err != nil {
		log.Fatal(err)
	}

	file.Write(out)
	file.Close()
	cl.wg.Done()
	if opts.Verbose {
		fmt.Println("done fetching logs from", client)
	}
}

func (cl *ClientLog) Wait() {
	cl.wg.Wait()
}

// filter for date range, convert date string and filter out some noise
func (cl *ClientLog) Filter(start string, end string) []string {
	loglines := make([]string, 0, 1000)

	re := regexp.MustCompile("Vector Engine MMM-Command|journalbeat|ENCRYPT_METHOD|Accepted publickey|sshd")

	for f, _ := range cl.clients {
		cf, err := os.Open(clientlogpathprefix + f)
		if err == nil {
			defer cf.Close()

			scanner := bufio.NewScanner(cf)
			for scanner.Scan() {
				line := scanner.Text()
				line = convertdate(line)
				if line[:15] >= start && line[:15] <= end {
					if !re.MatchString(line) {
						loglines = append(loglines, line+"\n")
						// fmt.Println("+", line)
					}
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

		} // err
	} // clients

	return loglines
}

var monthmap = map[string]string{
	"Jan": "01/",
	"Feb": "02/",
	"Mar": "03/",
	"Apr": "04/",
	"May": "05/",
	"Jun": "06/",
	"Jul": "07/",
	"Aug": "08/",
	"Sep": "09/",
	"Oct": "10/",
	"Nov": "11/",
	"Dec": "12/",
}

func convertdate(line string) string {
	month := line[0:3]
	line = monthmap[month] + line[4:]
	return line
}
