//go:generate protoc banner_rotation_service.proto --go_out=./generated --go-grpc_out=./generated --proto_path=../../api/

package internalgrpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"
	"github.com/oleglarionov/otusgolang_finalproject/internal/usecase"
)

type BannerRotationServerImpl struct {
	grpcgenerated.UnimplementedBannerRotationServiceServer

	useCase usecase.BannerRotation
}

func NewBannerRotationServerImpl(useCase usecase.BannerRotation) *BannerRotationServerImpl {
	return &BannerRotationServerImpl{
		useCase: useCase,
	}
}

func (s *BannerRotationServerImpl) ChooseBanner(
	ctx context.Context,
	req *grpcgenerated.ChooseBannerRequest,
) (*grpcgenerated.ChooseBannerResponse, error) {
	bannerID, err := s.useCase.ChooseBanner(ctx, req.SlotId, req.UserGroupId)
	if err != nil {
		return nil, err
	}

	bannerFound := true
	if bannerID == "" {
		bannerFound = false
	}

	return &grpcgenerated.ChooseBannerResponse{
		BannerId:    bannerID,
		BannerFound: bannerFound,
	}, nil
}

func (s *BannerRotationServerImpl) AddBanner(
	ctx context.Context,
	req *grpcgenerated.AddBannerRequest,
) (*empty.Empty, error) {
	err := s.useCase.AddBanner(ctx, req.SlotId, req.BannerId)
	return &empty.Empty{}, err
}

func (s *BannerRotationServerImpl) RemoveBanner(ctx context.Context, req *grpcgenerated.RemoveBannerRequest) (*empty.Empty, error) {
	err := s.useCase.RemoveBanner(ctx, req.SlotId, req.BannerId)
	return &empty.Empty{}, err
}

func (s *BannerRotationServerImpl) RegisterClick(ctx context.Context, req *grpcgenerated.RegisterClickRequest) (*empty.Empty, error) {
	err := s.useCase.RegisterClick(ctx, req.SlotId, req.BannerId, req.UserGroupId)
	return &empty.Empty{}, err
}
