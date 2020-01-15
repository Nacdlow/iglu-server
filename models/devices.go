package models

type DeviceType int64

const (
	Light = iota
	TempControl
	TV
	Speaker
)

type Devices struct {
	DeviceID    int64 `xorm:"pk"`
	RoomID      int64
	Type        DeviceType
	Description string
	Status      bool
	Temp        float64 `xorm:"null"`
}
