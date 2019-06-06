package config

import (
	"flag"
)

type Config struct {
	StandardDB string
	TestDB     string
}

var conf *Config

func init() {
	conf = &Config{}
	flag.StringVar(&conf.StandardDB, "standard-db", "root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.StringVar(&conf.TestDB, "test-db", "root:@tcp(127.0.0.1:4000)/?charset=utf8&parseTime=True&loc=Local", "the compare plugin")
	flag.Parse()
}

func GetConf() *Config {
	return conf
}
