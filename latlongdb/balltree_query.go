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

const Earth_Radius_KM = 6372.8

func Rdist(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	sin_0 := math.Sin(0.5 * (lat2 - lat1))
	sin_1 := math.Sin(0.5 * (long2 - long1))
	return (sin_0 * sin_0 + math.Cos(lat1) * math.Cos(lat2) * sin_1 * sin_1)
}

func Dist(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	return 2 * math.Asin(math.Sqrt(Rdist(lat1, lon1, lat2, lon2)))
}

func Radian(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
	
func Query(db_index *leveldb.DB, label string, lat float64, long float64, 
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
			results_left := Query(db_index, label + "-left", lat, long, min_distance, max_distance, g)
			g.Add(1)
			results_right := Query(db_index, label + "-right", lat, long, min_distance, max_distance, g)
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

func main() {
	g := &sync.WaitGroup{}
	g.Add(1)
	index := "test1"
	min_radius_km := 3.0 / Earth_Radius_KM
	max_radius_km := 10.0 / Earth_Radius_KM
	lat := Radian(2.950815010581982)
	long := Radian(101.62843052319319)
	db_index, _ := leveldb.OpenFile("db-index/" + index, nil)
	results := Query(db_index, "root", lat, long, min_radius_km, max_radius_km, g)
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
	fmt.Println(sorted)
	
}