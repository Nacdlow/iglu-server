package models

// DeviceType is the type of smart device.
type DeviceType int64

const (
	Light = iota
	TempControl
	TV
	Speaker
)

// Device represents a smart home (Internet of Things) device, such as a light
// bulb, TV, temperature control (thermometer), etc.
type Device struct {
	DeviceID    int64 `xorm:"pk autoincr"`
	RoomID      int64
	Type        DeviceType
	Description string
	Status      bool
	Temp        float64 `xorm:"null"`
}
