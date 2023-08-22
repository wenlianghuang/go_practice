// 其實他的response主要還是body,header主要是網路的一些一定會存在的設定。
package main

import (
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 讀取request的Body
		body, _ := ioutil.ReadAll(r.Body)

		// 解析request的JSON數據
		var requestData RequestData
		err := json.Unmarshal(body, &requestData)
		if err != nil {
			fmt.Println("JSON unmarshaling failed:", err)
			return
		}

		// 修改request的數據
		requestData.Message = "Modified by server"
		requestData.Number += 10
		fmt.Println(requestData.Message)
		// 將respond數據轉為 JSON字串
		jsonData, err := json.Marshal(ResponseData{
			Message:   requestData.Message,
			Modified:  true,
			NewNumber: requestData.Number,
		})
		if err != nil {
			fmt.Println("JSON marshaling failed:", err)
			return
		}

		// 設置response的Header
		w.Header().Set("Content-Type", "application/json")

		// Response Body
		fmt.Fprint(w, string(jsonData))
	})

	//http.ListenAndServe("192.168.100.9:80", nil)
	http.ListenAndServe("192.168.100.9:80", nil)
}
