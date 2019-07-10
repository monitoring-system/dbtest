# dbtest
database test framework


# How to Use it

1. build rangden docker image
    ```bash
    cd randgen-server
    GOOS=linux go build -o ./randgen-server  ./main.go
    docker build -f Dockerfile -t "randgen-server:latest" .
    ```

2. start docker container
    ```bash
     docker run  -p 9080:9080 randgen-server 
    ```

3. build the dbtest
    ```bash
    go build main.go -o ./dbtest
    ```

4. start dbtest server
    ```bash
    ./dbtest start --standard-db=root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local --test-db=root:@tcp(127.0.0.1:4000)/?charset=utf8&parseTime=True&loc=Local
    ```

5. submit randgen yy/zz file
    ```bash
    # single yy zz file
    ./dbtest add --yy=randgen/examples/example.yy --zz=randgen/examples/example.zz
    # or  add all yy zz file in a directory
    ./dbtest add --loadpath=randgen/examples/
    ```
6.  watch the test status
    ```bash
    ./dbtest watch
    ```
    
7. check test log and data, logs and data can be found in results directory
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
    
# Run with docker 

1. build binary
    ```bash
    make 
    ```
2. build docker image
   ```bash
   cd docker
   sh pack.sh
   ```
3. start docker    
    ```bash
     docker run -p 8080:8080 framework-dbtest /root/dbtest start --standard-db="root:@tcp(172.16.4.65:31175)/?charset=utf8&parseTime=True&loc=Local" --test-db="root:@tcp(172.16.4.65:30453)/?charset=utf8&parseTime=True&loc=Local"
    ```
    
    
# Generate yy zz with template

[doc](randgen-server/autoyyzz)

# Add filter

build your go plugin by `go build -buildmode=plugin`, then put the so file in 
`plugin-filters` director.

Your plugin must inclue a method called `Filter`, then must have signature as 

```go
func(errMsg string, source string) bool
```

or 

```go
func(vInTiDB interface{}, vInMySQL interface{}, colType *sql.ColumnType) bool
```

[一个示例](filter/filter.go) 中的filterNumberPrecision方法