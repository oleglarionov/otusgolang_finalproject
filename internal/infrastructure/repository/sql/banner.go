package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type BannerRepository struct {
	db *sqlx.DB
}

func NewBannerRepository(db *sqlx.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (r *BannerRepository) AddBanner(ctx context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	panic("implement me")
}

func (r *BannerRepository) RemoveBanner(ctx context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	panic("implement me")
}

func (r *BannerRepository) GetBanners(ctx context.Context, slot banerrotation.SlotID) ([]banerrotation.BannerID, error) {
	panic("implement me")
}
