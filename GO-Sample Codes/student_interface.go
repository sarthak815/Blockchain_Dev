package main

import "fmt"

type Student interface {
	grade() int
	semester() int
}

type UGStudent struct {
	Grade    int
	Semester int
}

func (student UGStudent) grade() int {
	return student.Grade
}

func (student UGStudent) semester() int {
	return student.Semester
}

func main() {
	var st Student

	st = UGStudent{10, 4}
	fmt.Println("Grade is:", st.grade())
	fmt.Println("Semester is:", st.semester())
}
