package models

type Statistic struct {
	Datetime int64 `xorm:"pk"`
	Powergen float64
	Powercon float64
	StatID   int64
}
