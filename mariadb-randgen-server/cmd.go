package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/randgen"
	"github.com/monitoring-system/dbtest/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

var (
	mysqlPort   int
	mysqlHost   string
	mysqlUser   string
	mysqlPasswd string
	serverPort int
)

func main() {
	rootCmd := &cobra.Command{
		Use:"m-randgen-server",
		Short:"mariadb randgen server extension to generate sqls according to yy zz",
		Run: func(cmd *cobra.Command, args []string) {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?readTimeout=2s",
				mysqlUser,
				mysqlPasswd,
				mysqlHost,
				mysqlPort)

			dbiPrefix := fmt.Sprintf("dbi:mysql:host=%s:port=%d:user=%s:password=%s:database=",
				mysqlHost,
				mysqlPort,
				mysqlUser,
				mysqlPasswd)

			db, err := util.OpenDBWithRetry("mysql", dsn)
			if err != nil {
				logrus.Warn(err)
				os.Exit(1)
			}

			s := &server{dbiPrefix, db}

			err = s.StartServer()
			if err != nil {
				logrus.Printf("%v\n", err)
			}
		},
	}

	rootCmd.Flags().IntVar(&mysqlPort, "mysql-port", 3306, "mysql port")
	rootCmd.Flags().StringVar(&mysqlHost, "mysql-host", "localhost", "mysql host")
	rootCmd.Flags().StringVar(&mysqlUser, "mysql-user", "root", "mysql user")
	rootCmd.Flags().StringVar(&mysqlPasswd, "mysql-passwd", "", "mysql password")
	rootCmd.Flags().IntVar(&serverPort, "port", 9080, "server listen port")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}

type server struct {
	dbiPrefix string
	db *sql.DB
}

func (this *server) StartServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/loaddata", this.LoadData)
	err := http.ListenAndServe(":9080", mux)
	if err != nil {
		return err
	}
	return nil
}

// default zz file content if zz param is empty
const defaultZz = `$tables = {
    rows => [10, 30, 40, 60, 90],
};

$fields = {
         types => ['bigint', 'float', 'double', 'decimal(40, 20)', 'decimal(10, 4)', 'decimal(6, 3)', 'char(20)', 'varchar(20)'],
         #indexes => ['undef', 'key'],
         sign => ['signed', 'unsigned'],
         #charsets => ['utf8']
};

$data = {
         numbers => ['null', 'tinyint', 'smallint',
         '12.991', '1.009', '-9.1823',
         '-111.1212', '12.98731', '1.098781',
         '0.112345', '-0.987103', '-0.000000001',
         '0.00000001', '0.999999999', '-0.999999999',
         'decimal',
         ],
         strings => ['null', 'letter', 'english', 'string(15)'],
}`

func (this *server) LoadData(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	payload := LoadDataRequest{}
	json.NewDecoder(r.Body).Decode(&payload)

	loader := randgen.NewMariaGenerator(this.db, payload.DB, this.dbiPrefix)

	zzContent := payload.ZZ
	if strings.TrimSpace(zzContent) == "" {
		zzContent = defaultZz
	}

	sqls, err := loader.LoadData(payload.ZZ, payload.Yy, payload.DB, payload.Queries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write(MustJosnMarshal(&LoadDataResponse{sqls}))
}

func MustJosnMarshal(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

type LoadDataRequest struct {
	Yy      string `json:"yy"`
	ZZ      string `json:"zz"`
	DB      string `json:"db"`
	Queries int    `json:"queries"`
}

type LoadDataResponse struct {
	SQLs    []string `json:"sql"`
}