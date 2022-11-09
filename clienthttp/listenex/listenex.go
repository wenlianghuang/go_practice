package listenex

import (
	"fmt"
	"net"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
	CONN_TYPE = "tcp"
)

func Listenserver() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		panic(err)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			panic(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	fmt.Println("Request Content Length: ", reqLen)

	conn.Write(([]byte("Message received.")))
	conn.Close()
}
