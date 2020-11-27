package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
)

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	nqsfiles, _ := filepath.Glob("/home/hobel/big/var/opt/nec/nqsv/batch_server_log*")

	alllogs := NewAllogs()

	wg := &sync.WaitGroup{}
	for _, file := range nqsfiles {
		wg.Add(1)
		go read_nqs_log(file, os.Args[1:], &alllogs, wg)
	}
	wg.Wait()

	for _, l := range alllogs.logs {
		fmt.Printf(l)
	}
}
