/*
Goroutine 拉資料出 Channel 時的等待情境
main 因拉取的時候 calculate 還沒將資料推入 Channel 中，因此 main 會被迫等待，因此 main 的最後一行 fmt.println 沒有馬上輸出在畫面上
calculate 執行並且計算完成
calculate 將 FINISH 推入 Channel
calculate 執行完成
main 拉取了 Channel 中的資料並且執行完成
*/

package goroutinefolder

import (
	"fmt"
	"time"
)

func Channelblockpull() {
	ch := make(chan string)

	go func() {
		fmt.Println("calculate goroutine starts calculating")
		time.Sleep(time.Second) // Heavy calculation
		fmt.Println("calculate goroutine ends calculating")

		ch <- "FINISH"

		fmt.Println("calculate goroutine finished")
	}()

	fmt.Println("main goroutine is waiting for channel to receive value")
	fmt.Println(<-ch) // goroutine 執行會在此被迫等待
	fmt.Println("main goroutine finished")
}
