// Package writer provides basic functions to write in success.log and error.log files
package main

import (
	"fmt"
	"github.com/LepikovStan/backlinkCrawler/lib"
	"io"
	"os"
	"sync"
)

type ILogger interface {
	Log(args ...string)
}

type Logger struct {
	mu    sync.Mutex
	files map[string]*os.File
}

func wlog(f io.Writer, args ...string) {
	for i := 0; i < len(args); i++ {
		fmt.Fprintln(f, args[i])
	}
}

func (l *Logger) Init() {
	l.files = map[string]*os.File{
		"result": lib.CreateLog("result.log"),
		"error":  lib.CreateLog("error.log"),
	}
}

func (l *Logger) Log(level string, args ...string) {
	l.mu.Lock()
	wlog(l.files[level], args...)
	l.mu.Unlock()
}
