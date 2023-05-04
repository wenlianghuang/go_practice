package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	data1 := map[string]interface{}{
		"name": "Matt",
		"age":  30,
	}

	data2 := map[string]interface{}{
		"name": "Jack",
		"age":  45,
	}

	//combine map[string]interface{}--1 and map[string]interface{}--2
	datas := []map[string]interface{}{data1, data2}

	jsonData, err := json.Marshal(datas)
	if err != nil {
		fmt.Println(err)
		return
	}

	url := "http://localhost:8000/receive"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
	res := string(body)
	fmt.Println(res)

	var jsonDataSingle interface{}
	// body = []byte
	err = json.Unmarshal(body, &jsonDataSingle)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonStr, err := json.MarshalIndent(jsonDataSingle, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err2 := ioutil.WriteFile("data.json", jsonStr, 0644)
	if err2 != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Save to finalRes.json successfully")

	/*
		var obj interface{}
		if err3 := json.Unmarshal([]byte(res), &obj); err != nil {
			fmt.Println(err3)
			return
		}
		age := obj.(map[string]interface{})["message"].(string)
		fmt.Println(age)
	*/
}
