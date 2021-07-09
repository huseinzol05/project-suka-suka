package main

import (
    "os"
    "fmt"
    "log"
    "sync"
    "math"
    "sort"
    "strconv"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/mkmik/argsort"
    "github.com/julienschmidt/httprouter"
    "github.com/syndtr/goleveldb/leveldb"
)

const Earth_Radius_KM = 6372.8

type DB struct {
    db *leveldb.DB
    mu sync.Mutex
}

func Rdist(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	sin_0 := math.Sin(0.5 * (lat2 - lat1))
	sin_1 := math.Sin(0.5 * (long2 - long1))
	return (sin_0 * sin_0 + math.Cos(lat1) * math.Cos(lat2) * sin_1 * sin_1)
}

func Rdist_to_dist(rdist float64) float64 {
	return 2 * math.Asin(math.Sqrt(rdist))
}

func Dist(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	return 2 * math.Asin(math.Sqrt(Rdist(lat1, lon1, lat2, lon2)))
}

func Radian(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}

func NodeBallTree(db *leveldb.DB, db_index *leveldb.DB, keys *[]string, label string, g *sync.WaitGroup){
	var result map[string]interface{}
	lat := 0.0
	long := 0.0
	min_lat := 9999.99
	min_long := 9999.99
	max_lat := -9999.99
	max_long := -9999.99
	var lats []float64
	var longs []float64
	for i := 0; i < len(*keys); i++ {
		value, _ := db.Get([]byte((*keys)[i]), nil)
		json.Unmarshal(value, &result)
		lat_ := result["lat"].(float64)
		long_ := result["long"].(float64)
		lat += lat_
		long += long_
		if lat_ < min_lat {
			min_lat = lat_
		}
		if long_ < min_long {
			min_long = long_
		}
		if lat_ > max_lat {
			max_lat = lat_
		}
		if long_ > max_long {
			max_long = long_
		}
		lats = append(lats, lat_)
		longs = append(longs, long_)
	}
	order_lat := argsort.Sort(sort.Float64Slice(lats))
	order_long := argsort.Sort(sort.Float64Slice(longs))
	lat /= float64(len(*keys))
	long /= float64(len(*keys))
	radius := -9999.99

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		json.Unmarshal(value, &result)
		lat_ := result["lat"].(float64)
		long_ := result["long"].(float64)
		rdist := Rdist(lat_, long_, lat, long)
		if rdist > radius {
			radius = rdist
		}
	}
	iter.Release()

	radius = Rdist_to_dist(radius)
	minus_lat := max_lat - min_lat
	minus_long := max_long - min_long
	order := &order_lat
	if minus_long > minus_lat {
		order = &order_long
	}

	sorted_keys := make([]string, len(*keys))
	for i := 0; i < len(*keys); i++ {
		sorted_keys[i] = (*keys)[(*order)[i]]
	}

	var data map[string]interface{}
	got_children := true

	if len(order_lat) == 1 {
		value, _ := db.Get([]byte((*keys)[(*order)[0]]), nil)
		json.Unmarshal(value, &data)
		got_children = false
	}
	fmt.Println(label, len(*keys), len(order_lat), got_children)
	left_key := label + "-left"
	right_key := label + "-right" 
	output := map[string]interface{}{
        "left": left_key,
		"right": right_key,
		"lat": lat,
		"long": long,
		"radius": radius,
		"data": data,
	}
	b, _ := json.Marshal(output)
	db_index.Put([]byte(label), b, nil)

	if got_children {
		median := len(sorted_keys) / 2
		left_sorted_keys := sorted_keys[:median]
		right_sorted_keys := sorted_keys[median:]
		g.Add(1)
		go NodeBallTree(db, db_index, &left_sorted_keys, left_key, g)
		g.Add(1)
		go NodeBallTree(db, db_index, &right_sorted_keys, right_key, g)
	}

	defer g.Done()
}

