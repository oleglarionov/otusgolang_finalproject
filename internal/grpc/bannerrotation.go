package internalgrpc

import grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"

type BannerRotationServerImpl struct {
	grpcgenerated.UnimplementedBannerRotationServiceServer
}

func NewBannerRotationServerImpl() *BannerRotationServerImpl {
	return &BannerRotationServerImpl{}
}
