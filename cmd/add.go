package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/monitoring-system/dbtest/api/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var AddTestCmd = &cobra.Command{
	Use:   "add",
	Short: "add a test to the db test web server",
	Long:  "add a test to the db test web server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		addTest()
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

func init() {
	AddTestCmd.Flags().StringVar(&yyFile, "yy", "", "randgen yy file path")
	AddTestCmd.Flags().StringVar(&zzFile, "zz", "", "randgen zz file path")
	AddTestCmd.Flags().IntVar(&loop, "loop", 1, "how many round the test should be run")
	AddTestCmd.Flags().IntVar(&loopInterval, "loop-interval", 30, "sleep time(seconds) between each loop")
	AddTestCmd.Flags().IntVar(&queryCount, "query-count", 1000, "the generated sql count")
	AddTestCmd.Flags().StringVar(&queryLoader, "query-loader", "randgen", "query loaders, split by comma")
	AddTestCmd.Flags().StringVar(&dataLoader, "data-loader", "randgen", "data loaders, split by comma")
	AddTestCmd.Flags().StringVar(&name, "name", "console", "data loaders, split by comma")
}

type randgenConf struct {
	yyFile string
	zzFile string
	loop int
}

type RandgenConfOpt struct {
	conf randgenConf
	LoopInterval int
	QueryCount int
	QueryLoader string
	DataLoader string
	Name string
}

func NewRandgenConf(yyFile string, zzFile string, loop int) *RandgenConfOpt {
	return &RandgenConfOpt{
		conf: randgenConf{
			yyFile: yyFile,
			zzFile: zzFile,
			loop: loop,
		},
		LoopInterval: 30,
		QueryCount: 1000,
		QueryLoader: "randgen",
		DataLoader: "randgen",
		Name: "console",
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
