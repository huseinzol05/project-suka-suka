package main

import (
	"fmt"
	"encoding/binary"
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
	
	db, _ := leveldb.OpenFile("db/b", nil)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 12)
	_ = db.Put([]byte("foo"), buf, nil)
	data, err := db.Get([]byte("foo1"), nil)
	var num uint64 = 0
	if err == nil {
		num = binary.BigEndian.Uint64(data)
	}
	fmt.Printf("%d", num)
	fmt.Printf("%s %s", []byte("foo"), buf)

}