# dbtest
database test framework


# build rangden docker image
```bash
cd randgen
make
docker build -f Dockerfile -t "randgen-server:latest" .
```

# start docker container
```bash
 docker run  -p 9080:9080 randgen-server 
```

# build the dbtest
```bash
go build main.go -o ./dbtest
```

# start dbtest server
```bash
./dbtest start --standard-db=root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local --test-db=root:@tcp(127.0.0.1:4000)/?charset=utf8&parseTime=True&loc=Local
```

# submit a new test
```bash
./dbtest add --yy=randgen/examples/example.yy --zz=randgen/examples/example.zz
```

# add all yy zz file in a directory
```bash
./dbtest add --loadpath=randgen/examples/
```
# watch the test status
```bash
./dbtest watch
```
# test log and data

logs and data can be found in results directory
```
results/
└── logs # base dir
    ├── 1
    │   ├── 1.log  # test logs
    │   ├── 1.query  # all exectued queries
    │   └── 1.sql  # all data that is inserted into db
    └── 2
        ├── 1.log
        ├── 1.query
        └── 1.sql

```