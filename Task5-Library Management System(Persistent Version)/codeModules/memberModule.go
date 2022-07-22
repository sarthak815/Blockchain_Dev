package codeModules

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

//Member Structure
type Member struct {
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	BooksBorrowed []Books `json:"books"`
}

//Borrower struct used to meet API endpoints to borrow or return a book
type Borrower struct {
	//Name of user wishing to borrow or return a book
	Name string `json:"name"`
	//Book type user wishes to borrow or return
	BookT BookType `json:"type"`
	//Book name user wishes to borrow or return
	BookN string `json:"bookname"`
}

//Constructor initialising member struct
func (member *Member) Init(name string, age int) {
	member.Name = name
	member.Age = age

}

//CheckMemberValidityApi checks if the member details entered are valid and not conflicting with other user details
func CheckMemberValidityApi(mname string, lib *Library, db *badger.DB) bool {
	fmt.Println(mname)
	bfound := false //Denotes book validity
	//checks if new member's requested name is already being used in cache memory
	for i := range lib.Members {
		if lib.Members[i].Name == mname {
			bfound = true
			return bfound
		}
	}
	//in case new member's requested name is not used in cache, it verifies with database as well
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
