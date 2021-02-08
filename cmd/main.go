package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	internalgrpc "github.com/oleglarionov/otusgolang_finalproject/internal/grpc"
	"github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/streamer"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type App struct {
	server *internalgrpc.Server
}

func NewApp(server *internalgrpc.Server) *App {
	return &App{server: server}
}

func (a *App) Serve() error {
	return a.server.Serve()
}

func (a *App) Stop() {
	a.server.Stop()
}

func main() {
	viper.AutomaticEnv()
	cfg := Config{
		ServerPort: viper.Get("SERVER_PORT").(string),
		DBDsn:      viper.Get("DB_DSN").(string),
		Rabbit: streamer.AMQPConfig{
			Dsn:   viper.Get("RABBIT_DSN").(string),
			Queue: viper.Get("RABBIT_QUEUE").(string),
		},
	}

	app, cleanup, err := setup(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	go signalHandler(ctx, app, cancel)

	go func() {
		err := app.Serve()
		if err == nil || errors.Is(err, grpc.ErrServerStopped) {
			log.Println("grpc server stopped")
		} else {
			log.Fatal("failed to start grpc server: " + err.Error())
		}
	}()

	<-ctx.Done()
}

func signalHandler(ctx context.Context, app *App, cancel context.CancelFunc) {
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	select {
	case <-signals:
		log.Println("stop signal received")
		app.Stop()
	case <-ctx.Done():
	}
}
