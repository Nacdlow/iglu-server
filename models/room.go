package models

import (
	"errors"
)

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

func GetRooms() (room []Room) {
	engine.Find(&room)
	return room
}

func AddRoom(r *Room) (err error) {
	_, err = engine.Insert(r)
	return
}

func HasRoom(id int64) (has bool) {
	has, _ = engine.Get(&Room{RoomID: id})
	return
}

func UpdateRoom(r *Room) (err error) {
	_, err = engine.Id(r.RoomID).Update(r)
	return
}

func UpdateRoomCols(r *Room, cols ...string) (err error) {
	_, err = engine.ID(r.RoomID).Cols(cols...).Update(r)
	return
}
