package main

import (
	"Task5-Library_Management_System/codeModules"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

//homePage sets the text to be displayed on the localhost page
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//createNewBook creates a new book object and stores it to the database
func createNewBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Books struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var bookJson codeModules.Books
	json.Unmarshal(reqBody, &bookJson)
	// update our bookJSON variable to hold the object received over API
	//checkBookValidityApi checks if the book received is valid to be stored to the db
	b := codeModules.CheckBookValidityApi(bookJson.B_Name, &lib, db)
	if b {
		log.Println("Book found in DB already")
		return
	}

	//In case book is of physical type ensures only one copy is present
	if bookJson.B_type > 1 && bookJson.B_type <= 3 {
		bookJson.Capacity = 1

	}
	lib.BooksAvailable = append(lib.BooksAvailable, bookJson)
	fmt.Println("Details of new book: ")
	codeModules.PrintBookDetails(&bookJson)
	//writeBooksToDB stores the newly added book to the database
	writeBooksToDB()
	//returns the json object as received to the api
	json.NewEncoder(w).Encode(bookJson)
}

//createNewMember creates a new member object and stores it to the database
func createNewMember(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new member struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var memberJson codeModules.Member
	json.Unmarshal(reqBody, &memberJson)
	//checkMemberValidityApi ensures that member is valid to be added to the DB
	b := codeModules.CheckMemberValidityApi(memberJson.Name, &lib, db)
	if b {
		log.Println("Member found in DB already")
		return
	}
	//If valid the new member is appended to the lib struct
	lib.Members = append(lib.Members, memberJson)
	//writeMembersToDB writes the newly added member to the BadgerDB
	writeMembersToDB()
	json.NewEncoder(w).Encode(memberJson)
}

//borrowBook takes json input of type Books and allows an user to borrow a book while performing necessary checks
func borrowBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Borrower struct

	reqBody, _ := ioutil.ReadAll(r.Body)
	var borrower codeModules.Borrower
	json.Unmarshal(reqBody, &borrower)
	name := borrower.Name

	// verifies member details
	b, verifiedMember := codeModules.CheckUserValidity(name, &lib, db) //checks username validity and number of books that user has borrowed is below 5
	if b {                                                             //if user is valid and registered
		fmt.Println("Identity verified")
		bname := borrower.BookN //stores book name to be borrowed

		bfound, bookFound := codeModules.CheckBookValidity(bname, &lib, verifiedMember, db) //checks if book is available then stores book pointed to bookFound
		if bfound {
			if !bookFound.Borrow() { //Borrow() checks if book id available to borrow
				fmt.Println("Book Unavailable")
				return
			}
			verifiedMember.BooksBorrowed = append(verifiedMember.BooksBorrowed, *bookFound)
			codeModules.PrintBookDetails(bookFound) //displays details of book borrowed
		}
		if !bfound { //in case book not present in struct
			fmt.Println("Book Not found/User already borrowed 5 books/Requested book borrowed!!")
		}

	} else { //in case user details not present in struct
		fmt.Println("Details not found or Already reached limit!")
	}

	writeMembersToDB()
	json.NewEncoder(w).Encode(borrower)
}

//returnBook takes json input of type Borrower and allows an user to return a borrowed book
func returnBook(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new borrower struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var borrower codeModules.Borrower
	json.Unmarshal(reqBody, &borrower)
	name := borrower.Name
	b, member := codeModules.CheckUserValidityReturn(name, &lib, db) //checks validity of user
	if b {
		fmt.Println("Identity verified")
		bfound, bookFound := codeModules.CheckUserBookValidity(borrower.BookN, lib, *member, db) //checks if book is available and borrowed then stores book pointed to bookFound
		if bfound {
			bookFound.Return()                                                             //modifies book borrowed value
			idx := codeModules.BookIndex(member.BooksBorrowed, *bookFound)                 //obtains index of borrowed book in user struct
			member.BooksBorrowed = codeModules.RemoveBookMember(member.BooksBorrowed, idx) //removes returned book from user struct
			fmt.Println("Book Returned")
		}
		if !bfound { //in case book not found in directory
			fmt.Println("Book Not found or not borrowed!!")
		}

	} else { //if user details not present in library
		fmt.Println("Details not found!! Please Register yourself!")
	}
	writeBooksToDB()
	writeMembersToDB()
	json.NewEncoder(w).Encode(borrower)
}

//writeBooksToDB saves all books data in Library struct to BadgerDB
func writeBooksToDB() {
	for i := range lib.BooksAvailable {
		var bookBytes bytes.Buffer // Stand-in for the bookBytes.
		// Create an encoder and send a value.
		enc := gob.NewEncoder(&bookBytes)
		err := enc.Encode(lib.BooksAvailable[i])
		if err != nil {
			log.Fatal("encode:", err)
		}
		//creates a new transaction
		txn := db.NewTransaction(true)
		defer txn.Discard()
		e := badger.NewEntry([]byte(lib.BooksAvailable[i].Name()), bookBytes.Bytes())
		_ = txn.SetEntry(e)

		_ = txn.Commit()

		fmt.Println("Inserted books")
	}
}

//writeMembersToDB saves all members data in Library struct to BadgerDB
func writeMembersToDB() {
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
}

//handleRequests contains all the functions necessary for handling API calls using gorilla mux
func handleRequests() {
	//creates a gorilla mux to handle different paths to access variety ogf functions
	myRouter := mux.NewRouter().StrictSlash(true)
	//homepage to verify the api is working
	myRouter.HandleFunc("/", homePage)
	//Endpoint to insert book
	myRouter.HandleFunc("/book", createNewBook).Methods("POST")
	//Endpoint to insert new user
	myRouter.HandleFunc("/user", createNewMember).Methods("POST")
	//Endpoint to borrow a book
	myRouter.HandleFunc("/borrow", borrowBook).Methods("POST")
	//Endpoint to return a borrowed book
	myRouter.HandleFunc("/return", returnBook).Methods("POST")
	//sets the port number to listed to requests
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
