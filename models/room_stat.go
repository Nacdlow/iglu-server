package models

import (
	"errors"

	"github.com/brianvoe/gofakeit/v4"
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
	CreatedUnix int64 `xorm:"created"`
	UpdatedUnix int64 `xorm:"updated"`
}

// GetFakeRoomStat returns a randomly generated room statistic. This is used
// for testing purposes.
func GetFakeRoomStat() (r *RoomStat) {
	// We will not use gofakeit.Struct as all fields are using ranges, which is
	// not possible to do with gofakeit's struct tag support.
	r = new(RoomStat)
	r.StatID = int64(gofakeit.Number(0, 9))
	r.RoomID = int64(gofakeit.Number(0, 9))
	r.Temperature = gofakeit.Float64Range(13, 26)
	r.Humidity = gofakeit.Float64Range(30, 75)
	r.OpenWindows = int64(gofakeit.Number(0, 3))
	return
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
func GetRoomStats() (roomStats []RoomStat, err error) {
	err = engine.Find(&roomStats)
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

// DeleteRoomStat deletes a RoomStat from the database.
func DeleteRoomStat(id int64) (err error) {
	_, err = engine.ID(id).Delete(&RoomStat{})
	return
}
