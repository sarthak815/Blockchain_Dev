package main

import (
	"bytes"
	"fmt"
)

// The Vector type has unexported fields, which the package cannot access.
// We therefore write a BinaryMarshal/BinaryUnmarshal method pair to allow us
// to send and receive the type with the gob package. These interfaces are
// defined in the "encoding" package.
// We could equivalently use the locally defined GobEncode/GobDecoder
// interfaces.
type Vector struct {
	x, y, z int
}

func (v Vector) MarshalBinary() ([]byte, error) {
	// A simple encoding: plain text.
	var b bytes.Buffer
	fmt.Fprintln(&b, v.x, v.y, v.z)
	return b.Bytes(), nil
}

// UnmarshalBinary modifies the receiver so it must take a pointer receiver.
func (v *Vector) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &v.x, &v.y, &v.z)
	return err
	//}
	//
	//// This example transmits a value that implements the custom encoding and decoding methods.
	//func main() {
	//	db, err := badger.Open(badger.DefaultOptions("C:\\Users\\Sanjay\\OneDrive - Manipal Academy of Higher Education\\Documents\\GitHub\\Blockchain_Dev\\Task5-Library Management System(Persistent Version)\\tmp\\badger"))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	err = db.View(func(txn *badger.Txn) error {
	//		opts := badger.DefaultIteratorOptions
	//		opts.PrefetchSize = 10
	//		it := txn.NewIterator(opts)
	//		defer it.Close()
	//		for it.Rewind(); it.Valid(); it.Next() {
	//			item := it.Item()
	//			k := item.Key()
	//			err := item.Value(func(v []byte) error {
	//				fmt.Printf("key=%s, value=%s\n", k, v)
	//				return nil
	//			})
	//			if err != nil {
	//				return err
	//			}
	//		}
	//		return nil
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//err = db.Update(func(txn *badger.Txn) error {
	//	// Start a writable transaction.
	//	txn = db.NewTransaction(true)
	//	defer txn.Discard()
	//	e := badger.NewEntry([]byte("Hello"), []byte("Yesssssssss"))
	//	err := txn.SetEntry(e)
	//	return err
	//	err = txn.Commit()
	//	return err
	//})
	//fmt.Println("Inserted Members")
	//defer db.Close()
}
