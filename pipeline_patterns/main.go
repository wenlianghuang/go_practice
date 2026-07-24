package main

import (
	"fmt"
)

func main() {
	fmt.Println("--- 使用 done channel ---")
	runWithDone()

	fmt.Println("--- 使用 context ---")
	runWithContext()

	fmt.Println("主程式結束")
}
