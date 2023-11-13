package click

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

func Click() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		//chromedp.Value(`#example-After textarea`, &example),
		chromedp.Text(".Documentation-function > p", &example),
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Go's time.After example:\n%s", example)
}
