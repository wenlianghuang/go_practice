/*
Goroutine 推資料入 Channel 時的等待情境大略如下
calculate 會先執行並且計算完成
calculate 將 FINISH 訊號推入 Channel
但由於目前 main 還未拉取 Channel 中的資料，所以 calculate 會被迫等待，因此 calculate 的最後一行 fmt.Println("main goroutine finished") 沒有馬上輸出在畫面上
main 拉取了 Channel 中的資料
calculate 執行fmt.Println("main goroutine finished") 並結束
main 執行完成
*/
package goroutinefolder

import (
	"fmt"
	"time"
)

func Channelblockpush() {
	ch := make(chan string)

	go func() {
		fmt.Println("calculate goroutine starts calculating")
		time.Sleep(time.Second) // Heavy calculation
		fmt.Println("calculate goroutine ends calculating")

		ch <- "FINISH" // goroutine 執行會在此被迫等待

		fmt.Println("calculate goroutine finished")
	}()

	time.Sleep(2 * time.Second) // 使 main 比 goroutine 慢
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println("main goroutine finished")
}
