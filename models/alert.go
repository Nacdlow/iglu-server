package models

// Category is a ranking of how severe the aleart is
type Category int64

// Category enums from lowest to highest.
const (
	Low    = iota // 0
	Medium        // 1
	High          // 2
)

//Alert struct represents the Alert title, AlertCategory (severity), Alert Message,
// and also the time the Alert was sent to the user
type Alert struct {
	AlertID         int64 `xorm:"pk autoincr"`
	Time            int64 `xorm:"index"`
	Username        string
	AlertTitle      string
	AlertMessage    string
	AlertImportance Category
}
