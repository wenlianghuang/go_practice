package main

import (
	"fmt"
	"os/exec"
)

func main() {
<<<<<<< HEAD
	cmd := exec.Command("calc")
=======
	// 在Windows上执行打开计算器的命令
	cmd := exec.Command("calc")

	// 执行命令
>>>>>>> 43c5e91252e895e66ac8bc4cfc68cbdcc66685a9
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
