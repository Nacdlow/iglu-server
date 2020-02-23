package models

import (
	"errors"

	"github.com/brianvoe/gofakeit/v4"
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
	RoomID      int64  `xorm:"pk autoincr" fake:"skip"`
	RoomName    string `fake:"{hipster.word}{address.street_suffix}"`
	Description string `fake:"{hacker.ingverb} {hacker.noun} {hacker.adjective}"`
	RoomType    RType  `fake:"skip"`
	WindowCount int64  `xorm:"null" fake:"skip"`
	IsSubRoom   bool
	PartOfRoom  int64  `xorm:"null" fake:"skip"`
	CreatedUnix int64  `xorm:"created"`
	UpdatedUnix int64  `xorm:"updated"`
	CurrentTemp int64  `xorm:"null"`
	HasLight    bool   `xorm:"-"`
	MainLight   Device `xorm:"-"`
	HasTemp     bool   `xorm:"-"`
	MainTemp    Device `xorm:"-"`
}

// LoadMainDevices loads the main light and temperature control variables of
// the room.
func (r *Room) LoadMainDevices() {
	devices, err := GetDevices()
	if err != nil {
		panic(err)
	}
	light, temp := false, false
	for _, l := range devices {
		if l.RoomID == r.RoomID && l.Type == Light && l.IsMainLight {
			r.MainLight = l
			r.HasLight = true
			light = true
		}
		if l.RoomID == r.RoomID && l.Type == TempControl {
			r.MainTemp = l
			r.HasTemp = true
			temp = true
		}
		if light && temp {
			break
		}
	}
}

// GetFakeRoom returns a new randomly generated Room. This is used for testing
// purposes.
func GetFakeRoom() (r *Room) {
	r = new(Room)
	gofakeit.Struct(r)
	r.RoomType = RType(gofakeit.Number(0, 9)) // This must match number of enums!
	r.WindowCount = int64(gofakeit.Number(0, 4))
	return
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
func GetRooms() (room []Room, err error) {
	err = engine.Find(&room)
	return
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

// DeleteRoom deletes a Room from the database.
func DeleteRoom(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Room{})
	return
}
