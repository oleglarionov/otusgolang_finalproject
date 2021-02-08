package sql

import (
	"context"

	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"github.com/pkg/errors"
)

type BannerRepository struct {
	dbConnector DBConnector
}

func NewBannerRepository(dbConnector DBConnector) *BannerRepository {
	return &BannerRepository{dbConnector: dbConnector}
}

func (r *BannerRepository) AddBanner(ctx context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.ExecContext(
		ctx,
		"insert into slot_banners (slot_id, banner_id) values ($1, $2)", slot, banner,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *BannerRepository) RemoveBanner(ctx context.Context, slot banerrotation.SlotID, banner banerrotation.BannerID) error {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.ExecContext(
		ctx,
		"delete from slot_banners where slot_id = $1 and banner_id = $2", slot, banner,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *BannerRepository) GetBanners(ctx context.Context, slot banerrotation.SlotID) ([]banerrotation.BannerID, error) {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var banners []banerrotation.BannerID
	rows, err := db.QueryxContext(ctx, "select banner_id from slot_banners where slot_id = $1", slot)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var banner banerrotation.BannerID
		err = rows.Scan(&banner)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		banners = append(banners, banner)
	}

	return banners, nil
}
