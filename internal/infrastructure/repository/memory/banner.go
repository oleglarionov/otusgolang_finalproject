package memory

import (
	"context"

	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type BannerRepository struct {
	data map[banerrotation.SlotID]map[banerrotation.BannerID]bool
}

func NewBannerRepository() *BannerRepository {
	data := make(map[banerrotation.SlotID]map[banerrotation.BannerID]bool)
	return &BannerRepository{data: data}
}

func (r *BannerRepository) AddBanner(_ context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	_, ok := r.data[slot]
	if !ok {
		r.data[slot] = make(map[banerrotation.BannerID]bool)
	}

	r.data[slot][banner] = true
	return nil
}

func (r *BannerRepository) RemoveBanner(ctx context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	delete(r.data[slot], banner)
	return nil
}

func (r *BannerRepository) GetBanners(ctx context.Context, slot banerrotation.SlotID) ([]banerrotation.BannerID, error) {
	l := len(r.data[slot])
	banners := make([]banerrotation.BannerID, 0, l)
	for banner := range r.data[slot] {
		banners = append(banners, banner)
	}

	return banners, nil
}
