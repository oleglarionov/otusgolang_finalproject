package memory

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBannerRepository_AddBanner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewBannerRepository()
	err := repo.AddBanner(ctx, "slot-1", "banner-1")
	require.NoError(t, err)
}

func TestBannerRepository_RemoveBanner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewBannerRepository()
	err := repo.AddBanner(context.Background(), "slot-1", "banner-1")
	require.NoError(t, err)

	err = repo.RemoveBanner(ctx, "slot-1", "banner-1")
	require.NoError(t, err)
}

func TestBannerRepository_GetBanners(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewBannerRepository()
	err := repo.AddBanner(context.Background(), "slot-1", "banner-1")
	require.NoError(t, err)

	err = repo.AddBanner(context.Background(), "slot-1", "banner-2")
	require.NoError(t, err)

	_, err = repo.GetBanners(ctx, "slot-1")
	require.NoError(t, err)
}
