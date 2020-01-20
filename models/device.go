package models

import (
	"errors"
)

// DeviceType is the type of smart device.
type DeviceType int64

// DeviceType enums.
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

// GetDevice gets a device based on its ID from the database.
func GetDevice(id int64) (*Device, error) {
	d := new(Device)
	has, err := engine.ID(id).Get(d)
	if err != nil {
		return d, err
	} else if !has {
		return d, errors.New("Device does not exist")
	}
	return d, nil
}

// GetDevices returns an array of all devices from the database.
func GetDevices() (devices []Device) {
	engine.Find(&devices)
	return
}

// AddDevice adds a Device in the database.
func AddDevice(d *Device) (err error) {
	_, err = engine.Insert(d)
	return
}

// HasDevice returns whether a device is in the database or not.
func HasDevice(id int64) (has bool) {
	has, _ = engine.Get(&Device{DeviceID: id})
	return
}

// UpdateDevice updates a Device in the database.
func UpdateDevice(d *Device) (err error) {
	_, err = engine.Id(d.DeviceID).Update(d)
	return
}

// UpdateDeviceCols will update the columns of an item even if they are empty.
func UpdateDeviceCols(d *Device, cols ...string) (err error) {
	_, err = engine.ID(d.DeviceID).Cols(cols...).Update(d)
	return
}
