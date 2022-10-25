package goroutinefolder

import "fmt"

func ping(ping chan<- string, msg string) {
	ping <- msg
}

func pong(ping <-chan string, pong chan<- string) {
	msg := <-ping
	pong <- msg
}
func Directions() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
