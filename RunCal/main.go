package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("calc")
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}
