package banerrotation

import (
	"context"

	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/algorithm"
	"github.com/pkg/errors"
)

var ErrNoBanners = errors.New("no banners")

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

	if len(banners) == 0 {
		return "", ErrNoBanners
	}

	counters, err := s.counterRepository.GetCounters(ctx, slot, userGroup, banners)
	if err != nil {
		return "", err
	}

	l := len(counters)
	stats := make([]algorithm.Ucb1ElStat, 0, l)

	for _, counter := range counters {
		curIncome := 0.0
		if counter.Views != 0 {
			curIncome = float64(counter.Clicks) / float64(counter.Views)
		}

		curStat := algorithm.Ucb1ElStat{
			AvgIncome: curIncome,
			Attempts:  counter.Views,
		}
		stats = append(stats, curStat)
	}

	j := algorithm.Ucb1(stats)

	bannerID := counters[j].Banner
	err = s.counterRepository.IncrementViews(ctx, slot, userGroup, bannerID)
	if err != nil {
		return "", err
	}

	return bannerID, nil
}
