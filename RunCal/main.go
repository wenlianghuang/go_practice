package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 在Windows上执行打开计算器的命令
	cmd := exec.Command("calc")

	// 执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
