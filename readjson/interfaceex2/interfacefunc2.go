package interfaceex2

import (
	"encoding/json"
	"fmt"
	"log"
)

func InterfaceFunc2() {
	/*
		jsonStr := `
		{
			"name": "Matt",
			"Age": 30,
			"Education": [
				"NTU",
				"KC",
				"LM"
			]
		}
		`
	*/
	jsonStr := `
	{
		"name": "Matt",
		"Age": 30,
		"Education": [
			"NTU",
			"KC",
			"LM",
			"NTUEE"
		]
	}
	`

	var data interface{}

	err := json.Unmarshal([]byte(jsonStr), &data)

	if err != nil {
		log.Fatal(err)
	}

	m := data.(map[string]interface{})
	fmt.Println("Name: ", m["name"].(string))
	fmt.Println("Age: ", m["Age"].(float64))

	educations := m["Education"].([]interface{})
	for _, education := range educations {
		fmt.Println("-", education.(string))
	}
}
