package main

import (
	"fmt"
	"encoding/binary"
    "github.com/syndtr/goleveldb/leveldb"
)

func main() {
	
	db, _ := leveldb.OpenFile("db/test1", nil)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 0)
	// _ = db.Put([]byte("foo"), buf, nil)
	data, _ := db.Get([]byte("count"), nil)
	num := binary.BigEndian.Uint64(data)
	fmt.Printf("%d", num)

	buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 0)
	data, _ = db.Get([]byte(buf), nil)
	fmt.Printf("%s", data)
	// fmt.Printf("%s", data)
	// var num uint64 = 0
	// if err == nil {
	// 	num = binary.BigEndian.Uint64(data)
	// }
	// fmt.Printf("%d", num)
	// fmt.Printf("%s %s", []byte("foo"), buf)

}

// import (
// 	"fmt"
// )

// func main (){
// 	a := 2
// 	a = 4
// 	fmt.Printf("%d", a)
// }