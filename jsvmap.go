package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

// return a map jsvid->execution host
func jsvmap() map[string]string {
	jsvmap := make(map[string]string)

	out, err := exec.Command("/opt/nec/nqsv/bin/qstat", "-Po", "-Sn", "-Fjsvno,ehost").Output()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Fields(line)
		jsvmap[s[0]] = s[1]
	}

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	return jsvmap
}
