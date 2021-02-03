package event

import (
	"time"

	"github.com/oleglarionov/otusgolang_finalproject/internal/domain/banerrotation"
)

type Type string

const (
	View  Type = "view"
	Click Type = "click"
)

type Event struct {
	Type      Type
	Slot      banerrotation.SlotID
	Banner    banerrotation.BannerID
	UserGroup banerrotation.UserGroupID
	Time      time.Time
}

type Streamer interface {
	Push(event Event) error
}
