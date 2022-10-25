package main

import (
	"fmt"
	"reversesimple/reversefunc"
)

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Reverse Navie: %+v\n", reversefunc.Reverse_Naive(input))
	fmt.Printf("Reverse Memory: %+v\n", reversefunc.Reverse_Memory(input))
	fmt.Printf("Reverse Incorporating Concurrency: %+v\n", reversefunc.Reverese_Incoporating_Concurrency(input))
	reversefunc.Reverse_Interger()
}
