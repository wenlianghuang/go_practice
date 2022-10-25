package mapstruct

import (
	"fmt"
	"sort"
)

type Student struct {
	ID   int
	Name string
}

func MapStruct() {
	var student1 Student
	student1.ID = 1
	student1.Name = "Daniel"

	student2 := new(Student)
	student2.ID = 2
	student2.Name = "Sam"

	student3 := Student{3, "Allen"}

	students := make(map[string]Student)
	students["s1"] = student1
	students["s2"] = *student2
	students["s3"] = student3

	var names []string

	for studentIndex, studentobj := range students {
		fmt.Println(studentIndex)
		fmt.Println(studentobj)
		names = append(names, studentobj.Name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Println(name)
	}
}
