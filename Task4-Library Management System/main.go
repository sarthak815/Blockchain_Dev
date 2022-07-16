package main

import (
	"bufio" // To read lines with whitespace
	"fmt"
	"log"
	"os"
)

//*********************ENUMS*********************

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

//*********************STRUCTS*********************
//Book struct to hold digital and physical books
type Books struct {
	B_type   BookType
	B_Name   string
	B_Author string
	Capacity int
	Borrowed int
}

//Member Structure
type Member struct {
	Name          string
	Age           int
	BooksBorrowed []Books
}

//Library structure
type Library struct {
	BooksBorrowed []Books
	Members       []Member
}

//*********************STRUCT CONSTRUCTORS*********************
//Constructor initialising member struct
func (member *Member) Init(name string, age int) {
	member.Name = name
	member.Age = age

}

//Book constructor to create a book object for digital and physical
func (book *Books) Init(btype BookType, name, author string, capacity int) {
	book.B_type = btype
	book.B_Name = name
	book.B_Author = author
	book.Capacity = capacity
	book.Borrowed = 0 //Initially always available to borrow
}

//*********************INTERFACES*********************
//book interface with all functions
type Book interface {
	Booktype() string
	Name() string
	Author() string
	Borrow() bool //checks if available to borrow
	Return()      //returns a particular book
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

//*********************FUNCTIONS USED IN MAIN FUNC*********************
//removeBookMember() removes a particular book from member's books slice
func removeBookMember(slice []Books, s int) []Books {
	return append(slice[:s], slice[s+1:]...)
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

//Checks if the user wishing to borrow is registered and eligible
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
	if b { //used for out of index error handling
		if len(member.BooksBorrowed) == 5 {
			b = false
		}
	}
	return b, member
}

//Checks if the user wishing to return is eligible
func checkUserValidityReturn(name string, lib Library) (bool, *Member) {
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

//checks if book borrowed and present in directory
func checkUserBookValidity(bname string, lib Library, member Member) (bool, *Books) {
	bfound := false //Denotes book validity
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow or return
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bookFound = &lib.BooksBorrowed[i]
			bfound = true
			break
		}
	}
	if bfound { //used for out of index error handling
		for i := range member.BooksBorrowed { // checks if user has borrowed same book
			if member.BooksBorrowed[i].B_Name == bookFound.B_Name {
				borrowed = true
			}
		}
	}
	if !borrowed { //incase book not borrowed by user
		bfound = false
	}
	return bfound, bookFound
}

//checks if username is unique as it acts as primary key
func checkUserNameValidity(username string, lib Library) bool {
	b := true
	for i := range lib.Members {
		if username == lib.Members[i].Name {
			b = false
			break
		}
	}
	return b
}

//Prints details of the book
func printBookDetails(bookFound *Books) {
	fmt.Println("Details of the book: ")
	fmt.Println("Book Type: " + bookFound.Booktype())
	fmt.Println("Book Name: " + bookFound.Name())
	fmt.Println("Book Author: " + bookFound.Author())
	fmt.Println("Book Issued")

}

