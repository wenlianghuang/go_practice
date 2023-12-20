package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func GetHttpHtmlContent(url string, selector string, sel interface{}) (string, error) {

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))

	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)

	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func GetSpecialData(htmlContent string, selector string) ([][]string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var results [][]string

	// 使用選擇器尋找所有匹配的元素
	dom.Find(selector).Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr.alt-row").Each(func(j int, trSelection *goquery.Selection) {
			var row []string
			trSelection.Find("td").Each(func(k int, tdSelection *goquery.Selection) {
				row = append(row, tdSelection.Text())
			})

			results = append(results, row)
		})
	})

	return results, nil
}

func main() {
	param := `document.querySelector("body")`
	selector := "#CPHB1_gv > tbody"
	url := "https://histock.tw/stock/rank.aspx?p=all"
	html, _ := GetHttpHtmlContent(url, selector, param)
	res, _ := GetSpecialData(html, selector)
	f, err := os.OpenFile("./data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	w := csv.NewWriter(f)
	for _, row := range res {
		if len(row) == 0 || len(row) == 1 {
			continue
		}
		w.Write(row)
		//fmt.Printf("%v\n", row)

	}

}
