package writejson

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func WritejsonFunc() {

	//Create a slice of Person structs
	people := []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}

	//Marshal the data into a JSON byte slice
	data, err := json.MarshalIndent(people, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the data to a file named "write.json"
	//err = ioutil.WriteFile("write.json", data, 0644)
	err = os.WriteFile("write.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully wrote data to file!")
}
