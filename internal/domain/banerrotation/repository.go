package banerrotation

import "context"

type CounterRepository interface {
	GetCounters(ctx context.Context, slot SlotID, userGroup UserGroupID, banners []BannerID) ([]Counter, error)
	IncrementViews(ctx context.Context, slot SlotID, userGroup UserGroupID, banner BannerID) error
	IncrementClicks(ctx context.Context, slot SlotID, userGroup UserGroupID, banner BannerID) error
}

type BannerRepository interface {
	AddBanner(ctx context.Context, slot SlotID, banner BannerID) error
	RemoveBanner(ctx context.Context, slot SlotID, banner BannerID) error
	GetBanners(ctx context.Context, slot SlotID) ([]BannerID, error)
}
