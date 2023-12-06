// map[string]interface: string=> 鍵 interface => 值
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//var m map[string]interface{}
	m := make(map[string]interface{})
	m["name"] = "say"
	age := 18
	m["age"] = age
	m["addr"] = "China"
	print_map(m)
	fmt.Println()

	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println("err:", err)

	} else {
		fmt.Println("data:", data)
	}

	m1 := make(map[string]interface{})
	err = json.Unmarshal(data, &m1)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		print_map(m1)
	}

	fmt.Println()
	value, ok := m1["key1"]
	if ok {
		fmt.Println(value.(string))
	} else {
		fmt.Println("key1 不存在")
	}
}

func print_map(m map[string]interface{}) {
	fmt.Println("enter print_map##########")
	for k, v := range m {
		switch value := v.(type) {
		case nil:
			fmt.Println(k, "is nil", "null")
		case string:
			fmt.Println(k, "is string", value)
		case int:
			fmt.Println(k, "is int", value)
		case float64:
			fmt.Println(k, "is float64", value)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range value {
				fmt.Println(i, u)
			}
		case map[string]interface{}:
			fmt.Println(k, "is an map:")
			print_map(value)
		default:
			fmt.Println(k, "is unknown type", fmt.Sprintf("%T", v))
		}
	}
	fmt.Println("out print_map##########")
}
