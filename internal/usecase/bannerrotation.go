package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type BannerRotation interface {
	AddBanner(ctx context.Context, slotID string, bannerID string) error
	RemoveBanner(ctx context.Context, slotID string, bannerID string) error
	ChooseBanner(ctx context.Context, slotID string, userGroupID string) (string, error)
	RegisterClick(ctx context.Context, slotID string, bannerID string, userGroupID string) error
}

type BannerRotationImpl struct {
	chooser           banerrotation.Chooser
	bannerRepository  banerrotation.BannerRepository
	counterRepository banerrotation.CounterRepository
	streamer          event.Streamer
}

func NewBannerRotationImpl(
	chooser banerrotation.Chooser,
	bannerRepository banerrotation.BannerRepository,
	counterRepository banerrotation.CounterRepository,
	streamer event.Streamer,
) *BannerRotationImpl {
	return &BannerRotationImpl{
		chooser:           chooser,
		bannerRepository:  bannerRepository,
		counterRepository: counterRepository,
		streamer:          streamer,
	}
}

func (u *BannerRotationImpl) AddBanner(ctx context.Context, slotID string, bannerID string) error {
	return u.bannerRepository.AddBanner(
		ctx,
		banerrotation.SlotID(slotID),
		banerrotation.BannerID(bannerID),
	)
}

func (u *BannerRotationImpl) RemoveBanner(ctx context.Context, slotID string, bannerID string) error {
	return u.bannerRepository.RemoveBanner(
		ctx,
		banerrotation.SlotID(slotID),
		banerrotation.BannerID(bannerID),
	)
}

func (u *BannerRotationImpl) ChooseBanner(ctx context.Context, slotID string, userGroupID string) (string, error) {
	slot := banerrotation.SlotID(slotID)
	userGroup := banerrotation.UserGroupID(userGroupID)

	banner, err := u.chooser.ChooseBanner(ctx, slot, userGroup)
	if err != nil {
		if errors.Is(err, banerrotation.ErrNoBanners) {
			return "", nil
		}

		return "", err
	}

	err = u.streamer.Push(event.Event{
		Type:      event.View,
		Slot:      slot,
		Banner:    banner,
		UserGroup: userGroup,
		Time:      time.Now(),
	})
	if err != nil {
		log.Println(err)
	}

	return string(banner), nil
}

func (u *BannerRotationImpl) RegisterClick(ctx context.Context, slotID string, bannerID string, userGroupID string) error {
	slot := banerrotation.SlotID(slotID)
	userGroup := banerrotation.UserGroupID(userGroupID)
	banner := banerrotation.BannerID(bannerID)

	err := u.counterRepository.IncrementClicks(
		ctx,
		banerrotation.SlotID(slotID),
		banerrotation.UserGroupID(userGroupID),
		banerrotation.BannerID(bannerID),
	)
	if err != nil {
		return err
	}

	err = u.streamer.Push(event.Event{
		Type:      event.Click,
		Slot:      slot,
		Banner:    banner,
		UserGroup: userGroup,
		Time:      time.Now(),
	})
	if err != nil {
		log.Println(err)
	}

	return nil
}
