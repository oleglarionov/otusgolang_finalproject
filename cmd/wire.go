//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	internalgrpc "github.com/oleglarionov/otusgolang_finalproject/internal/grpc"
	grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"
	"github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/repository/sql"
	"github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/streamer"
	"github.com/oleglarionov/otusgolang_finalproject/internal/usecase"
)

func setup(cfg Config) (*App, error) {
	wire.Build(
		NewApp,
		grpcServerProvider,
		wire.Bind(new(grpcgenerated.BannerRotationServiceServer), new(*internalgrpc.BannerRotationServerImpl)),
		internalgrpc.NewBannerRotationServerImpl,
		wire.Bind(new(usecase.BannerRotation), new(*usecase.BannerRotationImpl)),
		usecase.NewBannerRotationImpl,
		wire.Bind(new(banerrotation.Chooser), new(*banerrotation.ChooserImpl)),
		banerrotation.NewChooserImpl,
		wire.Bind(new(banerrotation.CounterRepository), new(*sql.CounterRepository)),
		sql.NewCounterRepository,
		wire.Bind(new(banerrotation.BannerRepository), new(*sql.BannerRepository)),
		sql.NewBannerRepository,
		wire.Bind(new(event.Streamer), new(*streamer.AMQPStreamer)),
		streamer.NewAMQPStreamer,
		dbProvider,
	)
	return nil, nil
}

func grpcServerProvider(cfg Config, service grpcgenerated.BannerRotationServiceServer) *internalgrpc.Server {
	return internalgrpc.NewServer(cfg.ServerPort, service)
}

func dbProvider(cfg Config) (*sqlx.DB, error) {
	return sql.NewDB(cfg.DBDsn)
}
