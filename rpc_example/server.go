// can be run using go run server.go
// go run server.go calculator.go
package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	calculator := new(Calculator)

	err := rpc.Register(calculator)
	if err != nil {
		log.Fatal(
			"Error registering calculator",
		)
	}

	listener, err := net.Listen("tcp", ":8222")
	if err != nil {
		log.Fatal(
			"Error creating listener",
		)
	}
	log.Println("RPC server is running on port 8222")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
