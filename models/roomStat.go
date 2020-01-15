package models

type RoomStat struct {
	RstatID int64 `xorm:"pk autoincr"`
	RoomID int64 
	TheRoom Room `xorm:"-"`
	Temperature float64
	Humidity float64
	OpenWindows int64 
}