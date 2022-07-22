package codeModules

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
)

//Library structure
type Library struct {
	//BooksAvailable stores all the books entered in the application
	BooksAvailable []Books
	//Members stores all member details entered in the application
	Members []Member
}

//Checks if the user wishing to borrow is registered and eligible
func CheckUserValidity(name string, lib *Library, db *badger.DB) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	//checks for the user in libraru cache
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			if len(member.BooksBorrowed) >= 5 {
				b = false
			}
			return b, member
		}
	}
	//in case not found in cache searches in badgerDB
	if !b {
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
						//adds user to cache if found in db
						lib.Members = append(lib.Members, *member)
						//obtains memory address of member stored in cache
						for i := range lib.Members {
							if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
								member = &lib.Members[i]
								b = true
								if len(member.BooksBorrowed) >= 5 {
									b = false
								}
								break
							}
						}
						return nil
					})
					if err != nil {
						return err
					}

				}

			}
			return nil
		}); err != nil {
			fmt.Println("DB Reading Error on library module")
		}
	}

	return b, member
}

//Checks if the user wishing to return is eligible
func CheckUserValidityReturn(name string, lib *Library, db *badger.DB) (bool, *Member) {
	b := false         //Denotes user validity
	var member *Member //Used to return member object of user that wishes to borrow or return
	//checks for user object in cache memory
	for i := range lib.Members {
		if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
			member = &lib.Members[i]
			b = true
			return b, member

		}
	}
	//checks in db if user data is not in cache
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
					//searches for the address to the member in library after it has been fetched from the db
					for i := range lib.Members {
						if lib.Members[i].Name == name { // checks validity by name, uses name as primary key
							member = &lib.Members[i]

						}
					}
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
func CheckUserNameValidity(username string, lib Library, db *badger.DB) bool {
	b := true
	//searches for user presence in the cache memory
	for i := range lib.Members {
		if username == lib.Members[i].Name {
			b = false
			break
		}
	}
	//if user is not found in cache it looks fo it in the db
	if b {
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
	}

	return b
}

//checks if book borrowed and present in directory
func CheckUserBookValidity(bname string, lib Library, member Member, db *badger.DB) (bool, *Books) {
	bfound := false //Denotes book availability
	borrowed := false
	var bookFound *Books //Used to return book object that user wishes to borrow or return
	//searches for the book in cache memory
	for i := range lib.BooksAvailable {
		if lib.BooksAvailable[i].B_Name == bname {
			bookFound = &lib.BooksAvailable[i]
			bfound = true
			return bfound, bookFound
		}
	}
	//checks in db if book is not found in the cache memory
	//uses a readonly db function of badgerDB
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

					// Create an decoder to decode the received value from the db
					enc := gob.NewDecoder(bytes.NewBuffer(v))
					err := enc.Decode(&bookFound)
					if err != nil {
						log.Fatal("Error in decoding user validity return:", err)
					}
					lib.BooksAvailable = append(lib.BooksAvailable, *bookFound)
					//obtains address to the book in cache memory
					for i := range lib.BooksAvailable {
						if lib.BooksAvailable[i].B_Name == bname {
							bookFound = &lib.BooksAvailable[i]

						}
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
	//used for out of index error handling
	if bfound {
		// checks if user has borrowed same book
		for i := range member.BooksBorrowed {
			if member.BooksBorrowed[i].B_Name == bookFound.B_Name {
				borrowed = true
			}
		}
	}
	//sets book status to false in case the user has not borrowed the book
	if !borrowed {
		bfound = false
	}
	return bfound, bookFound
}
