package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	str := `{"name":"Matt","age":30}`

	var obj interface{}
	if err := json.Unmarshal([]byte(str), &obj); err != nil {
		fmt.Println(err)
		return
	}

	age := obj.(map[string]interface{})["age"].(float64)
	fmt.Println(age)
}