func BallTree(db *leveldb.DB, db_index *leveldb.DB, g *sync.WaitGroup){
	var result map[string]interface{}
	iter := db.NewIterator(nil, nil)
	lat := 0.0
	long := 0.0
	count := 0
	min_lat := 9999.99
	min_long := 9999.99
	max_lat := -9999.99
	max_long := -9999.99
	var lats []float64
	var longs []float64
	var keys []string
	for iter.Next() {
		key := string(iter.Key())
		value := iter.Value()
		json.Unmarshal(value, &result)
		fmt.Println(count, key)
		lat_ := result["lat"].(float64)
		long_ := result["long"].(float64)
		lat += lat_
		long += long_
		if lat_ < min_lat {
			min_lat = lat_
		}
		if long_ < min_long {
			min_long = long_
		}
		if lat_ > max_lat {
			max_lat = lat_
		}
		if long_ > max_long {
			max_long = long_
		}
		count += 1
		lats = append(lats, lat_)
		longs = append(longs, long_)
		keys = append(keys, key)
	}
	order_lat := argsort.Sort(sort.Float64Slice(lats))
	order_long := argsort.Sort(sort.Float64Slice(longs))
	iter.Release()
	lat /= float64(count)
	long /= float64(count)
	radius := -9999.99

	iter = db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		json.Unmarshal(value, &result)
		lat_ := result["lat"].(float64)
		long_ := result["long"].(float64)
		rdist := Rdist(lat_, long_, lat, long)
		if rdist > radius {
			radius = rdist
		}
	}
	iter.Release()

	radius = Rdist_to_dist(radius)
	minus_lat := max_lat - min_lat
	minus_long := max_long - min_long
	order := &order_lat
	if minus_long > minus_lat {
		order = &order_long
	}

	sorted_keys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		sorted_keys[i] = keys[(*order)[i]]
	}

	var data map[string]interface{}
	if len(order_lat) == 1 {
		value, _ := db.Get([]byte(keys[(*order)[0]]), nil)
		json.Unmarshal(value, &data)
	}
	left_key := "root-left"
	right_key := "root-right"
	output := map[string]interface{}{
        "left": left_key,
		"right": right_key,
		"lat": lat,
		"long": long,
		"radius": radius,
		"data": data,
	}
	b, _ := json.Marshal(output)
	db_index.Put([]byte("root"), b, nil)

	median := len(sorted_keys) / 2
	left_sorted_keys := sorted_keys[:median]
	right_sorted_keys := sorted_keys[median:]
	g.Add(1)
	go NodeBallTree(db, db_index, &left_sorted_keys, left_key, g)
	g.Add(1)
	go NodeBallTree(db, db_index, &right_sorted_keys, right_key, g)

	defer g.Done()
}

	
func QueryTree(db_index *leveldb.DB, label string, lat float64, long float64, 
	min_distance float64, max_distance float64, g *sync.WaitGroup) []map[string]interface{} {
	var results []map[string]interface{}
	var result map[string]interface{}
	var data map[string]interface{}
	value, _ := db_index.Get([]byte(label), nil)
	json.Unmarshal(value, &result)
	data_ := result["data"]
	if data_ == nil {
		distance_ := Dist(result["lat"].(float64), result["long"].(float64), lat, long)
		rad := result["radius"].(float64)
		if (distance_ >= min_distance - rad) && (distance_ <= max_distance + rad){
			g.Add(1)
			results_left := QueryTree(db_index, label + "-left", lat, long, min_distance, max_distance, g)
			g.Add(1)
			results_right := QueryTree(db_index, label + "-right", lat, long, min_distance, max_distance, g)
			results = append(results, results_left...)
			results = append(results, results_right...)
		}
	} else {
		data = result["data"].(map[string]interface{})
		distance_ := Dist(data["lat"].(float64), data["long"].(float64), lat, long)
		if (distance_ >= min_distance) && (distance_ <= max_distance) {
			data["distance"] = distance_ * Earth_Radius_KM
			results = append(results, data)
		}
	}
	defer g.Done()
	return results
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
    g := &sync.WaitGroup{}
    g.Add(1)
    db, _ := leveldb.OpenFile("db/" + index, nil)
	db_index, _ := leveldb.OpenFile("db-index/" + index, nil)
	BallTree(db, db_index, g)
	g.Wait()
    fmt.Fprintf(w, "{\"status\": \"%s success.\"}", index)
    defer db.Close()
}

func Query(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
    min_distance, ok := r.URL.Query()["min_distance"]
    if !ok || len(min_distance[0]) < 1 {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"`min_distance` is missing.\"}")
        return
    }

    max_distance, ok := r.URL.Query()["max_distance"]
    if !ok || len(max_distance[0]) < 1 {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"`max_distance` is missing.\"}")
        return
    }

    lat, ok := r.URL.Query()["lat"]
    if !ok || len(lat[0]) < 1 {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"`lat` is missing.\"}")
        return
    }

    long, ok := r.URL.Query()["long"]
    if !ok || len(long[0]) < 1 {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"`long` is missing.\"}")
        return
    }

    min_radius_km, _ := strconv.ParseFloat(min_distance[0], 64)
    min_radius_km /= Earth_Radius_KM
    max_radius_km, _ := strconv.ParseFloat(max_distance[0], 64) 
    max_radius_km /= Earth_Radius_KM
    lat_float, _ := strconv.ParseFloat(lat[0], 64)
    long_float, _ := strconv.ParseFloat(long[0], 64)

    index := ps.ByName("index")
    directory := "db-index/" + index
    if _, err := os.Stat(directory); os.IsNotExist(err) {
        w.WriteHeader(400)
        fmt.Fprintf(w, "{\"error\": \"%s not exist, not able to query.\"}", index)
        return
    }
    lat_float = Radian(lat_float)
    long_float = Radian(long_float)
    g := &sync.WaitGroup{}
	g.Add(1)
	db_index, _ := leveldb.OpenFile("db-index/" + index, nil)
	results := QueryTree(db_index, "root", lat_float, long_float, min_radius_km, max_radius_km, g)
	g.Wait()
	distances := make([]float64, len(results))
	for i := 0; i < len(results); i++ {
		distances[i] = results[i]["distance"].(float64)
	}
	order := argsort.Sort(sort.Float64Slice(distances))
	sorted := make([]map[string]interface{}, len(results))
	for i := 0; i < len(results); i++ {
		sorted[i] = results[order[i]]
    }
    json_data, _ := json.Marshal(sorted)
    fmt.Fprintf(w, string(json_data))
    defer db_index.Close()
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
