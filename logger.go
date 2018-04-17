// Package writer provides basic functions to write in success.log and error.log files
package main

import (
	"fmt"
	"github.com/LepikovStan/backlinkCrawler/lib"
	"io"
	"os"
	"sync"
)

func fmtResultLog(msg *Backlink) []string {
	result := make([]string, len(msg.BLList)+1)
	result[0] = msg.Url
	for i := 0; i < len(msg.BLList); i++ {
		result[i+1] = fmt.Sprintf("    %s", msg.BLList[i].Url)
	}
	return result
}

func fmtErrorLog(msg *Backlink) []string {
	return []string{
		msg.Url,
		fmt.Sprintf("    %s", msg.Error),
	}
}

type ILogger interface {
	Log(args ...string)
}

type Logger struct {
	mu    sync.Mutex
	files map[string]*os.File
}

func wlog(f io.Writer, args ...string) {

	defer func() {
		fmt.Println("Log Done")
	}()
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

func (l *Logger) Log(level string, msg *Backlink) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Println("Log")
	args := []string{}
	switch level {
	case "result":
		args = fmtResultLog(msg)
	case "error":
		args = fmtErrorLog(msg)
	}
	wlog(l.files[level], args...)
}
