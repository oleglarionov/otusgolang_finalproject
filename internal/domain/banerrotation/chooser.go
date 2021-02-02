package banerrotation

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/algorithm"
)

type Chooser interface {
	ChooseBanner(ctx context.Context, slotID SlotID, userGroupID UserGroupID) (BannerID, error)
}

type ChooserImpl struct {
	bannerRepository  BannerRepository
	counterRepository CounterRepository
}

func NewChooserImpl(bannerRepository BannerRepository, counterRepository CounterRepository) *ChooserImpl {
	return &ChooserImpl{
		bannerRepository:  bannerRepository,
		counterRepository: counterRepository,
	}
}

func (s *ChooserImpl) ChooseBanner(ctx context.Context, slot SlotID, userGroup UserGroupID) (BannerID, error) {
	banners, err := s.bannerRepository.GetBanners(ctx, slot)
	if err != nil {
		return "", err
	}

	counters, err := s.counterRepository.GetCounters(ctx, slot, userGroup, banners)
	if err != nil {
		return "", err
	}

	l := len(counters)
	avgIncome := make([]float64, 0, l)
	nj := make([]uint64, 0, l)
	n := uint64(0)

	for _, counter := range counters {
		curIncome := 0.0
		if counter.Views != 0 {
			curIncome = float64(counter.Clicks) / float64(counter.Views)
		}

		avgIncome = append(avgIncome, curIncome)
		nj = append(nj, counter.Views)
		n += counter.Views
	}

	j := algorithm.Ucb1(avgIncome, nj, n)

	bannerID := counters[j].BannerID
	err = s.counterRepository.IncrementViews(ctx, slot, userGroup, bannerID)
	if err != nil {
		return "", err
	}

	return bannerID, nil
}
