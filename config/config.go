package config

import (
	"flag"
)

type Config struct {
	Loop         int
	LoopInterval int

	DataLoaders  string
	QueryLoaders string
	Comparor     string
	CellFilter   string

	StandardDB string
	TestDB     string
}

var conf *Config

func init() {
	conf = &Config{}
	flag.IntVar(&conf.Loop, "loop", 1, "the loop count")
	flag.IntVar(&conf.LoopInterval, "loop-interval", 10, "the second to sleep after a loop is finished")
	flag.StringVar(&conf.DataLoaders, "data-loaders", "dummy", "a list of data loader names split by comma")
	flag.StringVar(&conf.QueryLoaders, "query-loaders", "dummy", "a list of query loader names split by comma")
	flag.StringVar(&conf.Comparor, "comparor", "standard", "the compare plugin")
	flag.StringVar(&conf.Comparor, "cell-filter", "standard", "the cell filter plugin")
	flag.StringVar(&conf.StandardDB, "standard-db", "root:@tcp(127.0.0.1:3306)/tep?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.StringVar(&conf.TestDB, "test-db", "root:@tcp(127.0.0.1:4000)/tep?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.Parse()
}

func GetConf() *Config {
	return conf
}
