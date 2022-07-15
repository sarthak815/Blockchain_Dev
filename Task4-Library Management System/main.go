package main

import (
	"bufio"
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
	eBook BookType = iota
	Audiobook
	Hardback
	Paperback
	Encyclopedia
	Magazine
	Comic
)

type Books struct {
	B_type   BookType
	B_Name   string
	B_Author string
	Capacity int
	Borrowed int
}

type Book interface {
	Booktype() string
	Name() string
	Author() string
	Borrow() bool
}

func createPhybook(btype BookType, name, author string) Books {
	var phybook Books
	phybook.B_type = btype
	phybook.B_Name = name
	phybook.B_Author = author
	phybook.Capacity = 1
	phybook.Borrowed = 0
	return phybook
}
func createDigbook(btype BookType, name, author string, capacity int) Books {
	var digBook Books
	digBook.B_type = btype
	digBook.B_Name = name
	digBook.B_Author = author
	digBook.Capacity = capacity
	digBook.Borrowed = 0
	return digBook
}

func (b Books) Booktype() string {
	return [...]string{"eBook", "Audiobook", "Hardback", "Paperback", "Encyclopedia", "Magazine", "Comic"}[b.B_type]
}
func (b Books) Name() string {
	return b.B_Name
}
func (b Books) Author() string {
	return b.B_Author
}
func (b *Books) Borrow() bool {
	borrowed := true
	if b.Capacity <= b.Borrowed {
		borrowed = false
	} else {
		b.Borrowed++
	}
	return borrowed
}

type Member struct {
	Name          string
	Age           int
	BooksBorrowed []Books
}

func createMember(name string, age int) Member {
	var member Member
	member.Name = name
	member.Age = age
	return member
}

type Library struct {
	BooksBorrowed []Books
	Members       []Member
}

func main() {
	var lib Library
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter 1.Enter Book in LibraryDB\n2.Enter Member Details\n3.Borrow a book\n4.Exit")
		var n int
		fmt.Scanln(&n)
		if n == 1 {
			var bookS Books
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
				bookS = createDigbook(btype, name, author, capac)
				fmt.Println(bookS)
				lib.BooksBorrowed = append(lib.BooksBorrowed, bookS)

			}
			if btype > 1 && btype <= 3 {
				bookS = createPhybook(btype, name, author)
				fmt.Println(bookS)

			}
			if btype > 3 {
				var capac int
				fmt.Println("Enter borrowing limit for digital: ")
				fmt.Scanln(&capac)
				bookS = createDigbook(btype, name+"Digital", author, capac)
				lib.BooksBorrowed = append(lib.BooksBorrowed, bookS)
				bookS = createPhybook(btype, name+"Physical", author)
				lib.BooksBorrowed = append(lib.BooksBorrowed, bookS)
			}
			fmt.Println(lib)
		}
		if n == 2 {
			var member Member
			fmt.Println("Enter name:")
			scanner.Scan()
			name := scanner.Text()
			var age int
			fmt.Scanln(&age)
			member = createMember(name, age)
			lib.Members = append(lib.Members, member)
			fmt.Println(lib)

		}
		if n == 3 {
			var vermember *Member
			fmt.Println("Enter your name: ")
			scanner.Scan()
			name := scanner.Text()
			b := false
			for i := range lib.Members {
				if lib.Members[i].Name == name {
					vermember = &lib.Members[i]
					b = true
					break
				}
			}
			if b {
				fmt.Println("Identity verified")
				var btype BookType
				fmt.Println("Enter BookType: \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
				fmt.Scanln(&btype)
				if btype > 4 {
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text()
					var digPhy int
					fmt.Println("Enter Copy Type:\n1.Digital\n2.Physical")
					fmt.Scanln(&digPhy)
					if digPhy == 1 {
						bname += "Digital"
					}
					if digPhy > 1 || digPhy < 1 {
						bname += "Physical"
					}
					bfound := false
					var bookFound *Books
					for i := range lib.BooksBorrowed {
						if lib.BooksBorrowed[i].B_Name == bname {
							bookFound = &lib.BooksBorrowed[i]
							bfound = true
							break
						}
					}
					if bfound {
						if !bookFound.Borrow() {
							fmt.Println("Book Unavailable")
							continue
						}
						fmt.Println(*bookFound)
						fmt.Println("Details of the book: ")
						fmt.Println("Book Type: " + bookFound.Booktype())
						fmt.Println("Book Name: " + bookFound.Name())
						fmt.Println("Book Author: " + bookFound.Author())
						vermember.BooksBorrowed = append(vermember.BooksBorrowed, *bookFound)
						fmt.Println("Book Issued")
						fmt.Println(lib)
					}
					if !bfound {
						fmt.Println("Book Not found!!")
					}
				} else {
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text()
					bfound := false
					var bookFound *Books
					for i := range lib.BooksBorrowed {
						if lib.BooksBorrowed[i].B_Name == bname {
							bookFound = &lib.BooksBorrowed[i]
							bfound = true
							break
						}
					}
					if bfound {
						if !bookFound.Borrow() {
							fmt.Println("Book Unavailable")
							continue
						}
						fmt.Println(*bookFound)
						fmt.Println("Details of the book: ")
						fmt.Println("Book Type: " + bookFound.Booktype())
						fmt.Println("Book Name: " + bookFound.Name())
						fmt.Println("Book Author: " + bookFound.Author())
						vermember.BooksBorrowed = append(vermember.BooksBorrowed, *bookFound)
						fmt.Println("Book Issued")
						fmt.Println(lib)
					}
					if !bfound {
						fmt.Println("Book Not found!!")
					}
				}
			} else {
				fmt.Println("Details not found!! Please Register yourself!")
			}
		}
		if n == 4 {
			break
		}

	}

}
