package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"github.com/pkg/errors"
)

type CounterRepository struct {
	dbConnector DBConnector
}

func NewCounterRepository(dbConnector DBConnector) *CounterRepository {
	return &CounterRepository{dbConnector: dbConnector}
}

func (r *CounterRepository) GetCounters(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banners []banerrotation.BannerID) ([]banerrotation.Counter, error) {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bannersMap := make(map[banerrotation.BannerID]bool, len(banners))
	for _, banner := range banners {
		bannersMap[banner] = true
	}

	sql, args, err := sqlx.In("select * "+
		"from counters "+
		"where slot_id = ? "+
		"and user_group_id = ? "+
		"and banner_id in (?)",
		slot, userGroup, banners,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sql = db.Rebind(sql)
	rows, err := db.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	counters := make([]banerrotation.Counter, 0, len(banners))
	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var counter banerrotation.Counter
		err = rows.StructScan(&counter)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		counters = append(counters, counter)
		delete(bannersMap, counter.Banner)
	}

	for banner := range bannersMap {
		counter := banerrotation.Counter{
			Slot:      slot,
			UserGroup: userGroup,
			Banner:    banner,
			Views:     0,
			Clicks:    0,
		}
		counters = append(counters, counter)
	}

	return counters, nil
}

func (r *CounterRepository) IncrementViews(
	ctx context.Context,
	slot banerrotation.SlotID,
	userGroup banerrotation.UserGroupID,
	banner banerrotation.BannerID,
) error {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.ExecContext(ctx,
		"insert into counters"+
			"(slot_id, banner_id, user_group_id, views, clicks) "+
			"values ($1, $2, $3, 1, 0) "+
			"on conflict (slot_id, banner_id, user_group_id) "+
			"do update "+
			"set views = counters.views+1",
		slot, banner, userGroup,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *CounterRepository) IncrementClicks(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	db, err := r.dbConnector.GetConn()
	if err != nil {
		return errors.WithStack(err)
	}

	result, err := db.ExecContext(ctx,
		"update counters set clicks=clicks+1 where slot_id=$1 and banner_id=$2 and user_group_id=$3",
		slot, banner, userGroup,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}
	if affected == 0 {
		return errors.New("counter not updated")
	}

	return nil
}
