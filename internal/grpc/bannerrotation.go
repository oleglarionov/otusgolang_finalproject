package internalgrpc

import (
	"context"
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

	return &grpcgenerated.ChooseBannerResponse{
		BannerId: bannerID,
	}, nil
}
