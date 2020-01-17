package models

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
