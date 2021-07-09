# LatLongDB

Haversine Ball Tree index DB to query and aggregate based on distance (KM).

## how-to

1. Run main,

```bash
go get -u github.com/julienschmidt/httprouter \
github.com/syndtr/goleveldb/leveldb \
github.com/google/uuid \
github.com/mkmik/argsort
go run main.go
```

## tests

1. Run [insert.ipynb](insert.ipynb).

```python
index = 'test1'
samples = 200
batch_size = 10000
```

2. Run [balltree_index.go](balltree_index.go),

```bash
go run balltree_index.go
```

You should see indexing branches,

```text
root-right-right-right-right 13 13 true
root-right-right-right-right-left 6 6 true
root-right-right-right-right-left-left 3 3 true
root-right-right-right-right-left-left-left 1 1 false
root-right-right-right-right-left-left-right 2 2 true
```

3. Run [balltree_query.go](balltree_query.go),

```bash
go run balltree_query.go
```

I hard coded the variables,

```golang
index := "test1"
radius_km := 10.0
lat := Radian(2.950815010581982)
long := Radian(101.62843052319319)
```

Output,
```text
map[distance:9.411105869317856 lat:0.05253125325274705 long:1.7726908484321657 no:191] map[distance:9.441697545710548 lat:0.050183971479504035 long:1.7744293264197177 no:104] map[distance:9.447640848458185 lat:0.05131428468218985 long:1.7722781470727473 no:0] map[distance:9.479414939814848 lat:0.05001680244903595 long:1.773842803987455 no:137] map[distance:9.504295807966539 lat:0.05296700407868428 long:1.7734740468784655 no:25] map[distance:9.547934571092823 lat:0.050021309131340265 long:1.7735182204361621 no:86] map[distance:9.613456584741758 lat:0.05285792894644737 long:1.773089868491407 no:173] map[distance:9.642893157965883 lat:0.04998859245945505 long:1.7737802737982844 no:50] map[distance:9.643623272029211 lat:0.051786477335728956 long:1.7752388698647812 no:109] map[distance:9.666685260254345 lat:0.04998548620335271 long:1.7736979825702508 no:178] map[distance:9.680463281473921 lat:0.05215014377689979 long:1.7751261184089222 no:132] map[distance:9.71408350210781 lat:0.05115616232494017 long:1.7752373684417275 no:28] map[distance:9.798368334210416 lat:0.05024849858277584 long:1.7728584257721378 no:182] map[distance:9.81709056164358 lat:0.05030970527208331 long:1.7727733381994215 no:120]]
```

## Why

1. Practice my Golang and geospatial indexing.