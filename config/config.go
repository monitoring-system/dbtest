package config

type Config struct {
	StandardDB string
	TestDB     string
}

var Conf *Config

func GetConf() *Config {
	return Conf
}
