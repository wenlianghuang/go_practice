package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// The connected folder is "httpPOSTUpload"
func main() {
	// 開啟檔案控制代碼
	file, err := os.Open("./example.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 建立 multipart/form-data 的請求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	//POST file name "fileName, like the <input name="fileName" /> from html"
	part, err := writer.CreateFormFile("fileName", "example.json")
	if err != nil {
		panic(err)
	}
	io.Copy(part, file)
	writer.Close()

	// 建立 POST 請求
	// 20230428 This is for local "httpPOSTMain"
	//req, err := http.NewRequest("POST", "http://localhost:9090/upload", body)
	// 20230428 This is for remote "Uranos"
	req, err := http.NewRequest("POST", "https://uranos.acer.com/Uranos/testpy/upload.php", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 讀取回應
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}
