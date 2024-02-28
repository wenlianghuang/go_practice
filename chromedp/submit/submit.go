package submit

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

func Submit() {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()

	var res string
	err := chromedp.Run(ctx, submit(`https://github.com/search`, `//input[@name="q"]`, `chromedp`, &res))
	//err := chromedp.Run(ctx, submit(`https://www.twse.com.tw/zh/trading/historical/stock-day.html`, `//input[@name="q"]`, "2049", &res))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("got: `%s`", strings.TrimSpace(res))

}

func submit(urlstr, sel, q string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel),
		chromedp.SendKeys(sel, q),
		chromedp.Submit(sel),
		chromedp.WaitVisible(`//*[contains(., 'repository results')]`),
		chromedp.Text(`(//*//ul[contains(@class, "repo-list")]/li[1]//p)[1]`, res),
	}
}
