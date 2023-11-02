package interfacefolder

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human  //匿名欄位
	school string
	loan   float32
}

type Employee struct {
	Human   //匿名欄位
	company string
	money   float32
}

// Human 實現 SayHi 方法
func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

// Human 實現 Sing 方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

// Employee 過載 Human 的 SayHi 方法
func (e Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone)
}

// Interface Men 被 Human,Student 和 Employee 實現
// 因為這三個型別都實現了這兩個方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}

func Interfacevalue() {
	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 100}
	sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 1000}
	tom := Employee{Human{"Tom", 37, "222-444-XXX"}, "Things Ltd.", 5000}

	//定義 Men 型別的變數 i

	var i Men

	//i 能儲存 Student

	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")

	//i 也能儲存 Employee

	i = tom
	fmt.Println("This is tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	//定義了 slice Men
	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	//這三個都是不同型別的元素，但是他們實現了 interface 同一個介面
	x[0], x[1], x[2] = paul, sam, mike

	for _, value := range x {
		value.SayHi()
	}
}
