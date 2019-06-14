package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/monitoring-system/dbtest/api"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/monitoring-system/dbtest/config"
	"github.com/monitoring-system/dbtest/db"
	"github.com/monitoring-system/dbtest/executor"
	"github.com/monitoring-system/dbtest/util"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var StartCmd = &cobra.Command{
	Use:   "StartServer [options ]",
	Short: "StartServer the db test web server",
	Long:  "StartServer the db test web server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

var loaderPath string

func init() {
	StartCmd.Flags().StringVar(&loaderPath, "loadpath", "", "randgen yy/zz directory path")
}

func StartServer() {
	initDatabase(config.GetConf().StandardDB)
	engine := gin.Default()

	MySQL, err1 := util.OpenDBWithRetry("mysql", config.GetConf().StandardDB)
	TiDB, err2 := util.OpenDBWithRetry("mysql", config.GetConf().TestDB)
	if err1 != nil || err2 != nil {
		log.Fatal("can not connect to db", zap.Error(err1), zap.Error(err2))
	}

	server := api.NewServer(executor.New(MySQL, TiDB))

	engine.POST("/tests", server.NewTest)
	engine.GET("/tests", server.ListTest)
	engine.GET("/tests/:id", server.GetTest)
	engine.GET("/results", server.ListTestResult)
	engine.GET("/results/:id/detail", server.ListLoopResult)

	engine.POST("/addfilter", server.AddFilter)
	go func() {
		LoadRandgenConf()
	}()

	log.Fatal("StartServer server failed", zap.String("err", engine.Run("0.0.0.0:8080").Error()))

}

var models = []interface{}{&types.Test{}, &types.TestResult{}, &types.LoopResult{}}

func initDatabase(dsn string) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("can not parse the mysql configuration")
	}
	if cfg.DBName == "" {
		db, err := util.OpenDBWithRetry("mysql", cfg.FormatDSN())
		if err != nil {
			panic("can not connect to database")
		}
		db.Exec("create database if not exists dbtest")
		db.Close()
		cfg.DBName = "dbtest"
	}
	db.InitDatabase(cfg.FormatDSN(), models)
}

func LoadRandgenConf() {
	if loaderPath == "" {
		return
	}

	for{
		_, err := http.Get("http://localhost:8080")
		if err == nil {
			break
		}

		log.Warn("Wait server started:")
		time.Sleep(time.Second)
	}

	infos, err := ioutil.ReadDir(loaderPath)
	if err != nil {
		log.Error("load conf files failed", zap.Error(err))
		return
	}

	submmitedFiles := make(map[string]struct{})
	for _, info := range infos {
		valid, newName := rename(info.Name())
		if !valid {
			continue
		}

		if _, ok := submmitedFiles[newName]; ok {
			continue
		}
		submmitedFiles[newName] = struct{}{}

		yy := fmt.Sprintf("%s/%s", loaderPath, rebuildName(newName, "yy"))
		zz := fmt.Sprintf("%s/%s", loaderPath, rebuildName(newName, "zz"))

		if !PathExist(yy) || !PathExist(zz) {
			log.Info(fmt.Sprintf("%s/%s", loaderPath, yy))
			continue
		}

		log.Info("submit rangen file", zap.String("yy", yy), zap.String("zz", zz))
		NewRandgenConf(yy, zz, 1000).Add()
	}
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func rebuildName(name string, suffix string) string {
	return fmt.Sprintf("%s.%s", name, suffix)
}

func rename(fileName string) (bool, string){
	ext := path.Ext(fileName)
	return ext == ".zz" || ext == ".yy", strings.TrimSuffix(fileName, ext)
}
