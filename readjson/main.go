package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"readjson/interfaceex"
	"readjson/interfaceex2"
	"readjson/usersfunc"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func processPerson(p Person) {
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}

func main() {
	// Read JSON file
	data, err := ioutil.ReadFile("people.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal JSON data
	var people []Person
	err = json.Unmarshal(data, &people)
	if err != nil {
		log.Fatal(err)
	}

	// Process each person
	for _, p := range people {
		processPerson(p)
	}

	usersfunc.UsersTest()
	//writejson.WritejsonFunc()
	interfaceex.Interfacefunc()
	interfaceex2.InterfaceFunc2()
}
