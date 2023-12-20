package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"readjson/interfaceex"
	"readjson/interfaceex2"
	"readjson/usersfunc"
	"readjson/writejson"
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
	//data, err := ioutil.ReadFile("people.json")
	if _, err := os.Stat("write.json"); err == nil {
		os.Remove("./write.json")
	} else {
		fmt.Println("File 'write.json does not exist\n")
	}
	data, err := os.ReadFile("people.json")
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
	writejson.WritejsonFunc()
}
