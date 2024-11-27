package main

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	// 啟動瀏覽器
	_ := context.Background()
	u := launcher.New().MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	// 創建新的頁籤
	page := browser.MustPage()

	// 前往目標網頁
	page.MustNavigate("https://example.com")

	// 等待網頁完全加載
	page.MustWaitLoad()

	// 使用 CSS 選擇器選取元素並取得其文字內容
	element := page.MustElement("h1")
	text := element.MustText()

	fmt.Println("網頁標題:", text)

	// 關閉瀏覽器
	browser.MustClose()
}
