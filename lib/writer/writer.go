// Package writer provides basic functions to write in success.log and error.log files
package writer

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Writer is the type, contains methods for write results of pages crawling
type Writer struct {
	successmux  sync.Mutex
	errmux      sync.Mutex
	successFile *os.File
	errorFile   *os.File
}

// Init function creates success.log and error.log files to write results
func (w *Writer) Init(dir string) {
	successFileName := fmt.Sprintf("%s/%s", dir, "success.log")
	successFile, err := os.Create(successFileName)
	if err != nil {
		log.Fatal(err)
	}
	w.successFile = successFile

	errorFileName := fmt.Sprintf("%s/%s", dir, "error.log")
	errorFile, err := os.Create(errorFileName)
	if err != nil {
		log.Fatal(err)
	}
	w.errorFile = errorFile
}

// Destroy function write close success.log and error.log files
func (w *Writer) Destroy() {
	w.successFile.Close()
	w.errorFile.Close()
}

// WriteResult function write success result of crawling and parsing to success.log file
func (w *Writer) WriteResult(s string) {
	w.successmux.Lock()
	defer w.successmux.Unlock()
	if _, err := w.successFile.WriteString(s); err != nil {
		log.Fatal(err)
	}
}

// WriteError function write errors of crawling and parsing to error.log file
func (w *Writer) WriteError(s string) {
	w.errmux.Lock()
	defer w.errmux.Unlock()
	if _, err := w.errorFile.WriteString(s); err != nil {
		log.Fatal(err)
	}
}

// New function initialize new Writer instance and return pointer to it
func New(dir string) *Writer {
	wr := new(Writer)
	wr.Init(dir)
	return wr
}
