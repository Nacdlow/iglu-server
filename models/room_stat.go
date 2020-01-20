package models

import (
	"errors"
)

// RoomStat represents a room statistic, which is a part of a larger Statistic.
type RoomStat struct {
	RStatID     int64 `xorm:"pk autoincr"`
	StatID      int64 // the larger Statistic ID
	RoomID      int64
	Room        Room `xorm:"-"`
	Temperature float64
	Humidity    float64
	OpenWindows int64
}

// GetRoomStat gets a RoomStat based on its ID from the database.
func GetRoomStat(id int64) (*RoomStat, error) {
	r := new(RoomStat)
	has, err := engine.ID(id).Get(r)
	if err != nil {
		return r, err
	} else if !has {
		return r, errors.New("RoomStat does not exist")
	}
	return r, nil
}

// GetRoomStats returns an array of all RoomStats from the database.
func GetRoomStats() (roomStats []RoomStat) {
	engine.Find(&roomStats)
	return
}

// AddRoomStat adds an RoomStat in the database.
func AddRoomStat(r *RoomStat) (err error) {
	_, err = engine.Insert(r)
	return
}

// HasRoomStat returns whether an RoomStat is in the database or not.
func HasRoomStat(id int64) (has bool) {
	has, _ = engine.Get(&RoomStat{RStatID: id})
	return
}

// UpdateRoomStat updates an RoomStat in the database.
func UpdateRoomStat(r *RoomStat) (err error) {
	_, err = engine.Id(r.RStatID).Update(r)
	return
}

// UpdateRoomStatCols will update the columns of an item even if they are empty.
func UpdateRoomStatCols(r *RoomStat, cols ...string) (err error) {
	_, err = engine.ID(r.RStatID).Cols(cols...).Update(r)
	return
}
