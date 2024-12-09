// can be run using go run client.go
// go run client.go calculator.go
package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}

	args := Args{A: 10, B: 5}
	var reply int
	err = client.Call("Calculator.Add", args, &reply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result of Add: %d\n", reply)

	var divReply float64
	err = client.Call("Calculator.Divide", args, &divReply)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result of Divide: %.2f\n", divReply)
}
