package goroutinefolder

import "fmt"

func ping(pingCh chan<- string, msg string) {
	pingCh <- msg
}

func pong(pingCh <-chan string, pongCh chan<- string) {
	msg := <-pingCh
	pongCh <- msg
}

// Directions demonstrates the use of buffered channels to pass messages between
// two functions, `ping` and `pong`. It creates two channels, sends a message through
// the `ping` function, passes it to the `pong` function, and prints the final result
// received from the `pongs` channel.
func Directions() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
