package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	url := "https://uranos.acer.com/RAGLLM.pdf" // 这里替换成你 XAMPP 网站中 PDF 文件的 URL

	// 获取当前用户信息
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("获取当前用户信息失败：%s\n", err)
		return
	}

	// 构建下载文件夹路径
	downloadFolder := filepath.Join(currentUser.HomeDir, "Downloads")

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("HTTP GET 请求失败：%s\n", err)
		return
	}
	defer response.Body.Close()

	// 创建一个文件来保存下载的 PDF 文件
	out, err := os.Create(filepath.Join(downloadFolder, "RAGLLMPaper.pdf"))
	if err != nil {
		fmt.Printf("创建文件失败：%s\n", err)
		return
	}
	defer out.Close()

	// 将 HTTP 响应体中的内容保存到 PDF 文件中
	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Printf("保存 PDF 文件失败：%s\n", err)
		return
	}

	fmt.Println("PDF 文件已成功下载到下载文件夹")
}
