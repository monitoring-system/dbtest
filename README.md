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

# start dbtest server
```bash
go run main.go 
```

# call api to submit a new test
URL: 
```bash
POST 127.0.0.1:8080/tests
```
payload:
```json
{
  "testName":"rand1",
  "dataLoader":"name",
  "queryLoader":"query", 
 "yy":"query:\n    select ;\n\nselect:\n    SELECT coalesce\n    FROM _table\n    WHERE condition ;\n\ncoalesce:\n    COALESCE( _field , 0) | COALESCE( _field_list ) ;\n\ncondition:\n    _field IS NULL | _field = 1111 | _field = 'hello' ;",
 "zz":"$tables = {\n        rows =\u003e [0, 1, 10, 100],\n        partitions =\u003e [ undef , 'KEY (pk) PARTITIONS 2' ]\n};\n\n$fields = {\n        types =\u003e [ 'int', 'char', 'enum', 'set' ],\n        indexes =\u003e [undef, 'key' ],\n        null =\u003e [undef, 'not null'],\n        default =\u003e [undef, 'default null'],\n        sign =\u003e [undef, 'unsigned'],\n        charsets =\u003e ['utf8', 'latin1']\n};\n\n$data = {\n        numbers =\u003e [ 'digit', 'null', undef ],\n        strings =\u003e [ 'letter', 'english' ],\n        blobs =\u003e [ 'data' ],\n\ttemporals =\u003e ['date', 'year', 'null', undef ]\n}\n",
 "loop":1,
 "queries":1000, 
 "loopInterval": 30
}

```