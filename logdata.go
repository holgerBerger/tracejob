package main

import (
	"sync"
)

type allLogs struct {
	logs  []string
	files map[string]bool
	mutex sync.RWMutex
}

func NewAllogs() allLogs {
	var al allLogs
	al.logs = make([]string, 0, 1000)
	al.files = make(map[string]bool)
	return al
}

func (al *allLogs) Append(lines []string, filename string) {
	al.mutex.Lock()
	al.logs = append(al.logs, lines...)
	al.files[filename] = true
	al.mutex.Unlock()
}
