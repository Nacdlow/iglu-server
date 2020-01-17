package models

// RType is the room type.
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

// Room represents a room in a house, it includes a description and type of the
// room.
type Room struct {
	RoomID      int64 `xorm:"pk autoincr"`
	RoomName    string
	Description string
	RoomType    RType
}