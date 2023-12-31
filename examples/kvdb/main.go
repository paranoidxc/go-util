package main

import (
	"fmt"
	"time"

	"github.com/paranoidxc/go-util/kvdb"
)

func main() {

	db := kvdb.Open(kvdb.WithDBFileName("hello.data"))
	defer db.Close()
	db.Printf()

	db.Put("hello", []byte("world"+" "+time.Now().Format("2006-01-02 15:04:05")))
	value, ok := db.Get("hello")
	if ok != nil {
		fmt.Println("Key not found")
	} else {
		fmt.Printf("get val:%s from DB by k:hello\n", string(value))
	}

	db.Put("test", []byte("test will delete"))
	value, ok = db.Get("test")
	if ok != nil {
		fmt.Println("Key not found")
	} else {
		fmt.Printf("get val:%s from DB by k:test\n", string(value))
	}

	fmt.Println("db delete val by key:test")
	if err := db.Delete("test"); err != nil {
		fmt.Println("db.Delete(test) err", err)
	}

	db.Printf()

	db.Tx(func(tx *kvdb.Tx) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()

		val, _ := tx.Get("hello")
		fmt.Println("tx get [hello] val:", string(val))
		tx.Put("tx_test", []byte("val from tx"))
		val, _ = tx.Get("tx_test")
		fmt.Println("tx get [tx_text] val:", string(val))
	})
}
