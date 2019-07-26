module github.com/monitoring-system/dbtest

go 1.12

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/a8m/rql v1.1.0
	github.com/buger/goterm v0.0.0-20181115115552-c206103e1f37
	github.com/dqinyuan/go-mysqldump v0.2.5
	github.com/fatih/color v1.7.0
	github.com/gin-contrib/pprof v1.2.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/jinzhu/gorm v1.9.9
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/pingcap/errors v0.11.0
	github.com/pingcap/log v0.0.0-20190307075452-bd41d9273596
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.4.1
	github.com/sergi/go-diff v1.0.0
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.3.0
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	go.uber.org/zap v1.10.0
	gopkg.in/go-playground/assert.v1 v1.2.1
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
