package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCascadeDeleteRoom(t *testing.T) {
	assert := assert.New(t)
	engine := SetupTestEngine()
	defer engine.Close()

	mainRoom := GetFakeRoom()
	assert.Nil(AddRoom(mainRoom))
	devices, err := GetDevices()
	assert.Nil(err, "Failed to get list of devices, when there are no devices.")
	assert.Equal(0, len(devices), "List of all devices must be zero, as we didn't add any.")

	const numOfDevices = 5
	for i := 0; i < 5; i++ {
		dev := GetFakeDevice()
		dev.RoomID = mainRoom.RoomID
		assert.Nil(AddDevice(dev))
	}

	devices, err = GetDevices()
	assert.Nil(err)
	assert.Equal(numOfDevices, len(devices), "List of devices should higher.")
	assert.Nil(DeleteRoom(mainRoom.RoomID))
	devices, err = GetDevices()
	assert.Nil(err)

	assert.False(HasRoom(mainRoom.RoomID), "Deleted room shouldn't exist")
	assert.Equal(0, len(devices), "List of all devices must be zero, as we cascade delete.")

}
