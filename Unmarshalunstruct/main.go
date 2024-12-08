// 這邊可以比較在Unmarshal的已經定義好struct的情況下，如何取得JSON中的資料。
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// suppose this is the JSON data received from a client
	jsonData := `{
		"name": "Alice",
		"age": 25,
		"hobbies": ["gaming", "reading", "swimming"],
		"address": {
			"city": "San Francisco",
			"zipcode": 94107},
		"is_student": false}`
	// using map[string]interface{} to store the JSON data
	// interface{} 是Go中的通用類型。解析後需要將其轉換為實際的類型。
	var result map[string]interface{}

	// analyzing the result
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		fmt.Println("JSON decode failed: ", err)
		return
	}

	// dynamically accessing the JSON data
	fmt.Println("Name: ", result["name"].(string))
	fmt.Println("Age: ", result["age"].(float64))
	fmt.Println("Is Student: ", result["is_student"].(bool))
	fmt.Println("City: ", result["address"].(map[string]interface{})["city"].(string))
	fmt.Println("Zipcode: ", result["address"].(map[string]interface{})["zipcode"].(float64))

	// running a loop to access the hobbies
	fmt.Println("Hobbies:")
	for _, hobby := range result["hobbies"].([]interface{}) {
		fmt.Println(hobby.(string))
	}
}
