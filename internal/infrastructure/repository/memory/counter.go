package memory

import (
	"context"
	"sync"

	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type CounterKey struct {
	SlotID      banerrotation.SlotID
	UserGroupID banerrotation.UserGroupID
	BannerID    banerrotation.BannerID
}

type CounterRepository struct {
	mu   sync.RWMutex
	data map[CounterKey]banerrotation.Counter
}

func NewCounterRepository() *CounterRepository {
	data := make(map[CounterKey]banerrotation.Counter)
	return &CounterRepository{data: data}
}

func (r *CounterRepository) GetCounters(_ context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banners []banerrotation.BannerID) ([]banerrotation.Counter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counters := make([]banerrotation.Counter, 0, len(banners))

	for _, banner := range banners {
		counter := r.getCounter(slot, userGroup, banner)
		counters = append(counters, counter)
	}

	return counters, nil
}

func (r *CounterRepository) IncrementViews(_ context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter := r.getCounter(slot, userGroup, banner)
	counter.Views++

	r.saveCounter(counter)

	return nil
}

func (r *CounterRepository) IncrementClicks(_ context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter := r.getCounter(slot, userGroup, banner)
	counter.Clicks++

	r.saveCounter(counter)

	return nil
}

func (r *CounterRepository) getCounter(slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) banerrotation.Counter {
	key := CounterKey{
		SlotID:      slot,
		UserGroupID: userGroup,
		BannerID:    banner,
	}

	counter, ok := r.data[key]
	if !ok {
		return banerrotation.Counter{
			Slot:      slot,
			UserGroup: userGroup,
			Banner:    banner,
			Views:     0,
			Clicks:    0,
		}
	}

	return counter
}

func (r *CounterRepository) saveCounter(counter banerrotation.Counter) {
	key := CounterKey{
		SlotID:      counter.Slot,
		UserGroupID: counter.UserGroup,
		BannerID:    counter.Banner,
	}

	r.data[key] = counter
}
