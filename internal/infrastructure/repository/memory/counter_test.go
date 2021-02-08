package memory

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCounterRepository_IncrementViews(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewCounterRepository()

	err := repo.IncrementViews(ctx, "slot-1", "group-1", "banner-1")
	require.NoError(t, err)
}

func TestCounterRepository_IncrementClicks(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewCounterRepository()

	err := repo.IncrementClicks(ctx, "slot-1", "group-1", "banner-1")
	require.NoError(t, err)
}

func TestCounterRepository_GetCounters(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := NewCounterRepository()

	_, err := repo.GetCounters(ctx, "slot-1", "group-1", []banerrotation.BannerID{"banner-1"})
	require.NoError(t, err)
}
