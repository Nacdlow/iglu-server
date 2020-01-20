package models

import (
	"errors"
)

// RType is the room type.
type RType int64

// RType enums.
const (
	LoungeRoom = iota // iota acts as an auto-increment, Lounge = 0, Bedroom = 1, etc.
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

// GetRoom gets a Room based on its ID from the database.
func GetRoom(id int64) (*Room, error) {
	r := new(Room)
	has, err := engine.ID(id).Get(r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("Room does not exist")
	}
	return r, nil
}

// GetRooms returns an array of all Rooms from the database.
func GetRooms() (room []Room) {
	engine.Find(&room)
	return room
}

// AddRoom adds a Room in the database.
func AddRoom(r *Room) (err error) {
	_, err = engine.Insert(r)
	return
}

// HasRoom returns whether an Room is in the database or not.
func HasRoom(id int64) (has bool) {
	has, _ = engine.Get(&Room{RoomID: id})
	return
}

// UpdateRoom updates an Room in the database.
func UpdateRoom(r *Room) (err error) {
	_, err = engine.Id(r.RoomID).Update(r)
	return
}

// UpdateRoomCols will update the columns of an item even if they are empty.
func UpdateRoomCols(r *Room, cols ...string) (err error) {
	_, err = engine.ID(r.RoomID).Cols(cols...).Update(r)
	return
}
