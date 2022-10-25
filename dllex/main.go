package main

import (
	"dllex/dllfunc"
	"dllex/sllfunc"
	"fmt"
)

func main() {
	doublyList := dllfunc.InitDoublyList()
	fmt.Printf("Add Front Node: C\n")
	doublyList.AddFrontNodeDLL("C")
	fmt.Printf("Add Front Node: B\n")
	doublyList.AddFrontNodeDLL("B")
	fmt.Printf("Add Front Node: A\n")
	doublyList.AddFrontNodeDLL("A")
	fmt.Printf("Add End Node: D\n")
	doublyList.AddEndNodeDLL("D")
	fmt.Printf("Add End Node: E\n")
	doublyList.AddEndNodeDLL("E")

	fmt.Printf("Size of doubly linked ist: %d\n", doublyList.Size())

	err := doublyList.TraverseForward()
	if err != nil {
		panic(err)
	}

	err = doublyList.TraverseReverse()
	if err != nil {
		panic(err)
	}

	link := sllfunc.List{}
	//link := sllfunc.InitSingleList()
	link.Insert(5)
	link.Insert(9)
	link.Insert(13)
	link.Insert(22)
	link.Insert(28)
	link.Insert(36)

	fmt.Printf("\n==============================\n")
	link.Displayhead()
	link.Displaytail()
	link.Display()
	fmt.Printf("\n==============================\n")

	link.Reverse()
	link.Displayhead()
	link.Displaytail()
	fmt.Printf("\n==============================\n")
}
