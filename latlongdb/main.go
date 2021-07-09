package main

import (
    "os"
    "fmt"
    "log"
    "sync"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/julienschmidt/httprouter"
    "github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
    db *leveldb.DB
    mu sync.Mutex
}

type BallTree struct {
    Left *BallTree
    Right *BallTree
    Loc [2]float64
    Radius float64
    Data [1]map[string]interface{}
}

func Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "LatLongDB!\n")
}

func Insert(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    var results []map[string]interface{}
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal([]byte(body), &results)
    if err != nil {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"body not a JSON.\"}")
        return
    }
    index := ps.ByName("index")
    db, _ := leveldb.OpenFile("db/" + index, nil)
    dbm := DB{db: db}
    dbm.mu.Lock()
    var success_row []int
    var fail_row []int
    var fail_reasons []string
    for no, result := range results {
        s := "success"
        id := uuid.New().String()
        success := true
        _, ok_lat := result["lat"]
        _, ok_long := result["long"]
        _, ok_id := result["id"]
        if !ok_long || !ok_lat {
            s = "key `lat` or `long` not exist"
            success = false
        }
        if ok_id {
            id = fmt.Sprint(result["id"])
        }
        if success {
            success_row = append(success_row, no)
            b, _ := json.Marshal(result)
            dbm.db.Put([]byte(id), b, nil)
        } else {
            fail_row = append(fail_row, no)
            fail_reasons = append(fail_reasons, s)
        }
        
    }
    output := map[string]interface{}{
        "success_row": success_row,
        "fail_row": fail_row,
        "fail_reasons": fail_reasons,
    }
    json_data, _ := json.Marshal(output)
    fmt.Fprintf(w, string(json_data))
    dbm.mu.Unlock()
    defer dbm.db.Close()
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    index := ps.ByName("index")
    directory := "db/" + index
    if _, err := os.Stat(directory); os.IsNotExist(err) {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"%s not exist, not able to index.\"}", index)
        return
    }
    db, _ := leveldb.OpenFile(directory, nil)
    data, _ := db.Get([]byte("key"), nil)
    fmt.Fprintf(w, "index, %s!\n", data)
    defer db.Close()
}

func Query(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
    index := ps.ByName("index")
    fmt.Fprintf(w, "aggregate, %s", index)
}

func Aggregate(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
    index := ps.ByName("index")
    fmt.Fprintf(w, "aggregate, %s", index)
}

func main() {
    router := httprouter.New()
    router.GET("/", Root)
    router.POST("/:index/insert", Insert)
    router.GET("/:index/index", Index)
    router.GET("/:index/query", Query)
    router.GET("/:index/aggregate", Aggregate)
    log.Fatal(http.ListenAndServe(":8080", router))
}
