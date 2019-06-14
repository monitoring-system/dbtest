package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var AddTestCmd = &cobra.Command{
	Use:   "add",
	Short: "add a test to the db test web server",
	Long:  "add a test to the db test web server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if loaderPath == "" {
			addTest()
		} else {
			LoadRandgenConf()
		}
	},
}
var yyFile string
var zzFile string
var loop int
var loopInterval int
var queryCount int
var queryLoader string
var dataLoader string
var name string
var loaderPath string

func init() {
	AddTestCmd.Flags().StringVar(&yyFile, "yy", "", "randgen yy file path")
	AddTestCmd.Flags().StringVar(&zzFile, "zz", "", "randgen zz file path")
	AddTestCmd.Flags().IntVar(&loop, "loop", 1, "how many round the test should be run")
	AddTestCmd.Flags().IntVar(&loopInterval, "loop-interval", 30, "sleep time(seconds) between each loop")
	AddTestCmd.Flags().IntVar(&queryCount, "query-count", 1000, "the generated sql count")
	AddTestCmd.Flags().StringVar(&queryLoader, "query-loader", "randgen", "query loaders, split by comma")
	AddTestCmd.Flags().StringVar(&dataLoader, "data-loader", "randgen", "data loaders, split by comma")
	AddTestCmd.Flags().StringVar(&name, "name", "console", "data loaders, split by comma")
	AddTestCmd.Flags().StringVar(&loaderPath, "loadpath", "", "randgen yy/zz directory path")
}

type randgenConf struct {
	yyFile string
	zzFile string
	loop   int
}

type RandgenConfOpt struct {
	conf         randgenConf
	LoopInterval int
	QueryCount   int
	QueryLoader  string
	DataLoader   string
	Name         string
}

func NewRandgenConf(yyFile string, zzFile string, loop int) *RandgenConfOpt {
	return &RandgenConfOpt{
		conf: randgenConf{
			yyFile: yyFile,
			zzFile: zzFile,
			loop:   loop,
		},
		LoopInterval: loopInterval,
		QueryCount:   queryCount,
		QueryLoader:  "randgen",
		DataLoader:   "randgen",
		Name:         "console",
	}
}

func (c *RandgenConfOpt) Add() {
	yyContent, err := ioutil.ReadFile(c.conf.yyFile)
	if err != nil {
		log.Fatal(err)
	}
	zzContent, err := ioutil.ReadFile(c.conf.zzFile)
	if err != nil {
		log.Fatal(err)
	}
	payload := &types.Test{TestName: c.Name, Yy: string(yyContent), Zz: string(zzContent), Loop: c.conf.loop, LoopInterval: c.LoopInterval, QueryLoader: c.QueryLoader, DataLoader: c.DataLoader}
	resp, err := http.Post("http://localhost:8080/tests", "application/json",
		strings.NewReader(getPayload(payload)))

	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		fmt.Println("no response get from server")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal([]byte(body), payload)
	fmt.Println("add test successfully")
	fmt.Println(getPayload(payload))
}

func addTest() {
	opt := NewRandgenConf(yyFile, zzFile, loop)
	opt.LoopInterval = loopInterval
	opt.QueryCount = queryCount
	opt.QueryLoader = queryLoader
	opt.DataLoader = dataLoader
	opt.Name = name

	opt.Add()
}

func getPayload(payload *types.Test) string {
	bytes, _ := json.Marshal(payload)
	return string(bytes)
}

func LoadRandgenConf() {
	if loaderPath == "" {
		return
	}

	for {
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
		NewRandgenConf(yy, zz, loop).Add()
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

func rename(fileName string) (bool, string) {
	ext := path.Ext(fileName)
	return ext == ".zz" || ext == ".yy", strings.TrimSuffix(fileName, ext)
}
