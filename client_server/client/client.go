package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestData struct {
	Message string `json:"message"`
	Number  int    `json:"number"`
}

type ResponseData struct {
	Message   string `json:"message"`
	Modified  bool   `json:"modified"`
	NewNumber int    `json:"newNumber"`
}

func main() {
	url := "http://localhost:8080/"

	// 准备请求的数据
	data := RequestData{
		Message: "Hello from client!",
		Number:  42,
	}

	// 将请求数据转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
		return
	}

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("POST request failed:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应的 Body
	body, _ := ioutil.ReadAll(resp.Body)

	// 解析响应的 JSON 数据
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("JSON unmarshaling failed:", err)
		return
	}

	// 打印响应数据
	fmt.Println("Response Message:", responseData.Message)
	fmt.Println("Response Modified:", responseData.Modified)
	fmt.Println("Response NewNumber:", responseData.NewNumber)
}
