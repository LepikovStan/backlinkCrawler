package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const CHUNK = 100

func ReadFile(path string) []string {
	var result []string
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}

func ReadStream(r io.Reader) ([]byte, error) {
	var data, result []byte
	var error error
	for {
		data = make([]byte, CHUNK)
		count, err := r.Read(data)
		result = append(result, data[:count]...)
		if err != io.EOF {
			error = err
		}
		if err != nil || count == 0 {
			break
		}
	}
	return result, error
}

func CreateResultDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	resultDirName := fmt.Sprintf("%s/%s/%d", dir, "results", time.Now().Unix())

	if _, err := os.Stat(resultDirName); os.IsNotExist(err) {
		err = os.MkdirAll(resultDirName, 0755)
	}

	if err != nil {
		log.Fatal(err)
	}
	return resultDirName
}

func CreateLog(filename string) *os.File {
	filepath := fmt.Sprintf("%s/%s", CreateResultDir(), filename)
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
