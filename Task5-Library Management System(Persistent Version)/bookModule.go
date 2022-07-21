package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
)

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
	B_type   BookType `json:"booktype"`
	B_Name   string   `json:"name"`
	B_Author string   `json:"author"`
	Capacity int      `json:"capacity"`
	Borrowed int      `json:"borrowed"`
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
func checkBookValidity(bname string, lib *Library, member *Member, db *badger.DB) (bool, *Books) {
	bfound := false //Denotes book validity
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bookFound = &lib.BooksBorrowed[i]
			bfound = true
			return bfound, bookFound
		}
	}
	if err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if k == bname {
				err := item.Value(func(v []byte) error {
					// Create an encoder and send a value.
					enc := gob.NewDecoder(bytes.NewBuffer(v))
					err := enc.Decode(&bookFound)
					lib.BooksBorrowed = append(lib.BooksBorrowed, *bookFound)
					if err != nil {
						log.Fatal("Error in decoding user validity return:", err)
					}

					return nil
				})
				if err != nil {
					return err
				}
				bfound = true
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println("DB Reading Error on library module")
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
func checkBookValidityApi(bname string, lib *Library, db *badger.DB) bool {
	fmt.Println(bname)
	bfound := false //Denotes book validity
	for i := range lib.BooksBorrowed {
		if lib.BooksBorrowed[i].B_Name == bname {
			bfound = true
			return bfound
		}
	}
	if err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())
			fmt.Println(k)
			if k == bname {
				bfound = true
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println("DB Reading Error on library module")
	}

	return bfound
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
