package usecase

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/repository/memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBannerRotationImpl(t *testing.T) {
	banners := []initBanner{
		{"slot-1", "banner-1"},
		{"slot-1", "banner-2"},
		{"slot-1", "banner-3"},
		{"slot-2", "banner-1"},
		{"slot-2", "banner-4"},
	}

	t.Run("test views without clicks", func(t *testing.T) {
		uc := buildBannerRotationImpl(t, banners)
		ctx := context.Background()

		viewsByBanner := make(map[string]int)
		n := 1000
		for i := 0; i < n; i++ {
			banner, err := uc.ChooseBanner(ctx, "slot-1", "group-1")
			require.NoError(t, err)
			require.Subset(t,
				[]string{"banner-1", "banner-2", "banner-3", "banner-4"},
				[]string{banner},
			)
			viewsByBanner[banner] += 1
		}

		minCounterValue := n / len(viewsByBanner)
		maxCounterValue := minCounterValue
		if n%len(viewsByBanner) > 0 {
			maxCounterValue += 1
		}

		for _, views := range viewsByBanner {
			require.GreaterOrEqual(t, views, minCounterValue)
			require.LessOrEqual(t, views, maxCounterValue)
		}
	})

	t.Run("test with clicks on one banner", func(t *testing.T) {
		uc := buildBannerRotationImpl(t, banners)
		ctx := context.Background()

		viewsByBanner := make(map[string]int)
		n := 1000
		for i := 0; i < n; i++ {
			banner, err := uc.ChooseBanner(ctx, "slot-1", "group-1")
			require.NoError(t, err)
			viewsByBanner[banner] += 1

			if banner == "banner-1" {
				err = uc.RegisterClick(ctx, "slot-1", "group-1", banner)
				require.NoError(t, err)
			}
		}

		for _, views := range viewsByBanner {
			require.GreaterOrEqual(t, views, 1)
		}
	})
}

func buildBannerRotationImpl(t *testing.T, banners []initBanner) *BannerRotationImpl {
	bannerRepo := memory.NewBannerRepository()
	counterRepo := memory.NewCounterRepository()
	chooser := banerrotation.NewChooserImpl(bannerRepo, counterRepo)
	streamer := &streamerStub{}

	ctx := context.Background()
	for _, b := range banners {
		err := bannerRepo.AddBanner(ctx, b.slotID, b.bannerID)
		require.NoError(t, err)
	}

	return NewBannerRotationImpl(chooser, bannerRepo, counterRepo, streamer)
}

type initBanner struct {
	slotID   banerrotation.SlotID
	bannerID banerrotation.BannerID
}

type streamerStub struct {
}

func (s *streamerStub) Push(event event.Event) error {
	return nil
}
