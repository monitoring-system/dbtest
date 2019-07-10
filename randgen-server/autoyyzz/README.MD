
# 使用方式：

```shell
sh scan.sh
```

之后会在conf目录下生成相应的yy zz文件

xuhuaiyu目录存放的是之前xuhuaiyu所写的yy zz文件，这个目录下的yy文件统一使用的是table-abs.zz这个zz文件

有一些算子因为模板无法通用暂时无法加入：

 - ADDDATE
 - DateAdd
 - SubDate
 - DateSub
 - AddTime
 - FromUnixTime
 - NOW
 - SubTime
 - UnixTimestamp
 - Convert
 - Substring
 - Substr
 - SubstringIndex
 - TimestampAdd
 - TimestampDiff
 - Conv
 - PI
 - POW
 - POWER
 - RAND
 - TRUNCATE
 - CONVERT_TZ
 - CURDATE
 - CURRENT_DATE
 - CURRENT_TIME
 - CURRENT_TIMESTAMP
 - CURTIME
 - DATE_FORMAT
 - EXTRACT
 - GET_FORMAT
 - LOCALTIME
 - LOCALTIMESTAMP
 - MAKEDATE
 - MAKETIME
 - PERIOD_ADD
 - SYSDATE
 - TIME_FORMAT
 - UTC_DATE
 - UTC_TIME
 - UTC_TIMESTAMP
 - ELT
 - EXPORT_SET
 - FIELD
 - FORMAT
 - FROM_BASE64
 - INSERT
 - INSTR
 - LEFT
 - RIGHT
 - LOADFILE
 - LOCATE
 - LPAD
 - MID
 - MAKE_SET
 - POSITION
 - REPEAT
 - REPLACE
 - Rpad
 - CHAR
 - FIND_IN_SET
 - CONNECTION_ID
 - CURRENT_USER
 - DATABASE
 - SCHEMA
 - FOUND_ROWS
 - LAST_INSERT_ID
 - USER
 - VERSION
 - BENCHMARK


# 如何添加算子

在operator_list文件中新增一行，内容为"算子 参数个数"，如果参数个数是任意个的话，可以直接写作'unbound'

# 单独使用模板替换功能：

```bash
python template.py SUM 2
```

其中sum就是你想要测试的算子，2是该算子可以接受的参数数目。生成的结果文件
为`conf/SUM.yy`和`conf/SUM.zz`

如果要测试的算子的参数个数是任意的，比如`COALESCE`，可以将参数数量设置为'unbound'：

```bash
python template.py COALESCE unbound
```

它将会生成2个参数的COALESCE，3个参数的COALESCE，5个参数的COALESCE
以及将全部field作为参数的COALESCE

上面的例子生成的`conf/COALESCE.yy`文件看起来如下：

```yacc
... operator ...(省略)

operator:
    COALESCE( _field , _field )
|   COALESCE( _field , _field , _field )
|   COALESCE( _field , _field , _field , _field , _field )
|   COALESCE( _fields )
```

默认使用的模板文件为`template/tmpl.yy`和`template/tmpl.zz`，还是以之前SUM算子为例：

 - `template/tmpl.zz`将会被直接拷贝过来作为`conf/SUM.zz`
 - 在`template/tmpl.yy`模板变量`$operator`将会被替换为`SUM( _field , _field )`，然后被拷贝到`conf/SUM.yy`
 
 
如果你不想要用默认的模板，可以使用模板变量`$operator`自己写一个模板，放知道当前的`template`目录，
然后通过`--tmpl`参数来指定模板。举个例子，`--tmpl a`表示使用模板`template/a.yy`和`template/a.zz`