package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"encoding/json"
	"github.com/mkmik/argsort"
    "github.com/syndtr/goleveldb/leveldb"
)

func Rdist(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	sin_0 := math.Sin(0.5 * (lat2 - lat1))
	sin_1 := math.Sin(0.5 * (long2 - long1))
	return (sin_0 * sin_0 + math.Cos(lat1) * math.Cos(lat2) * sin_1 * sin_1)
}

func Rdist_to_dist(rdist float64) float64 {
	return 2 * math.Asin(math.Sqrt(rdist))
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


func main() {
	index := "test1"
	g := &sync.WaitGroup{}
	g.Add(1)
	db, _ := leveldb.OpenFile("db/" + index, nil)
	db_index, _ := leveldb.OpenFile("db-index/" + index, nil)
	BallTree(db, db_index, g)
	g.Wait()
}