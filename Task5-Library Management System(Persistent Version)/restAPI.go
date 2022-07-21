package main

import (
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
	writeBooksToDB()
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
	writeMembersToDB()
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
	writeMembersToDB()
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
	writeBooksToDB()
	writeMembersToDB()
	json.NewEncoder(w).Encode(borrower)
}
func writeBooksToDB() {
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
}
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
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/book", createNewBook).Methods("POST")
	myRouter.HandleFunc("/user", createNewMember).Methods("POST")
	myRouter.HandleFunc("/borrow", borrowBook).Methods("POST")
	myRouter.HandleFunc("/return", returnBook).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
