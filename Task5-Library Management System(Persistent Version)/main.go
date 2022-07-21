package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var lib Library
var db *badger.DB

type Borrower struct {
	Name  string   `json:"name"`
	BookT BookType `json:"type"`
	BookN string   `json:"bookname"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
func createNewBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var bookJson Books
	json.Unmarshal(reqBody, &bookJson)
	// update our global Articles array to include
	// our new Article
	b := checkBookValidityApi(bookJson.B_Name, &lib, db)
	if b {
		log.Println("Book found in DB already")
		return
	}
	fmt.Println(bookJson)
	//Indicates book to be of Physical type
	if bookJson.B_type > 1 && bookJson.B_type <= 3 {
		bookJson.Capacity = 1

	}
	lib.BooksBorrowed = append(lib.BooksBorrowed, bookJson)
	fmt.Println(lib)
	for i := range lib.BooksBorrowed {
		var bookBytes bytes.Buffer // Stand-in for the bookBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&bookBytes)
		err := enc.Encode(lib.BooksBorrowed[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		txn := db.NewTransaction(true)
		defer txn.Discard()
		e := badger.NewEntry([]byte(lib.BooksBorrowed[i].Name()), bookBytes.Bytes())
		_ = txn.SetEntry(e)

		_ = txn.Commit()

		fmt.Println("Inserted books")
	}
	json.NewEncoder(w).Encode(bookJson)
}
func createNewMember(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var memberJson Member
	json.Unmarshal(reqBody, &memberJson)
	// update our global Articles array to include
	// our new Article
	b := checkMemberValidityApi(memberJson.Name, &lib, db)
	if b {
		log.Println("Member found in DB already")
		return
	}
	fmt.Println(memberJson)

	lib.Members = append(lib.Members, memberJson)
	fmt.Println(lib)
	for i := range lib.Members {
		var memberBytes bytes.Buffer // Stand-in for the memberBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&memberBytes)
		err := enc.Encode(lib.Members[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		txn := db.NewTransaction(true)
		defer txn.Discard()
		if err := txn.Set([]byte(lib.Members[i].Name), memberBytes.Bytes()); err != nil {
			log.Println("Commmit Error")
		}

		if err := txn.Commit(); err != nil {
			log.Println("Commmit Error")
		}

		fmt.Println("Inserted Members")
	}
	json.NewEncoder(w).Encode(memberJson)
}
func borrowBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var borrower Borrower
	json.Unmarshal(reqBody, &borrower)
	name := borrower.Name
	fmt.Println("Name:" + name)
	// update our global Articles array to include
	// our new Article
	b, verifiedMember := checkUserValidity(name, &lib, db) //checks username validity and number of books that user has borrowed is below 5
	b, verifiedMember = checkUserValidity(name, &lib, db)
	if b { //if user is valid and registered
		fmt.Println("Identity verified")
		bname := borrower.BookN //stores book name to be borrowed

		bfound, bookFound := checkBookValidity(bname, &lib, verifiedMember, db) //checks if book is available then stores book pointed to bookFound
		bfound, bookFound = checkBookValidity(bname, &lib, verifiedMember, db)
		if bfound {
			if !bookFound.Borrow() { //Borrow() checks if book id available to borrow
				fmt.Println("Book Unavailable")
				return
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
	fmt.Println(borrower)

	fmt.Println(lib)
	for i := range lib.Members {
		var memberBytes bytes.Buffer // Stand-in for the memberBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&memberBytes)
		err := enc.Encode(lib.Members[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		txn := db.NewTransaction(true)
		defer txn.Discard()
		if err := txn.Set([]byte(lib.Members[i].Name), memberBytes.Bytes()); err != nil {
			log.Println("Commmit Error")
		}

		if err := txn.Commit(); err != nil {
			log.Println("Commmit Error")
		}

		fmt.Println("Inserted Members")
	}
	json.NewEncoder(w).Encode(borrower)
}
func returnBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var borrower Borrower
	json.Unmarshal(reqBody, &borrower)
	name := borrower.Name
	fmt.Println("Name:" + name)
	b, member := checkUserValidityReturn(name, &lib, db) //checks validity of user
	if b {
		fmt.Println("Identity verified")
		bfound, bookFound := checkUserBookValidity(borrower.BookN, lib, *member, db) //checks if book is available and borrowed then stores book pointed to bookFound
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
	for i := range lib.BooksBorrowed {
		var bookBytes bytes.Buffer // Stand-in for the bookBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&bookBytes)
		err := enc.Encode(lib.BooksBorrowed[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		txn := db.NewTransaction(true)
		defer txn.Discard()
		e := badger.NewEntry([]byte(lib.BooksBorrowed[i].Name()), bookBytes.Bytes())
		_ = txn.SetEntry(e)

		_ = txn.Commit()

		fmt.Println("Inserted books")
	}
	for i := range lib.Members {
		var memberBytes bytes.Buffer // Stand-in for the memberBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&memberBytes)
		err := enc.Encode(lib.Members[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		txn := db.NewTransaction(true)
		defer txn.Discard()
		if err := txn.Set([]byte(lib.Members[i].Name), memberBytes.Bytes()); err != nil {
			log.Println("Commmit Error")
		}

		if err := txn.Commit(); err != nil {
			log.Println("Commmit Error")
		}

		fmt.Println("Inserted Members")
	}
	json.NewEncoder(w).Encode(borrower)
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/book", createNewBook).Methods("POST")
	myRouter.HandleFunc("/user", createNewMember).Methods("POST")
	myRouter.HandleFunc("/borrow", borrowBook).Methods("POST")
	myRouter.HandleFunc("/return", returnBook).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	db, _ = badger.Open(badger.DefaultOptions("C:\\Users\\Sanjay\\OneDrive - Manipal Academy of Higher Education\\Documents\\GitHub\\Blockchain_Dev\\Task5-Library Management System(Persistent Version)\\tmp\\badger"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	defer db.Close()
	//To keep program running for multiple operations
	for {
		fmt.Println("Enter 1.Enter Librarian Interface\n2.Enter User Interface\n3.Enter API\n4.Exit")
		var n int
		fmt.Scanln(&n) // Option choice stored\
		switch n {
		// Enter Library Management System
		case 1:
			fmt.Println("Enter 1.Enter Book in LibraryDB\n2.Exit")
			var n int
			fmt.Scanln(&n)
			switch n {
			//Case 1 enters new book in library
			case 1:
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
				// Indicates book to be of type eBook
				if bookType <= 1 {
					var capacity int
					fmt.Println("Enter borrowing limit: ") //Total copies available to borrow for digital
					fmt.Scanln(&capacity)
					newBook.Init(bookType, name, author, capacity)
					log.Println(newBook) //
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)

				}
				//Indicates book to be of Physical type
				if bookType > 1 && bookType <= 3 {
					newBook.Init(bookType, name, author, 1) // capaccity set to 1 as physical copy can only be 1 piece
					log.Println(newBook)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *newBook)

				}
				//Indicates book to be of Physical and Digital
				if bookType > 3 {
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

			// exit clause to quit library management interface
			case 2:
				break

			}
		//Enters user interface to register/borrow/return
		case 2:
			fmt.Println("1.Enter Member Details\n2.Borrow a book\n3.Return Book\n4.Exit")
			var n int
			fmt.Scanln(&n) //stores user choice
			switch n {
			case 1: //User registration portal
				member := new(Member) //creates new member
				fmt.Println("Enter name:")
				scanner.Scan()
				name := scanner.Text()
				b := checkUserNameValidity(name, lib, db) // checks if username is valid
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
			case 2: //Book borrowing portal
				var verifiedMember *Member //used to obtain user details
				var b bool                 //checks if username is valid
				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, verifiedMember = checkUserValidity(name, &lib, db) //checks username validity and number of books that user has borrowed is below 5
				b, verifiedMember = checkUserValidity(name, &lib, db)
				if b { //if user is valid and registered
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
					bfound, bookFound := checkBookValidity(bname, &lib, verifiedMember, db) //checks if book is available then stores book pointed to bookFound
					bfound, bookFound = checkBookValidity(bname, &lib, verifiedMember, db)
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
			case 3: //Portal for user to return book

				fmt.Println("Enter your name: ")
				scanner.Scan()
				name := scanner.Text()
				b, member := checkUserValidityReturn(name, &lib, db) //checks validity of user
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

					bfound, bookFound := checkUserBookValidity(bname, lib, *member, db) //checks if book is available and borrowed then stores book pointed to bookFound
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
			case 4: //exit clause for user portal
				break
			}

		case 3:
			handleRequests()
		//exit clause to close application
		case 4:
			for i := range lib.BooksBorrowed {
				var bookBytes bytes.Buffer // Stand-in for the bookBytes.
				// Create an encoder and send a value.
				enc := gob.NewEncoder(&bookBytes)
				err := enc.Encode(lib.BooksBorrowed[i])
				if err != nil {
					log.Fatal("encode:", err)
				}
				txn := db.NewTransaction(true)
				defer txn.Discard()
				e := badger.NewEntry([]byte(lib.BooksBorrowed[i].Name()), bookBytes.Bytes())
				_ = txn.SetEntry(e)

				_ = txn.Commit()

				fmt.Println("Inserted books")
			}
			for i := range lib.Members {
				var memberBytes bytes.Buffer // Stand-in for the memberBytes.
				// Create an encoder and send a value.
				enc := gob.NewEncoder(&memberBytes)
				err := enc.Encode(lib.Members[i])
				if err != nil {
					log.Fatal("encode:", err)
				}
				txn := db.NewTransaction(true)
				defer txn.Discard()
				if err := txn.Set([]byte(lib.Members[i].Name), memberBytes.Bytes()); err != nil {
					log.Println("Commmit Error")
				}

				if err := txn.Commit(); err != nil {
					log.Println("Commmit Error")
				}

				fmt.Println("Inserted Members")
			}
			os.Exit(0)
		}

	}
	//handleRequests()
}