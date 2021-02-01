package memory

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"sync"
)

type CounterRepository struct {
	mu   sync.RWMutex
	data map[banerrotation.SlotID]map[banerrotation.UserGroupID]map[banerrotation.BannerID]banerrotation.Counter
}

func NewCounterRepository() *CounterRepository {
	data := make(map[banerrotation.SlotID]map[banerrotation.UserGroupID]map[banerrotation.BannerID]banerrotation.Counter)
	return &CounterRepository{data: data}
}

func (r *CounterRepository) GetCounters(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banners []banerrotation.BannerID) ([]banerrotation.Counter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	counters := make([]banerrotation.Counter, 0, len(banners))

	for _, banner := range banners {
		counter := r.getCounter(slot, userGroup, banner)
		counters = append(counters, counter)
	}

	return counters, nil
}

func (r *CounterRepository) IncrementViews(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter := r.getCounter(slot, userGroup, banner)
	counter.Views++

	r.saveCounter(counter)

	return nil
}

func (r *CounterRepository) IncrementClicks(ctx context.Context, slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	counter := r.getCounter(slot, userGroup, banner)
	counter.Clicks++

	r.saveCounter(counter)

	return nil
}

func (r *CounterRepository) getCounter(slot banerrotation.SlotID, userGroup banerrotation.UserGroupID, banner banerrotation.BannerID) banerrotation.Counter {
	emptyCounter := banerrotation.Counter{
		SlotID:      slot,
		UserGroupID: userGroup,
		BannerID:    banner,
		Views:       0,
		Clicks:      0,
	}

	slotData, ok := r.data[slot]
	if !ok {
		return emptyCounter
	}

	userGroupData, ok := slotData[userGroup]
	if !ok {
		return emptyCounter
	}

	counter, ok := userGroupData[banner]
	if !ok {
		return emptyCounter
	}

	return counter
}

func (r *CounterRepository) saveCounter(counter banerrotation.Counter) {
	slotData, ok := r.data[counter.SlotID]
	if !ok {
		r.data[counter.SlotID] = make(map[banerrotation.UserGroupID]map[banerrotation.BannerID]banerrotation.Counter)
	}

	_, ok = slotData[counter.UserGroupID]
	if !ok {
		r.data[counter.SlotID][counter.UserGroupID] = make(map[banerrotation.BannerID]banerrotation.Counter)
	}

	r.data[counter.SlotID][counter.UserGroupID][counter.BannerID] = counter
}