func main() {
	var lib Library
	scanner := bufio.NewScanner(os.Stdin)
	// To keep program running for multiple operations
	for {
		fmt.Println("Enter 1.Enter Librarian Interface\n2.Enter User Interface\n3.Exit")
		var n int
		fmt.Scanln(&n) // Option choice stored
		// Enter Library Management System
		if n == 1 {
			fmt.Println("Enter 1.Enter Book in LibraryDB\n2.Exit")
			var n int
			fmt.Scanln(&n)
			if n == 1 {
				newBook := new(Books)
				//var typeBook string
				var bookType BookType
				fmt.Println("Enter Book Type \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
				fmt.Scanln(&bookType) //Stores book type to identify digital or physical
				bookType--            // Indexing of enum starts at 0
				fmt.Println("Enter Name of book: ")
				scanner.Scan()
				name := scanner.Text()
				fmt.Println("Enter Author of book: ")
				scanner.Scan()
				author := scanner.Text()
				if bookType <= 1 { // Indicates book to be of type eBook
					var capacity int
					fmt.Println("Enter borrowing limit: ") //Total copies available to borrow for digital
					fmt.Scanln(&capacity)
					newBook.Init(bookType, name, author, capacity)
					log.Println(newBook) //
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)

				}
				if bookType > 1 && bookType <= 3 { //Indicates book to be of Physical type
					newBook.Init(bookType, name, author, 1) // capaccity set to 1 as physical copy can only be 1 piece
					log.Println(newBook)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)

				}
				if bookType > 3 { //Indicates book to be of Physical and Digital
					var capacity int
					fmt.Println("Enter borrowing limit for digital: ")
					fmt.Scanln(&capacity)
					newBook.Init(bookType, name+"Digital", author, capacity) //Creating digital copy
					log.Println(newBook)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)
					newBook.Init(bookType, name+"Physical", author, 1) // Creating physical copy
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)
				}
				log.Println(lib) //logs changes made to library struct
			}
			if n == 2 { //exit clause to quit library management interface
				break
			}

		}
		//Enters user interface to register/borrow/return
		if n == 2 {
			fmt.Println("1.Enter Member Details\n2.Borrow a book\n3.Return Book\n4.Exit")
			var n int
			fmt.Scanln(&n) //stores user choice
			if n == 1 {    //User registration portal
				member := new(Member) //creates new member
				fmt.Println("Enter name:")
				scanner.Scan()
				name := scanner.Text()
				b := checkUserNameValidity(name, lib) // checks if username is valid
				if !b {
					fmt.Println("Username taken! Please Enter Full Name!")
					continue
				}
				var age int
				fmt.Println("Enter age:")
				fmt.Scanln(&age)
				member.Init(name, age)

				lib.Members = append(lib.Members, *member)
				log.Println(lib) //logs changes made to users

			}
			if n == 2 { //Book borrowing portal
				var verifiedMember *Member //used to obtain user details
				var b bool                 //checks if username is valid
				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, verifiedMember = checkUserValidity(name, lib) //checks username validity and number of books that user has borrowed is below 5
				if b {                                           //if user is valid and registered
					fmt.Println("Identity verified")
					var btype BookType
					fmt.Println("Enter BookType: \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
					fmt.Scanln(&btype) //checks for book type in case of book being available in both digital and physical versions
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text() //stores book name to be borrowed
					var digPhy int
					if btype > 4 { // In case of book type having physical and digital copies
						fmt.Println("Enter Copy Type:\n1.Digital\n2.Physical")
						fmt.Scanln(&digPhy)
						if digPhy == 1 {
							bname += "Digital" //appends digital to name in format in which stored in slice
						}
						if digPhy > 1 || digPhy < 1 {
							bname += "Physical" //appends physical to name in format in which stored in slice
						}
					}
					bfound, bookFound := checkBookValidity(bname, lib, verifiedMember) //checks if book is available then stores book pointed to bookFound
					if bfound {
						if !bookFound.Borrow() { //Borrow() checks if book id available to borrow
							fmt.Println("Book Unavailable")
							continue
						}
						verifiedMember.BooksBorrowed = append(verifiedMember.BooksBorrowed, *bookFound)
						printBookDetails(bookFound) //displays details of book borrowed
					}
					if !bfound { //in case book not present in struct
						fmt.Println("Book Not found/User already borrowed 5 books/Requested book borrowed!!")
					}

				} else { //in case user details not present in struct
					fmt.Println("Details not found or Already reached limit!")
				}
				fmt.Println(lib)
			}
			if n == 3 { //Portal for user to return book

				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, member := checkUserValidityReturn(name, lib) //checks validity of user
				if b {
					fmt.Println("Identity verified")
					var btype BookType
					fmt.Println("Enter BookType: \n1.Ebook\n2. AudioBook\n3. HardBack\n4. PaperBack\n5. Encyclopedia\n6. Magazine\n7. Comic: ")
					fmt.Scanln(&btype)
					fmt.Println("Enter Book name: ")
					scanner.Scan()
					bname := scanner.Text()
					var digPhy int
					if btype > 4 { //modifies book name in case of type where physical and digital copy is available

						fmt.Println("Enter Copy Type:\n1.Digital\n2.Physical")
						fmt.Scanln(&digPhy)
						if digPhy == 1 {
							bname += "Digital"
						}
						if digPhy > 1 || digPhy < 1 {
							bname += "Physical"
						}
					}

					bfound, bookFound := checkUserBookValidity(bname, lib, *member) //checks if book is available and borrowed then stores book pointed to bookFound
					if bfound {
						bookFound.Return()                                                 //modifies book borrowed value
						idx := bookIndex(member.BooksBorrowed, *bookFound)                 //obtains index of borrowed book in user struct
						member.BooksBorrowed = removeBookMember(member.BooksBorrowed, idx) //removes returned book from user struct
						fmt.Println("Book Returned")
						log.Println(lib)
					}
					if !bfound { //in case book not found in directory
						fmt.Println("Book Not found or not borrowed!!")
					}

				} else { //if user details not present in library
					fmt.Println("Details not found!! Please Register yourself!")
				}
			}
			if n == 4 { //exit clause for user portal
				break
			}
		}
		//exit clause to close application
		if n == 3 {
			break
		}

	}

}
