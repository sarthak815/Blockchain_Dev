package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	"github.com/near/borsh-go"
	"gopkg.in/yaml.v3"
)

type Person struct {
	Name      string
	Age       uint
	Ethnicity string
	Gender    string
	DOB       string
}

func main() {
	p := Person{"Manish", 23, "Indian", "Male", "1999-05-01"}
	// JSON
	jsonwire, err := json.Marshal(p)
	if err != nil {
		log.Fatalln("Error while JSON Marshalling", err)
	}
	// YAML
	yamlwire, err := yaml.Marshal(p)
	if err != nil {
		log.Fatalln("Error while YAML Marshalling", err)
	}
	// Borsh <- Near Protocol
	borshwire, err := borsh.Serialize(p)
	if err != nil {
		log.Fatalln("Error while Borsh Marshalling", err)
	}
	// Gob -> Go Native Encoding
	// Encoding
	buf := bytes.NewBuffer(make([]byte, 0))
	encoder := gob.NewEncoder(buf)
	if err := encoder.Encode(p); err != nil {
		log.Fatalln("Error while Gob Marshalling", err)
	}
	gobwire := buf.Bytes()
	// Decoding
	reader := bytes.NewReader(gobwire)
	decoder := gob.NewDecoder(reader)
	newperson := new(Person)
	if err := decoder.Decode(newperson); err != nil {
		log.Fatalln("Error while Gob Unmarshalling", err)
	}
	fmt.Println(newperson)
	fmt.Println("JSON:", jsonwire)
	fmt.Println("YAML:", yamlwire)
	fmt.Println("Borsh:", borshwire)
	fmt.Println("Gob:", gobwire)
}
