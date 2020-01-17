package models

import (
	"errors"
)

type Statistic struct {
	StatID   int64 `xorm:"pk"`
	Datetime int64
	Powergen float64
	Powercon float64
}

func GetStat(id int64) (*Statistic, error) {
	s := new(Statistic)
	has, err := engine.ID(id).Get(s)
	if err != nil {
		return s, err
	} else if !has {
		return s, errors.New("Statistic does not exist")
	}
	return s, nil
}

func GetStats() (stat []Statistic) {
	engine.Find(&stat)
	return
}

func AddStat(s *Statistic) (err error) {
	_, err = engine.Insert(s)
	return
}

func HasStat(id int64) (has bool) {
	has, _ = engine.Get(&Statistic{StatID: id})
	return
}

func UpdateStat(d *Statistic) (err error) {
	_, err = engine.Id(d.StatID).Update(d)
	return
}

func UpdateStatCols(d *Statistic, cols ...string) (err error) {
	_, err = engine.ID(d.StatID).Cols(cols...).Update(d)
	return
}
