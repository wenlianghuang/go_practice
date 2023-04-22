package interfaceex

import (
	"encoding/json"
	"fmt"
	"log"
)

func Interfacefunc() {
	jsonStr := `
	{
		"name": "John",
		"age": 30,
		"isMarried": true,
		"hobbies": [
			"reading",
			"traveling",
			"sports"
		],
		"address":{
			"street": "123 Main St",
			"city": "New York",
			"state": "NY",
			"zip": "10001"
		}
	}
	`

	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	//Access the data using type assertion, like {"a":"a-1"}
	m := data.(map[string]interface{})
	fmt.Println("Name: ", m["name"].(string))
	fmt.Println("Age: ", m["age"].(float64))
	fmt.Println("IsMarried:", m["isMarried"].(bool))

	//like {"b":["b-1","b-2"]}
	hobbies := m["hobbies"].([]interface{})
	fmt.Println("Hobbies:")
	for _, hobby := range hobbies {
		fmt.Println("-", hobby.(string))
	}

	//like {"c":"c-1"}
	address := m["address"].(map[string]interface{})
	fmt.Println("Address:")
	fmt.Println("Street:", address["street"].(string))
	fmt.Println("City:", address["city"].(string))
	fmt.Println("State:", address["state"].(string))
	fmt.Println("Zip:", address["zip"].(string))
}
