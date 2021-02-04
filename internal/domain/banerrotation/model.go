package banerrotation

type SlotID string
type UserGroupID string
type BannerID string

type Counter struct {
	Slot      SlotID      `db:"slot_id"`
	UserGroup UserGroupID `db:"user_group_id"`
	Banner    BannerID    `db:"banner_id"`
	Views     uint64      `db:"views"`
	Clicks    uint64      `db:"clicks"`
}
