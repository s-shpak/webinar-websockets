package server

import "flag"

type Config struct {
	Port string
}

func GetConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.Port, "p", "4443", "порт websocket-сервера")

	flag.Parse()
	return cfg
}
