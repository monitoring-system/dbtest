package main

import (
	"fmt"
	"github.com/monitoring-system/dbtest/randgen"
	"github.com/monitoring-system/dbtest/util"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"os"
)

type config struct {
	mysqlPort   int
	mysqlHost   string
	mysqlUser   string
	mysqlPasswd string
	serverPort  int
}

var rConfig = &config{}

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

func main() {
	rootCmd := &cobra.Command{
		Use:   "randgen-server",
		Short: "a server to generate sqls according to yy zz",
		Run: func(cmd *cobra.Command, args []string) {

			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?readTimeout=2s",
				rConfig.mysqlUser,
				rConfig.mysqlPasswd,
				rConfig.mysqlHost,
				rConfig.mysqlPort)

			dbiPrefix := fmt.Sprintf("dbi:mysql:host=%s:port=%d:user=%s:password=%s:database=",
				rConfig.mysqlHost,
				rConfig.mysqlPort,
				rConfig.mysqlUser,
				rConfig.mysqlPasswd)

			db, err := util.OpenDBWithRetry("mysql", dsn)
			if err != nil {
				log.Fatalf("connect mysql error %v\n", err)
			}

			s := &randgen.Server{
				db,
				dbiPrefix,
				defaultZz,
			}

			s.Listen(rConfig.serverPort)
		},
	}

	rootCmd.Flags().IntVar(&rConfig.mysqlPort, "mysql-port", 3306, "mysql port");
	rootCmd.Flags().StringVar(&rConfig.mysqlHost, "mysql-host", "localhost", "mysql host")
	rootCmd.Flags().StringVar(&rConfig.mysqlUser, "mysql-user", "root", "mysql user")
	rootCmd.Flags().StringVar(&rConfig.mysqlPasswd, "mysql-passwd", "", "mysql password")
	rootCmd.Flags().IntVar(&rConfig.serverPort, "port", 9080, "server listen port")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}
