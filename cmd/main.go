package main

import (
	"errors"
	"log"

	internalgrpc "github.com/oleglarionov/otusgolang_finalproject/internal/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type App struct {
	server *internalgrpc.Server
}

func NewApp(server *internalgrpc.Server) *App {
	return &App{server: server}
}

func main() {
	// todo: начать обрабатывать сигналы ОС

	viper.AutomaticEnv()

	serverPort := viper.Get("SERVER_PORT").(string)
	dbDsn := viper.Get("DB_DSN").(string)
	cfg := Config{
		ServerPort: serverPort,
		DBDsn:      dbDsn,
	}

	app, err := setup(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = app.server.Serve()
	if errors.Is(err, grpc.ErrServerStopped) {
		log.Println("grpc server stopped")
	} else {
		log.Fatal("failed to start grpc server: " + err.Error())
	}
}
