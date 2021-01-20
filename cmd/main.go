package main

import (
	"errors"
	internalgrpc "github.com/oleglarionov/otusgolang_finalproject/internal/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
)

type App struct {
	server *internalgrpc.Server
}

func NewApp(server *internalgrpc.Server) *App {
	return &App{server: server}
}

func main() {
	viper.AutomaticEnv()

	serverPort := viper.Get("SERVER_PORT").(string)
	cfg := Config{
		ServerPort: serverPort,
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
