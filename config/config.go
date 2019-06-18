package config

type Config struct {
	StandardDB     string
	TestDB         string
	TraceAllErrors bool
	RandGenServer  string
}

var Conf *Config

func GetConf() *Config {
	return Conf
}
