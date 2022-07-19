package main

import "fmt"

//book interface with all functions
type Book interface {
	Booktype() string
	Name() string
	Author() string
	Borrow() bool //checks if available to borrow
	Return()      //returns a particular book
}

//Book struct to hold digital and physical books
type Books struct {
	B_type   BookType
	B_Name   string
	B_Author string
	Capacity int
	Borrowed int
}

type BookType int

//book type enum
const (
	eBook        BookType = iota //digital
	Audiobook                    //digital
	Hardback                     //physical
	Paperback                    //physical
	Encyclopedia                 //both
	Magazine                     //both
	Comic                        //both
)

//Book constructor to create a book object for digital and physical
func (book *Books) Init(btype BookType, name, author string, capacity int) {
	book.B_type = btype
	book.B_Name = name
	book.B_Author = author
	book.Capacity = capacity
	book.Borrowed = 0 //Initially always available to borrow
}

//*********************INTERFACE FUNCTION IMPLEMENTATION*********************
//Booktype() returns the string value of book type
func (b Books) Booktype() string {
	return [...]string{"eBook", "Audiobook", "Hardback", "Paperback", "Encyclopedia", "Magazine", "Comic"}[b.B_type]
}

//Name() returns book name
func (b Books) Name() string {
	return b.B_Name
}

//Author() returns author name
func (b Books) Author() string {
	return b.B_Author
}

//Borrow() returns bool value indicating whether book is available to borrow
func (b *Books) Borrow() bool {
	borrowed := true
	if b.Capacity <= b.Borrowed {
		borrowed = false
	} else {
		b.Borrowed++
	}
	return borrowed
}

//Return() Returns a book object by decrementing the borrowed value from library struct
func (b *Books) Return() {
	b.Borrowed--
}

//BookIndex() returns the position of the book to be removed from the member book slice
func bookIndex(slice []Books, book Books) int {
	idx := -1
	for i := range slice {
		if slice[i].B_Name == book.B_Name {
			idx = i
			break
		}
	}
	return idx
}

//Checks if book is present in the library and user eligible to borrow
func checkBookValidity(bname string, lib Library, member *Member) (bool, *Books) {
	bfound := false //Denotes book validity
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bookFound = &lib.BooksBorrowed[i]
			bfound = true
			break
		}
	}
	if bfound { //used for out of index error handling
		for i := range member.BooksBorrowed { // checks if user has borrowed book already
			if member.BooksBorrowed[i].B_Name == bookFound.B_Name {
				borrowed = true
			}
		}
	}
	if borrowed {
		bfound = false
	}
	return bfound, bookFound
}

//Prints details of the book
func printBookDetails(bookFound *Books) {
	fmt.Println("Details of the book: ")
	fmt.Println("Book Type: " + bookFound.Booktype())
	fmt.Println("Book Name: " + bookFound.Name())
	fmt.Println("Book Author: " + bookFound.Author())
	fmt.Println("Book Issued")

}

//removeBookMember() removes a particular book from member's books slice
func removeBookMember(slice []Books, s int) []Books {
	return append(slice[:s], slice[s+1:]...)
}
