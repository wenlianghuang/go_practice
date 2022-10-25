package receiverex

import (
	"fmt"
	"math"
)

type Rectangle struct {
	length int
	width  int
}

type Circle struct {
	radias float64
}

func (r Rectangle) Area() int {
	return r.length * r.width
}

func (c Circle) Area() float64 {
	return math.Pi * c.radias * c.radias
}
func Receiver() {
	r := Rectangle{
		length: 10,
		width:  5,
	}
	fmt.Printf("Area of rectangle %d\n", r.Area())
	c := Circle{
		radias: 12,
	}
	fmt.Printf("Area of circle %.6v\n", c.Area())
}
