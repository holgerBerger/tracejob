package main

import (
	"sync"
)

type allLogs struct {
	logs  []string
	mutex sync.RWMutex
}

func NewAllogs() allLogs {
	var al allLogs
	al.logs = make([]string, 0, 1000)
	return al
}

func (al *allLogs) Append(lines []string) {
	al.mutex.Lock()
	al.logs = append(al.logs, lines...)
	al.mutex.Unlock()
}
