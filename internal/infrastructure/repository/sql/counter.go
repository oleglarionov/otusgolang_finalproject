package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type CounterRepository struct {
	db *sqlx.DB
}

func NewCounterRepository(db *sqlx.DB) *CounterRepository {
	return &CounterRepository{db: db}
}

func (r *CounterRepository) GetCounters(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banners []banerrotation.BannerID) ([]banerrotation.Counter, error) {
	panic("implement me")
}

func (r *CounterRepository) IncrementViews(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	panic("implement me")
}

func (r *CounterRepository) IncrementClicks(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	panic("implement me")
}
