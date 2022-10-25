package reversefunc

import (
	"fmt"
	"sync"
)

func Reverse_Naive(input []int) []int {
	var output []int

	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}

	return output
}

func Reverse_Memory(input []int) []int {
	inputLen := len(input)
	output := make([]int, inputLen)

	for i, n := range input {
		i := inputLen - i - 1
		output[i] = n

	}
	return output
}

func Reverese_Incoporating_Concurrency(input []int) []int {
	const batchSize = 1000
	inputLen := len(input)
	output := make([]int, inputLen)

	var wg sync.WaitGroup

	for i := 0; i < inputLen; i += batchSize {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			l := i + batchSize
			if l > inputLen {
				l = inputLen
			}

			for ; i < l; i++ {
				j := inputLen - i - 1
				n := input[i]

				output[j] = n
			}
		}(i)
	}
	wg.Wait()

	return output
}

func Reverse_Interger() {
	var n int
	var reverse int = 0
	fmt.Println("Enter a number")
	fmt.Scanln(&n) // Take input from user
	for n != 0 {
		reverse = reverse * 10
		reverse = reverse + n%10
		n = (n / 10)
	}
	fmt.Println(reverse)
}
