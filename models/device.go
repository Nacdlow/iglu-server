package models

import (
	"errors"
	"github.com/brianvoe/gofakeit/v4"
)

// DeviceType is the type of smart device.
type DeviceType int64

// DeviceType enums.
const (
	Light       = iota //0
	TempControl        //1
	TV                 //2
	Speaker            //3
)

// Device represents a smart home (Internet of Things) device, such as a light
// bulb, TV, temperature control (thermometer), etc.
type Device struct {
	DeviceID    int64      `xorm:"pk autoincr" fake:"skip"`
	RoomID      int64      `fake:"skip"`
	Type        DeviceType `fake:"skip"`
	Description string     `fake:"{lorem.word} {lorem.word} {lorem.word}"`
	Status      bool       `fake:"skip"`
	Temp        float64    `xorm:"null" fake:"skip"` // In Celcius
	Volume      int64      `xorm:"null" fake:"skip"`
	Brightness  int64      `xorm:"null" fake:"skip"`
	IsMainLight bool       `fake:"skip"` // Whether the light device is the room's main light source
}

// GetFakeDevice returns a new randomly created Device. This is used for
// testing purposes.
func GetFakeDevice() (d *Device) {
	d = new(Device)
	gofakeit.Struct(d)
	d.RoomID = int64(gofakeit.Number(0, 10))
	d.Type = DeviceType(gofakeit.Number(0, 3)) // This must match number of enums!
	d.Status = gofakeit.Bool()
	if d.Type == TempControl {
		d.Temp = gofakeit.Float64Range(18, 28)
	} else if d.Type == Speaker {
		d.Volume = int64(gofakeit.Number(0, 10))
	} else if d.Type == Light {
		d.Brightness = int64(gofakeit.Number(0, 10))
	}
	return
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

// DeleteDevice deletes a Device from the database.
func DeleteDevice(id int64) (err error) {
	_, err = engine.ID(id).Delete(&Device{})
	return
}
