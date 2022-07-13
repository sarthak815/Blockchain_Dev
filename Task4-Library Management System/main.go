package main

// The goal of this assignment is to implement a library management system which maintains an inventory of
// physical and digital books owned by a library as well as all members of the library.
// • Only a registered member of the library can borrow books from the library and can only borrow up
// to 5 books at a given point in time.
// • A book may be physical or digital in nature which lends itself to different borrowing constraints
// which is that physical books can only be borrowed by one member at a time while digital books
// contain a certain number of copies each of which can be borrowed by a member.
// • Only the member who has currently borrowed a book can return it to the library
// In an effort towards developing and modelling this library management system, you must use the Go
// Programming Language and the following language elements
// 1. Use enums to define a type BookType with variants such as eBook, Audiobook, Hardback,
// Paperback, Encyclopedia, Magazine, Comic, etc. Each of these types can be associated with
// either a physical or digital book (or both)
// 2. Book must be an Interface type with methods that return the kind, name and author of the book,
// whether it is a digital or physical book as well expose a method to set a borrower to it (returns a
// Boolean). When Borrow() is called on the Book to set the borrower
// a. PhysicalBook objects will only allow one to borrow and return false if already borrowed
// b. DigitalBook objects will allow multiple borrowers until their capacity is full and return
// false if capacity is full.
// 3. PhysicalBook and DigitalBook must be structs that satisfy the Book interface and implement
// their own constructors NewPhysicalBook and NewDigitalBook
// 4. Library must be a type that has methods to add new books to the inventory and register new
// members to the userbase. A member must have the ability to borrow a book from a Library and
// return it
// 5. You have freedom over the rest of the technical implementation details
type PhysicalBook struct {
	B_type   string
	B_Name   string
	B_Author string
	B_Borrow bool
}
type DigitalBook struct {
	B_type   string
	B_Name   string
	B_Author string
	Capacity int
	Borrowed int
	B_Borrow bool
}
type Book interface {
	Btype() string
	Name() string
	Author() string
	Borrow() bool
}
func (b PhysicalBook) Btype() string {
	return b.B_type
}
func (b PhysicalBook) Name() string {
	return b.B_Name
}
func (b PhysicalBook) Author() string {
	return b.B_Author
}
func (b PhysicalBook) Borrow() bool {
	return b.B_Borrow
}
func (b DigitalBook) Btype() string {
	return b.B_type
}
func (b DigitalBook) Name() string {
	return b.B_Name
}
func (b DigitalBook) Author() string {
	return b.B_Author
}
func (b DigitalBook) Borrow() bool {
	borrowed := true
	if b.Capacity <= b.Borrowed {
		borrowed = false
	} else {
		b.Borrowed++
	}
	return borrowed
}

type Member struct{
	Name string
	Age int
	Dig_borrowed DigitalBook[]
	Phy_borrowed PhysicalBook[]
	
}

type Library struct{
	Phy_Books PhysicalBook[]
	Dig_Books DigitalBook[]


}

func main() {

}
