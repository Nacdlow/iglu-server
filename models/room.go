package models

type RType int64

const (
	Lounge = iota //iota acts as an auto-increment, Lounge = 1, Bedroom = 2, etc.
	Bedroom
	Dining
	Kitchen
	Bathroom
	Hallway
	Storage
	Utility
	Garage
	Guest
)

type Room struct {
	RoomID      int64 `xorm:"pk autoincr"`
	RoomName    string
	Description string
	RoomType    RType
}
