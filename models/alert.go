package models

import (
	"errors"
)

// AlertCategory is a ranking of how severe the aleart is.
type AlertCategory int64

// AlertCategory enums from lowest to highest.
const (
	Low    = iota // 0
	Medium        // 1
	High          // 2
)

// Alert struct represents the Alert title, AlertCategory (severity), Alert Message,
// and also the time the Alert was sent to the user
type Alert struct {
	AlertID    int64         `xorm:"pk autoincr" json:"id" xml:"id,attr"`
	Time       int64         `json:"time" xml:"time"`
	Username   string        `xorm:"index" json:"username" xml:"username,attr"`
	Title      string        `json:"title" xml:"title"`
	Message    string        `json:"message" xml:"message"`
	Importance AlertCategory `json:"importance" xml:"importance"`
}

// GetAlert gets a alert based on its ID from the database.
func GetAlert(id int64) (*Alert, error) {
	d := new(Alert)
	has, err := engine.ID(id).Get(d)
	if err != nil {
		return d, err
	} else if !has {
		return d, errors.New("Alert does not exist")
	}
	return d, nil
}

// GetAlerts returns an array of all alerts from the database.
func GetAlerts() (alerts []Alert, err error) {
	err = engine.Find(&alerts)
	return
}

// AddAlert adds a Alert in the database.
func AddAlert(d *Alert) (err error) {
	_, err = engine.Insert(d)
	return
}

// HasAlert returns whether a alert is in the database or not.
func HasAlert(id int64) (has bool) {
	has, _ = engine.Get(&Alert{AlertID: id})
	return
}

// UpdateAlert updates a Alert in the database.
func UpdateAlert(d *Alert) (err error) {
	_, err = engine.Id(d.AlertID).Update(d)
	return
}

// UpdateAlertCols will update the columns of an item even if they are empty.
func UpdateAlertCols(d *Alert, cols ...string) (err error) {
	_, err = engine.ID(d.AlertID).Cols(cols...).Update(d)
	return
}

// DeleteAlert deletes a Alert from the database.
func DeleteAlert(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Alert{})
	return
}
