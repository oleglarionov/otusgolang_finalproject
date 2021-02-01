package banerrotation

type SlotID string
type UserGroupID string
type BannerID string

type Counter struct {
	SlotID
	UserGroupID
	BannerID
	Views  uint64
	Clicks uint64
}
