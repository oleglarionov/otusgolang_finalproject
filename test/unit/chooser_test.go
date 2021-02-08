package unit

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"github.com/oleglarionov/otusgolang_finalproject/internal/infrastructure/repository/memory"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChooserImpl_ChooseBanner(t *testing.T) {
	t.Run("choose banner from correct slot/group", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		chooser := buildChooser(ctx, t, []BannerStat{
			{
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-1",
				views:     1,
				clicks:    1,
			}, {
				slot:      "slot-1",
				userGroup: "group-2",
				banner:    "banner-2",
				views:     1,
				clicks:    1,
			}, {
				slot:      "slot-2",
				userGroup: "group-1",
				banner:    "banner-3",
				views:     1,
				clicks:    1,
			},
		})

		banner, err := chooser.ChooseBanner(ctx, "slot-2", "group-1")
		require.NoError(t, err)

		require.Equal(t, banerrotation.BannerID("banner-3"), banner)
	})

	t.Run("choose banner with zero views", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		chooser := buildChooser(ctx, t, []BannerStat{
			{
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-1",
				views:     1,
				clicks:    0,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-2",
				views:     2,
				clicks:    0,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-3",
				views:     0,
				clicks:    0,
			},
		})

		banner, err := chooser.ChooseBanner(ctx, "slot-1", "group-1")
		require.NoError(t, err)

		require.Equal(t, banerrotation.BannerID("banner-3"), banner)
	})

	t.Run("choose banner with min views", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		chooser := buildChooser(ctx, t, []BannerStat{
			{
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-1",
				views:     10,
				clicks:    0,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-2",
				views:     11,
				clicks:    0,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-3",
				views:     12,
				clicks:    0,
			},
		})

		banner, err := chooser.ChooseBanner(ctx, "slot-1", "group-1")
		require.NoError(t, err)

		require.Equal(t, banerrotation.BannerID("banner-1"), banner)
	})

	t.Run("choose banner with max clicks", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		chooser := buildChooser(ctx, t, []BannerStat{
			{
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-1",
				views:     100,
				clicks:    1,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-2",
				views:     100,
				clicks:    91,
			}, {
				slot:      "slot-1",
				userGroup: "group-1",
				banner:    "banner-3",
				views:     100,
				clicks:    92,
			},
		})

		banner, err := chooser.ChooseBanner(ctx, "slot-1", "group-1")
		require.NoError(t, err)

		require.Equal(t, banerrotation.BannerID("banner-3"), banner)
	})
}

func buildChooser(ctx context.Context, t *testing.T, bannerStats []BannerStat) *banerrotation.ChooserImpl {
	bannerRepo := memory.NewBannerRepository()
	counterRepo := memory.NewCounterRepository()

	for _, bannerStat := range bannerStats {
		err := bannerRepo.AddBanner(ctx, bannerStat.slot, bannerStat.banner)
		require.NoError(t, err)

		for i := uint64(0); i < bannerStat.views; i++ {
			err := counterRepo.IncrementViews(ctx, bannerStat.slot, bannerStat.userGroup, bannerStat.banner)
			require.NoError(t, err)
		}

		for i := uint64(0); i < bannerStat.clicks; i++ {
			err := counterRepo.IncrementClicks(ctx, bannerStat.slot, bannerStat.userGroup, bannerStat.banner)
			require.NoError(t, err)
		}
	}

	return banerrotation.NewChooserImpl(bannerRepo, counterRepo)
}

type BannerStat struct {
	slot      banerrotation.SlotID
	userGroup banerrotation.UserGroupID
	banner    banerrotation.BannerID
	views     uint64
	clicks    uint64
}
