package models

type RType int64

const (
	LoungeRoom = iota //iota acts as an auto-increment, Lounge = 1, Bedroom = 2, etc.
	BedroomRoom
	DiningRoom
	KitchenRoom
	BathroomRoom
	HallwayRoom
	StorageRoom
	UtilityRoom
	GarageRoom
	GuestRoom
)

type Room struct {
	RoomID      int64 `xorm:"pk autoincr"`
	RoomName    string
	Description string
	RoomType    RType
}
