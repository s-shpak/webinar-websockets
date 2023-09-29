package client

import "flag"

type Config struct {
	ServerPort string
	Msg        string
}

func GetConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.ServerPort, "p", "4443", "порт websocket-сервера")
	flag.StringVar(&cfg.Msg, "m", "Hi everyone!", "сообщение, которое клиент будет отправлять в чат")

	flag.Parse()
	return cfg
}
