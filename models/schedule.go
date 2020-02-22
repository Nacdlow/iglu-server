package models

import (
	"errors"
)

// ScheduleType some text
type ScheduleType int64

// Some text
const (
	Duration  = iota // 0
	TurnOnOff        // 1
)

// Schedule some text
type Schedule struct {
	ScheduleID  int64 `xorm:"pk autoincr"`
	DeviceID    string
	Cron        string `xorm:"null"`
	Start       int64  // Start time for task
	End         int64  // End time for task
	Username    string // Username of
	Type        ScheduleType
	Title       string
	Description string
}

// GetFakeSchedule some text
func GetFakeSchedule() (sch *Schedule) {
	sch = new(Schedule)
	// Some fake data generation here
	return
}

// GetSchedule some text
func GetSchedule(id int64) (*Schedule, error) {
	sch := new(Schedule)
	has, err := engine.ID(id).Get(sch)
	if err != nil {
		return sch, err
	} else if !has {
		return sch, errors.New("Schedule does not exist")
	}
	return sch, nil
}

// AddSchedule some text
func AddSchedule(sch *Schedule) (err error) {
	_, err = engine.Insert(sch)
	return
}

// UpdateSchedule updates an Schedule in the database
func UpdateSchedule(d *Schedule) (err error) {
	_, err = engine.Id(d.ScheduleID).Update(d)
	return
}

// DeleteSchedule deletes a Schedule from the database.
func DeleteSchedule(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Schedule{})
	return
}
