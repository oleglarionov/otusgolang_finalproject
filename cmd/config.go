package main

import "github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/streamer"

type Config struct {
	ServerPort string
	DBDsn      string
	Rabbit     streamer.AMQPConfig
}
