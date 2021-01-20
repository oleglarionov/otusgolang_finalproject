//+build wireinject

package main

import (
	"github.com/google/wire"
	internalgrpc "github.com/oleglarionov/otusgolang_finalproject/internal/grpc"
	grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"
)

func setup(cfg Config) (*App, error) {
	wire.Build(
		NewApp,
		grpcServerProvider,
		wire.Bind(new(grpcgenerated.BannerRotationServiceServer), new(*internalgrpc.BannerRotationServerImpl)),
		internalgrpc.NewBannerRotationServerImpl,
	)
	return nil, nil
}

func grpcServerProvider(cfg Config, service grpcgenerated.BannerRotationServiceServer) *internalgrpc.Server {
	return internalgrpc.NewServer(cfg.ServerPort, service)
}
