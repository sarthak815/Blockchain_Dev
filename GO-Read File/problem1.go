package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	//Read list of names into slice of strings
	//sort and overwrite existing file
	data, err := os.OpenFile("data.txt", os.O_RDWR, 0777)
	strArray := make([]string, 0)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		strArray = append(strArray, scanner.Text())
	}

	sort.Strings(strArray)
	newnames := strings.Join(strArray, "\n")
	data.Truncate(0)
	_, _ = data.WriteString(newnames)
}
