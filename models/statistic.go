package models

type Statistic struct {
	StatID   int64 `xorm:"pk"`
	Datetime int64
	Powergen float64
	Powercon float64
}
