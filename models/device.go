package models

type DeviceType int64

const (
	Light = iota
	TempControl
	TV
	Speaker
)

type Device struct {
	DeviceID    int64 `xorm:"pk autoincr"`
	RoomID      int64
	Type        DeviceType
	Description string
	Status      bool
	Temp        float64 `xorm:"null"`
}
