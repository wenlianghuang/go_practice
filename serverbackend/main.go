package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 設定服務器的 IP 和端口
	ip := "192.168.100.9" // 這裡替換為你的服務器 IP
	//ip := "142.251.42.238"
	port := "80" // 這裡替換為你的服務器端口

	// 設定路由和處理函數
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// 開始服務器
	addr := ip + ":" + port
	fmt.Printf("服務器運行於 %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("無法啟動服務器:", err)
	}
}
