package lib

import (
	"bufio"
	"fmt"
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
