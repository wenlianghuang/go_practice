package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	sublogrus "github.com/username/myproject/golog/Sublogrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

}
func main() {

	LogFilePath, errLogFilePath := os.Getwd()
	if errLogFilePath != nil {
		panic(errLogFilePath)
	}
	logf := sublogrus.Sublogrusfunc(LogFilePath)

	// 創建一個 JSON 格式的資料
	data := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	// 將資料轉換為 JSON 格式
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		logf.Error("無法將資料轉換為 JSON 格式: %v", err)
		log.Fatalf("無法將資料轉換為 JSON 格式: %v", err)
	}

	// 寫入 JSON 資料到檔案
	err = os.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		logf.Error("無法寫入 JSON 資料到檔案: %v", err)
		log.Fatalf("無法寫入 JSON 資料到檔案: %v", err)
	}

	//fmt.Println("成功創建並儲存 JSON 檔案")
	logf.Info("成功創建並儲存 JSON 檔案")
	// 讀取檔案並解析 JSON 資料
	fileContent, err := os.ReadFile("data.json")

	if err != nil {
		logf.Error("無法讀取檔案: %v", err)
		log.Fatalf("無法讀取檔案: %v", err)
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal(fileContent, &parsedData)
	if err != nil {
		logf.Error("無法解析 JSON 資料: %v", err)
		log.Fatalf("無法解析 JSON 資料: %v", err)
	}

	// 修改其中一個值
	parsedData["age"] = 35

	// 再次將資料轉換為 JSON 格式
	updatedJSONData, err := json.MarshalIndent(parsedData, "", "    ")
	if err != nil {
		logf.Error("無法將資料轉換為 JSON 格式: %v", err)
		log.Fatalf("無法將資料轉換為 JSON 格式: %v", err)
	}

	// 寫入更新後的 JSON 資料到檔案
	err = os.WriteFile("data.json", updatedJSONData, 0644)
	if err != nil {
		logf.Error("無法寫入更新後的 JSON 資料到檔案: %v", err)
		log.Fatalf("無法寫入更新後的 JSON 資料到檔案: %v", err)
	}

	//fmt.Println("成功更新 JSON 檔案的值")
	name := "Matt"
	logf.Info("成功更新" + name + "JSON 檔案的值")
}
