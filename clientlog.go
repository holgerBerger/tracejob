package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

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

	file, err := os.Create("/tmp/tracejob_temp/" + client)
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
