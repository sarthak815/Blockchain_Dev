package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func main() {
	//read from a file data.txt
	//contains a line separated (\n) list of names
	//read list of names into a slice of strings
	//sort it and overwrite the existing file

	data, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatalln(err)
	}

	strdata := string(data)
	names := strings.Split(strdata, "\n")

	sort.Strings(names)
	fmt.Println("Sorted :", names)

	newNames := strings.Join(names, "\n")
	err = ioutil.WriteFile("data.txt", []byte(newNames), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
