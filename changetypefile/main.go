package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	fileUrl := "https://uranos.acer.com/PreloadPN/APP/Acer-HQ1/A400CP/RV02RC/APP010DPZZ000C21_20101202094524.txt"
	fileName := "file.ini"

	// 建立 http request
	resp, err := http.Get(fileUrl)
	if err != nil {
		fmt.Println("Error while downloading", fileUrl, "-", err)
		return
	}
	defer resp.Body.Close()

	// 創建本地檔案並設置權限
	out, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer out.Close()
	out.Chmod(0644)

	// 將下載內容寫入本地檔案
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error while downloading", fileUrl, "-", err)
		return
	}

	// 將 .txt 檔案更改為 .ini 檔案
	if strings.HasSuffix(fileName, ".php") {
		newFileName := strings.TrimSuffix(fileName, ".php") + ".txt"
		err = os.Rename(fileName, newFileName)
		if err != nil {
			fmt.Println("Error while renaming file", fileName, "-", err)
			return
		}
		fileName = newFileName
	}

	fmt.Println("Downloaded file", fileName)
}
