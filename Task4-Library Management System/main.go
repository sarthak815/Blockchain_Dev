package main

import (
	"bufio" // To read lines with whitespace
	"fmt"
	"os"
)

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
type BookType int

const (
	eBook        BookType = iota //digital
	Audiobook                    //digital
	Hardback                     //physical
	Paperback                    //physical
	Encyclopedia                 //both
	Magazine                     //both
	Comic                        //both
)

//Book struct to hold digital and physical books
type Books struct {
	B_type   BookType
	B_Name   string
	B_Author string
	Capacity int
	Borrowed int
}

//book interface with all functions
type Book interface {
	Booktype() string
	Name() string
	Author() string
	Borrow() bool //checks if available to borrow
	Return()      //returns a particular book
}

//Book constructor to create a book object for digital and physical
func (book *Books) Init(btype BookType, name, author string, capacity int) {
	book.B_type = btype
	book.B_Name = name
	book.B_Author = author
	book.Capacity = capacity
	book.Borrowed = 0 //Initially always available to borrow
}

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

func removeBookMember(slice []Books, s int) []Books {
	return append(slice[:s], slice[s+1:]...)
}

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

//Checks if the user wishing to borrow or return is registered
func checkUserValidity(name string, lib Library) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			break
		}
	}
	return b, member
}

//Checks if book is present in the library
func checkBookValidity(bname string, lib Library) (bool, *Books) {
	bfound := false      //Denotes book validity
	var bookFound *Books //Used to return book object that user wishes to borrow or return
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bookFound = &lib.BooksBorrowed[i]
			bfound = true
			break
		}
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

//Member Structure
type Member struct {
	Name          string
	Age           int
	BooksBorrowed []Books
}

//Constructor initialising member struct
func (member *Member) Init(name string, age int) {
	member.Name = name
	member.Age = age

}

//Library strcuture
type Library struct {
	BooksBorrowed []Books
	Members       []Member
}

func main() {
	var lib Library
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter 1.Enter Librarian Interface\n2.Enter User Interface\n3.Exit")
		var n int
		fmt.Scanln(&n)
		if n == 1 {
			fmt.Println("Enter 1.Enter Book in LibraryDB\n2.Exit")
			var n int
			fmt.Scanln(&n)
			if n == 1 {
				bookS := new(Books)
				//var typeBook string
				var btype BookType
				fmt.Println("Enter Book Type \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
				fmt.Scanln(&btype)
				btype--
				fmt.Println("Enter Name of book: ")
				scanner.Scan()
				name := scanner.Text()
				fmt.Println("Enter Author of book: ")
				scanner.Scan()
				author := scanner.Text()

				if btype <= 1 {
					var capac int
					fmt.Println("Enter borrowing limit: ")
					fmt.Scanln(&capac)
					bookS.Init(btype, name, author, capac)
					fmt.Println(bookS)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *bookS)

				}
				if btype > 1 && btype <= 3 {
					bookS.Init(btype, name, author, 1)
					fmt.Println(bookS)

				}
				if btype > 3 {
					var capac int
					fmt.Println("Enter borrowing limit for digital: ")
					fmt.Scanln(&capac)
					bookS.Init(btype, name+"Digital", author, capac)
					fmt.Println(bookS)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *bookS)
					bookS.Init(btype, name+"Physical", author, 1)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *bookS)
				}
				fmt.Println(lib)
			}
			if n == 2 {
				break
			}

		}
		if n == 2 {
			fmt.Println("1.Enter Member Details\n2.Borrow a book\n3.Return Book\n4.Exit")
			var n int
			fmt.Scanln(&n)
			if n == 1 {
				member := new(Member)
				fmt.Println("Enter name:")
				scanner.Scan()
				name := scanner.Text()
				var age int
				fmt.Println("Enter age:")
				fmt.Scanln(&age)
				member.Init(name, age)
				lib.Members = append(lib.Members, *member)
				fmt.Println(lib)

			}
			if n == 2 {
				var verifiedMember *Member
				var b bool
				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, verifiedMember = checkUserValidity(name, lib)
				if b {
					fmt.Println("Identity verified")
					var btype BookType
					fmt.Println("Enter BookType: \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
					fmt.Scanln(&btype)
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text()
					var digPhy int
					if btype > 4 {

						fmt.Println("Enter Copy Type:\n1.Digital\n2.Physical")
						fmt.Scanln(&digPhy)
						if digPhy == 1 {
							bname += "Digital"
						}
						if digPhy > 1 || digPhy < 1 {
							bname += "Physical"
						}
					}
					bfound := false
					var bookFound *Books
					bfound, bookFound = checkBookValidity(bname, lib)
					if bfound {
						if !bookFound.Borrow() {
							fmt.Println("Book Unavailable")
							continue
						}
						fmt.Println(*bookFound)
						verifiedMember.BooksBorrowed = append(verifiedMember.BooksBorrowed, *bookFound)
						printBookDetails(bookFound)
						fmt.Println(lib)
					}
					if !bfound {
						fmt.Println("Book Not found!!")
					}

				} else {
					fmt.Println("Details not found!! Please Register yourself!")
				}
			}
			if n == 3 {

				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, member := checkUserValidity(name, lib)
				if b {
					fmt.Println("Identity verified")
					var btype BookType
					fmt.Println("Enter BookType: \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
					fmt.Scanln(&btype)
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text()
					var digPhy int
					if btype > 4 {

						fmt.Println("Enter Copy Type:\n1.Digital\n2.Physical")
						fmt.Scanln(&digPhy)
						if digPhy == 1 {
							bname += "Digital"
						}
						if digPhy > 1 || digPhy < 1 {
							bname += "Physical"
						}
					}
					bfound := false
					var bookFound *Books
					bfound, bookFound = checkBookValidity(bname, lib)
					if bfound {
						bookFound.Return()
						idx := bookIndex(member.BooksBorrowed, *bookFound)
						member.BooksBorrowed = removeBookMember(member.BooksBorrowed, idx)
						fmt.Println("Book Returned")
						fmt.Println(lib)
					}
					if !bfound {
						fmt.Println("Book Not found!!")
					}

				} else {
					fmt.Println("Details not found!! Please Register yourself!")
				}
			}
			if n == 4 {
				break
			}
		}
		if n == 3 {
			break
		}

	}

}
