package main

//create a struct called person with fields name and age
//read a list of names from a text file
//for each name on the text file, create a person object with a given name and random age[0-100]
//add each created person into a slice of persons
//marshall this slice into json bytes and write to a json file

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
)

type person struct {
	Name string
	Age  int
}

func main() {

	data, err := ioutil.ReadFile("names.txt")
	if err != nil {
		log.Fatalln(err)
	}

	names := strings.Split(string(data), "\n")

	people := make([]person, 0)

	for i := range names {
		//fmt.Println((names[i]))
		age := rand.Intn(100) + 0
		object := person{names[i], age}
		fmt.Printf("%v", object)
		//fmt.Println(object.Age)
		people = append(people, object)
	}
	for i := range people {
		fmt.Println(people[i].Name)
		fmt.Println(people[i].Age)
	}
}
