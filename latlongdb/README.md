# LatLongDB

Haversine Ball Tree index DB to query based on distance (KM).

## how-to

1. Run main,

```bash
go get -u github.com/julienschmidt/httprouter \
github.com/syndtr/goleveldb/leveldb \
github.com/google/uuid \
github.com/mkmik/argsort
go run main.go
```

### Insert

```bash
curl --location --request POST 'http://localhost:8080/test1/insert' \
--header 'Content-Type: application/json' \
--data-raw '[{
    "lat": 1, "long": 1, "data": 1
}]'
```

For bulk automatic inserts, check [insert.ipynb](insert.ipynb).

### Index

Before able to query, we need to index the data into Haversine Ball Tree first,

```bash
curl --request GET 'http://localhost:8080/test1/index'
```

```text
{"status": "test1 success."}
```

### Query

I want to query between 2KM and 6.4KM given a lat long,

```bash
time curl --request GET 'http://localhost:8080/test1/query?min_distance=2.0&max_distance=6.4&lat=2.950815010581982&long=101.62843052319319'
```

```text
time curl --request GET 'http://localhost:8080/test1/query?min_distance=2.0&max_distance=6.4&lat=2.950815010581982&long=101.62843052319319'
[{"distance":2.7435444445887867,"lat":0.051198664647091295,"long":1.774057179420072,"no":64},{"distance":2.908191238826303,"lat":0.05152282827772343,"long":1.7732942780718577,"no":66},{"distance":2.9317181243178023,"lat":0.05193458662198576,"long":1.7739059018276264,"no":79},{"distance":3.507168524801305,"lat":0.05111951985242149,"long":1.7733539637739804,"no":6},{"distance":3.7777900075606414,"lat":0.05186169237965742,"long":1.7742221286769562,"no":22},{"distance":4.123977332125158,"lat":0.051004354314852604,"long":1.7741656076241954,"no":90},{"distance":4.1841168418478,"lat":0.051892999414839665,"long":1.7732230026785,"no":60},{"distance":4.285398680702833,"lat":0.05106455168399484,"long":1.7732388569635495,"no":42},{"distance":4.468724842160898,"lat":0.051134292761990933,"long":1.7743489349145352,"no":2},{"distance":4.514926771470817,"lat":0.05212757078380303,"long":1.774082660067285,"no":76},{"distance":4.579484714434945,"lat":0.050788123678384846,"long":1.7738378296203465,"no":93},{"distance":4.698476735819023,"lat":0.050765947492623885,"long":1.7738019987425815,"no":13},{"distance":4.72191272415612,"lat":0.05129110975924112,"long":1.7744621345361506,"no":95},{"distance":5.068992707067398,"lat":0.052225952818291586,"long":1.7734220227025197,"no":53},{"distance":5.192150412617559,"lat":0.052303936773922394,"long":1.7736098597669325,"no":46},{"distance":5.6006044241166375,"lat":0.052022697865773015,"long":1.7744592265592674,"no":54},{"distance":5.724694769842948,"lat":0.05098346152590945,"long":1.7744856160475193,"no":74},{"distance":5.78003545647078,"lat":0.052390356234102604,"long":1.7735703558401288,"no":75},{"distance":6.007678145874556,"lat":0.05066778662816621,"long":1.7733100136898774,"no":72},{"distance":6.310308867273926,"lat":0.050554845585368484,"long":1.7734597481619931,"no":45}]
real	0m0.098s
user	0m0.004s
sys	0m0.009s
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