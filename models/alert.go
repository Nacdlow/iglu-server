package models

import (
	"time"
)

// Category is a ranking of how severe the aleart is
type Category int64

// Category enums from lowest to highest.
const (
	Low 			// 0
	Medium 			// 1
	High        	// 2
)

//Alert struct represents the Aleart title, AlertCategory (severity), Alert Message,
// and also the time the Alert was sent to the user
type Alert struct {
	AlertID int64	`xorm:"pk autoincr"`
	AlertCategory Category
	Time int64
	AlertTitle string
	AlertMessage string
}