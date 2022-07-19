package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
)

//Library structure
type Library struct {
	BooksBorrowed []Books
	Members       []Member
}

//Checks if the user wishing to borrow is registered and eligible
func checkUserValidity(name string, lib *Library, db *badger.DB) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			return b, member
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
			if k == name {
				err := item.Value(func(v []byte) error {
					// Create an encoder and send a value.
					enc := gob.NewDecoder(bytes.NewBuffer(v))
					err := enc.Decode(&member)

					if err != nil {
						log.Fatal("Error in decoding user validity:", err)
					}
					lib.Members = append(lib.Members, *member)
					return nil
				})
				if err != nil {
					return err
				}
				b = true
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println("DB Reading Error on library module")
	}
	if b { //used for out of index error handling
		if len(member.BooksBorrowed) == 5 {
			b = false
		}
	}
	return b, member
}

//Checks if the user wishing to return is eligible
func checkUserValidityReturn(name string, lib *Library, db *badger.DB) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			return b, member

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
			if k == name {
				err := item.Value(func(v []byte) error {
					// Create an encoder and send a value.
					enc := gob.NewDecoder(bytes.NewBuffer(v))
					err := enc.Decode(&member)
					if err != nil {
						log.Fatal("Error in decoding user validity return:", err)
					}
					lib.Members = append(lib.Members, *member)
					return nil
				})
				if err != nil {
					return err
				}
				b = true
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println("DB Reading Error on library module")
	}
	return b, member
}

//checks if username is unique as it acts as primary key
func checkUserNameValidity(username string, lib Library, db *badger.DB) bool {
	b := true
	for i := range lib.Members {
		if username == lib.Members[i].Name {
			b = false
			break
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
			if k == username {
				b = false
				break
			}

		}
		return nil
	}); err != nil {
		fmt.Println("DB Reading Error on library module")
	}

	return b
}

//checks if book borrowed and present in directory
func checkUserBookValidity(bname string, lib Library, member Member, db *badger.DB) (bool, *Books) {
	bfound := false //Denotes book validity
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow or return
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
					if err != nil {
						log.Fatal("Error in decoding user validity return:", err)
					}
					lib.BooksBorrowed = append(lib.BooksBorrowed, *bookFound)
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
