package main

import (
	"fmt"
	"encoding/json"
	// "encoding/binary"
    // "github.com/syndtr/goleveldb/leveldb"
)

func main() {
	
	// db, _ := leveldb.OpenFile("db/test1", nil)
	// buf := make([]byte, 8)
	// binary.BigEndian.PutUint64(buf, 0)
	// // _ = db.Put([]byte("foo"), buf, nil)
	// data, _ := db.Get([]byte("count"), nil)
	// num := binary.BigEndian.Uint64(data)
	// fmt.Printf("%d", num)

	// buf = make([]byte, 8)
	// binary.BigEndian.PutUint64(buf, 0)
	// data, _ = db.Get([]byte(buf), nil)
	// fmt.Printf("%s", data)
	// fmt.Printf("%s", data)
	// var num uint64 = 0
	// if err == nil {
	// 	num = binary.BigEndian.Uint64(data)
	// }
	// fmt.Printf("%d", num)
	// fmt.Printf("%s %s", []byte("foo"), buf)
	success_row := []int{1,2,3}
	fail_row := []int{2,3,4}
	fail_reasons := []string{"a", "b"}
	// map[string][]interface{}{
	output := map[string]interface{}{
        "success_row": success_row,
        "fail_row": fail_row,
        "fail_reasons": fail_reasons,
	}
	fmt.Printf("%#v\n",output)
	fmt.Println(output, output["fail_row"])
	json_data, err := json.Marshal(output)
	fmt.Println(err)
    fmt.Println(string(json_data))
}

// import (
// 	"fmt"
// )

// func main (){
// 	a := 2
// 	a = 4
// 	fmt.Printf("%d", a)
// }