package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// where all the file located?
const BASEPATH string = "/var/opt/nec/nqsv/"

// command line options
var opts struct {
	Days     int      `long:"days" short:"n" default:"1" description:"number of days to search into the past"`
	Before   int      `long:"before" short:"b" default:"0" description:"limit search of job end, job must be between -n and -b"`
	Archive  bool     `long:"archive" short:"a" description:"access archived logfiles as well (slower)"`
	NoServer bool     `long:"noserver" short:"s"  description:"do not read batch server logs"`
	NoJM     bool     `long:"nojm" short:"j"  description:"do not read job manipulator logs"`
	JSV      bool     `long:"jsv" short:"J" description:"also show JSV information (verbose)"`
	Clients  bool     `long:"client" short:"C" description:"fetch logs from involved NQSV clients"`
	Light    bool     `long:"light" short:"l"  description:"colorize output for light terminals"`
	Dark     bool     `long:"dark" short:"d"  description:"colorize output for dark terminals"`
	NoColor  bool     `long:"nocolor" short:"c" description:"do not colorize output (default)"`
	Filter   []string `long:"filter" short:"f" description:"filter out lines containing this word"`
	Filelist bool     `long:"filelist" short:"F" description:"show list of logfiles with matches"`
	Grep     []string `long:"grep" short:"g" description:"show only lines matching regexp"`
	Verbose  bool     `long:"verbose" short:"v" description:"be more verbose"`
}

var (
	JSVmap     map[string]string
	Clientlogs ClientLog
)

func main() {

	// args, err := flags.Parse(&opts)
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] requestid [...]"
	args, err := parser.Parse()
	if err != nil {
		//fmt.Println(err.Error())
		os.Exit(1)
	}

	fileindex := make(map[string]fileIndex)

	// do not read when not wanted server or jm logs
	if !opts.NoServer {
		indexFiles(&fileindex, []string{BASEPATH + "batch_server_log*"})
		if opts.Archive {
			indexArchiveFiles(&fileindex, []string{BASEPATH + "logarchive/" + "batch_server_log*.xz"})
		}
	}
	if !opts.NoJM {
		indexFiles(&fileindex, []string{BASEPATH + "nqs_jmd*"})
		if opts.Archive {
			indexArchiveFiles(&fileindex, []string{BASEPATH + "logarchive/" + "nqs_jmd*.xz"})
		}
	}

	// normalize job ids
	for i, _ := range args {
		if strings.Contains(args[i], ".") {
			args[i] = strings.Split(args[i], ".")[0] + "."
		} else {
			args[i] = args[i] + "."
		}
	}

	/*
		// debug
		for i := range fileindex {
			fmt.Printf("%s: %s - %s\n", i, fileindex[i].first.String(), fileindex[i].last.String())
		}
	*/

	// file lists
	nqsfiles, _ := filepath.Glob(BASEPATH + "batch_server_log*")
	jmfiles, _ := filepath.Glob(BASEPATH + "nqs_jmd*")

	// archive files
	anqsfiles := make([]string, 0)
	ajmfiles := make([]string, 0)
	if opts.Archive {
		anqsfiles, _ = filepath.Glob(BASEPATH + "logarchive/" + "batch_server_log*.xz")
		ajmfiles, _ = filepath.Glob(BASEPATH + "logarchive/" + "nqs_jmd*.xz")
	}

	// find which files to read, depending on days parameter
	filefilter := make(map[string]bool)
	now := time.Now()
	for fn, ft := range fileindex {
		if ft.first.AddDate(0, 0, opts.Days).After(now) {
			filefilter[fn] = true
			if ft.last.AddDate(0, 0, opts.Before).After(now) {
				filefilter[fn] = false
			}
		} else {
			filefilter[fn] = false
		}
		if opts.Verbose {
			fmt.Printf(" filefilter[%s]:%v\n", fn, filefilter[fn])
		}
	}

	// compile regexp early
	compiled_regexp := make([]*regexp.Regexp, 0, len(opts.Grep))
	if len(opts.Grep) > 0 {
		for _, reg := range opts.Grep {
			r, err := regexp.Compile(reg)
			if err != nil {
				fmt.Printf("could not compile regexp <%s>: %s\n", reg, err.Error())
			} else {
				compiled_regexp = append(compiled_regexp, r)
			}
		}
	}

	// get JSV mapping
	if opts.Clients {
		JSVmap = jsvmap()
		Clientlogs = NewClientlog()
	}

	// read log files
	alllogs := NewAllogs()

	wg := &sync.WaitGroup{}
	if !opts.NoServer {
		for _, file := range nqsfiles {
			if filefilter[file] {
				wg.Add(1)
				go read_nqs_log(file, args, &alllogs, wg, false)
			}
		}
		if opts.Archive {
			for _, file := range anqsfiles {
				if filefilter[file] {
					wg.Add(1)
					go read_nqs_log(file, args, &alllogs, wg, true)
				}
			}
		}
	}
	if !opts.NoJM {
		for _, file := range jmfiles {
			if filefilter[file] {
				wg.Add(1)
				go read_jm_log(file, args, &alllogs, wg, false)
			}
		}
		if opts.Archive {
			for _, file := range ajmfiles {
				if filefilter[file] {
					wg.Add(1)
					go read_jm_log(file, args, &alllogs, wg, true)
				}
			}
		}
	}
	wg.Wait()
	// ... done with logfiles

	// wait for all client logs
	if opts.Clients {
		Clientlogs.Wait()
	}

	// sort and print
	if !opts.NoColor && (opts.Dark || opts.Light) {
		initcolors()
	}
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
		grepped := true
		if len(opts.Grep) > 0 {
			grepped = false
			for _, regexp := range compiled_regexp {
				match := regexp.MatchString(l)
				if match {
					grepped = true
					break
				}
			}
		}
		if !filtered {
			if grepped {
				if !opts.NoColor && (opts.Dark || opts.Light) {
					colorize(&l)
				}
				fmt.Printf(l)
			}
		}
	}
	fmt.Printf("\033[0m")

	if opts.Filelist {
		fmt.Println("\nMinimal list of logfiles with matches:")
		for f, _ := range alllogs.files {
			fmt.Println(f)
		}
	}

}
