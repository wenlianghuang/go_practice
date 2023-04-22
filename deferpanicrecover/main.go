package main

import "fmt"

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from:", err)
		}
	}()

	fmt.Println("Starting the program...")
	// 調用引發 panic 的函數
	causePanic()
	fmt.Println("Program ended normally.")
}

func causePanic() {
	defer fmt.Println("deferred statements in causePanic()")

	fmt.Println("Start executing causePanic()...")

	panic("Somehing went wrong.")
	fmt.Println("End exeucting causePanic()...")
}
