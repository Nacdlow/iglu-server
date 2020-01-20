package models

import (
	"errors"
)

// Statistic represents a statistic log at a period of time.
type Statistic struct {
	StatID   int64 `xorm:"pk"`
	Datetime int64
	Powergen float64
	Powercon float64
}

// GetStat gets a Statistic based on its ID from the database.
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

// GetStats returns an array of all Statistics from the database.
func GetStats() (stat []Statistic) {
	engine.Find(&stat)
	return
}

// AddStat adds a Statistic in the database.
func AddStat(s *Statistic) (err error) {
	_, err = engine.Insert(s)
	return
}

// HasStat returns whether an Statistic is in the database or not.
func HasStat(id int64) (has bool) {
	has, _ = engine.Get(&Statistic{StatID: id})
	return
}

// UpdateStat updates an Statistic in the database.
func UpdateStat(d *Statistic) (err error) {
	_, err = engine.Id(d.StatID).Update(d)
	return
}

// UpdateStatCols will update the columns of an item even if they are empty.
func UpdateStatCols(d *Statistic, cols ...string) (err error) {
	_, err = engine.ID(d.StatID).Cols(cols...).Update(d)
	return
}

// DeleteStat deletes a Statistic from the database.
func DeleteStat(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Statistic{})
	return
}
