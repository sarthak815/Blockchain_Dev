package main

//Member Structure
type Member struct {
	Name          string
	Age           int
	BooksBorrowed []Books
}

//*********************STRUCT CONSTRUCTORS*********************
//Constructor initialising member struct
func (member *Member) Init(name string, age int) {
	member.Name = name
	member.Age = age

}
