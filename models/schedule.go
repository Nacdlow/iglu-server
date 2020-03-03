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
	ScheduleID  int64        `xorm:"pk autoincr" json:"id" xml:"id,attr"`
	DeviceID    string       `json:"deviceID" xml:"device_id,attr"`
	Cron        string       `xorm:"null" json:"cron,omitempty" xml:"schedule>cron,omitempty"`
	Start       int64        `json:"start" xml:"schedule>start"` // Start time for task
	End         int64        `json:"end" xml:"schedule>end"`     // End time for task
	Username    string       `json:"username" xml:"username"`    // Username of
	Type        ScheduleType `json:"type" xml:"type,attr"`
	Title       string       `json:"title" xml:"title"`
	Description string       `json:"description" xml:"description"`
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
