package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

const BASEPATH string = "/home/hobel/big/var/opt/nec/nqsv/"

var opts struct {
	Days     int      `long:"days" short:"n" description:"number of days to search"`
	NoServer bool     `long:"noserver" short:"s" default:"false" description:"do not read batch server logs"`
	NoJM     bool     `long:"nojm" short:"j" default:"false" description:"do not read job manipulator logs"`
	NoColor  bool     `long:"nocolor" short:"c" description:"do not colorize output"`
	Filter   []string `long:"filter" short:"f" description:"filter out lines containing this word"`
}

func main() {

	args, err := flags.Parse(&opts)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fileindex := make(map[string]fileIndex)

	if !opts.NoServer {
		indexFiles(&fileindex, []string{BASEPATH + "batch_server_log*"})
	}
	if !opts.NoJM {
		indexFiles(&fileindex, []string{BASEPATH + "nqs_jmd*"})
	}

	// debug
	for i := range fileindex {
		fmt.Printf("%s: %s - %s\n", i, fileindex[i].first.String(), fileindex[i].last.String())
	}

	// file lists
	nqsfiles, _ := filepath.Glob(BASEPATH + "batch_server_log*")
	jmfiles, _ := filepath.Glob(BASEPATH + "nqs_jmd*")

	alllogs := NewAllogs()

	// FIXME add filter logic to read only necessary files

	wg := &sync.WaitGroup{}
	if !opts.NoServer {
		for _, file := range nqsfiles {
			wg.Add(1)
			go read_nqs_log(file, args, &alllogs, wg)
		}
	}
	if !opts.NoJM {
		for _, file := range jmfiles {
			wg.Add(1)
			go read_jm_log(file, args, &alllogs, wg)
		}
	}
	wg.Wait()

	sort.SliceStable(alllogs.logs, func(i, j int) bool {
		return alllogs.logs[i][:15] < alllogs.logs[j][:15]
	})
	for _, l := range alllogs.logs {
		filtered := false
		if len(opts.Filter) > 0 {
			for _, fword := range opts.Filter {
				if strings.Contains(l, fword) {
					filtered = true
					break
				}
			}
		}
		if !filtered {
			fmt.Printf(l)
		}
	}

}
