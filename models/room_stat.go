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

func GetRoomStats() (roomStats []RoomStat) {
	engine.Find(&roomStats)
	return
}

func AddRoomStat(r *RoomStat) (err error) {
	_, err = engine.Insert(r)
	return
}

func HasRoomStat(id int64) (has bool) {
	has, _ = engine.Get(&RoomStat{RStatID: id})
	return
}

func UpdateRoomStat(r *RoomStat) (err error) {
	_, err = engine.Id(r.RStatID).Update(r)
	return
}

func UpdateRoomStatCols(r *RoomStat, cols ...string) (err error) {
	_, err = engine.ID(r.RStatID).Cols(cols...).Update(r)
	return
}
