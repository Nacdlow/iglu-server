package models

// Category is a ranking of how severe the aleart is
type Category int64

// Category enums from lowest to highest.
const (
	Low    = iota // 0
	Medium        // 1
	High          // 2
)

// Alert struct represents the Alert title, AlertCategory (severity), Alert Message,
// and also the time the Alert was sent to the user
type Alert struct {
	AlertID    int64    `xorm:"pk autoincr" json:"id" xml:"id,attr"`
	Time       int64    `json:"time" xml:"time"`
	Username   string   `xorm:"index" json:"username" xml:"username,attr"`
	Title      string   `json:"title" xml:"title"`
	Message    string   `json:"message" xml:"message"`
	Importance Category `json:"importance" xml:"importance"`
}
