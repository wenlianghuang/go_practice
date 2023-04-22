package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 設定下載檔案的 URL 和 POST 資料
	url := "http://localhost:8080/download"
	postData := "key=value"

	// 發送 POST 請求
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 讀取檔案內容
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得檔案的 MIME 類型
	contentType := http.DetectContentType(content)

	// 寫入檔案
	filePath := "example.pdf"
	err = ioutil.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 設定檔案下載時的 MIME 類型
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))

	// 使用 http.ServeFile 函數提供檔案內容
	http.ServeFile(w, r, filePath)
}

func main() {
	http.HandleFunc("/download", downloadHandler)
	http.ListenAndServe(":8080", nil)
}
