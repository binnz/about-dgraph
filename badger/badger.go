package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions(""))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	txn := db.NewTransaction(true)
	defer txn.Discard()
	txn.Set([]byte("bin"), []byte("21"))
	txn.Commit()
	txnr := db.NewTransaction(false)
	item, err := txnr.Get([]byte("bin"))
	// val, err := item.Value()
	fmt.Println("haha")
	res := make([]byte, 100)
	fmt.Println(item.ValueCopy(res))
	fmt.Println(string(res))
	fmt.Println("heh")
}
