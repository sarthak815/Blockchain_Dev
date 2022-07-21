package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

//Member Structure
type Member struct {
	Name          string `json:`
	Age           int
	BooksBorrowed []Books
}

//*********************STRUCT CONSTRUCTORS*********************
//Constructor initialising member struct
func (member *Member) Init(name string, age int) {
	member.Name = name
	member.Age = age

}
func checkMemberValidityApi(mname string, lib *Library, db *badger.DB) bool {
	fmt.Println(mname)
	bfound := false //Denotes book validity
	for i := range lib.Members {
		if lib.Members[i].Name == mname {
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
			if k == mname {
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
