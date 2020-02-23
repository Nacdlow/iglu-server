package models

import (
	"errors"
)

// ScheduleType is used to store the different types of schedules that the user can set
type ScheduleType int64

// The enumerations for the types of schedules that the user may set
const (
	Cron      = iota // 0
	TurnOnOff        // 1
)

// Schedule represents the parameters that a schedule may hold, such as the time that a device may turn on and off
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

// GetFakeSchedule creates a fake schedule
func GetFakeSchedule() (sch *Schedule) {
	sch = new(Schedule)
	// Some fake data generation here
	return
}

// GetSchedule gets the schedule
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

// AddSchedule makes a new schedule
func AddSchedule(sch *Schedule) (err error) {
	_, err = engine.Insert(sch)
	return
}

// UpdateSchedule updates the values of the schedule
func UpdateSchedule(d *Schedule) (err error) {
	_, err = engine.Id(d.ScheduleID).Update(d)
	return
}

// DeleteSchedule deletes the Schedule from the database
func DeleteSchedule(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Schedule{})
	return
}
