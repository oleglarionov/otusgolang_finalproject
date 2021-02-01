package usecase

import (
	"context"
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
	"log"
	"time"
)

type BannerRotation interface {
	AddBanner(ctx context.Context, slotID string, bannerID string) error
	RemoveBanner(ctx context.Context, slotID string, bannerID string) error
	ChooseBanner(ctx context.Context, slotID string, userGroupID string) (string, error)
	CountClick(ctx context.Context, slotID string, bannerID string, userGroupID string) error
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

func (u *BannerRotationImpl) ChooseBanner(ctx context.Context, userGroupID string, slotID string) (string, error) {
	slot := banerrotation.SlotID(slotID)
	userGroup := banerrotation.UserGroupID(userGroupID)

	banner, err := u.chooser.ChooseBanner(ctx, slot, userGroup)
	if err != nil {
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

func (u *BannerRotationImpl) CountClick(ctx context.Context, slotID string, userGroupID string, bannerID string) error {
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
