package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserFavRooms(t *testing.T) {
	assert := assert.New(t)
	engine := SetupTestEngine()
	defer engine.Close()
	// Add rooms
	AddRoom(&Room{
		RoomName:    "Main living",
		Description: "Our family living room",
		RoomType:    LoungeRoom,
	})
	AddRoom(&Room{
		RoomName:    "Main kitchen",
		Description: "Our family kitchen",
		RoomType:    KitchenRoom,
	})
	AddRoom(&Room{
		RoomName:    "Outdoor garage",
		Description: "The outdoor garage",
		RoomType:    GarageRoom,
	})
	// Add user
	AddUser(&User{
		Username:     "az40",
		FirstName:    "Alakbar",
		LastName:     "Zeynalzade",
		Role:         AdminRole,
		FavRoomsList: []int64{1, 2},
	})
	// Getting the user
	user, err := GetUser("az40")
	assert.Nil(err)
	assert.Equal("Alakbar", user.FirstName, "First name must match!")
	assert.Equal("Zeynalzade", user.LastName, "Last name must match!")
	assert.Equal(AdminRole, int(user.Role), "Role must match!")
	assert.Equal(2, len(user.FavRoomsList), "Length of favourite rooms list must be 2!")
	assert.Equal(2, len(user.FavRooms), "Length of favourite rooms must be 2!")
	assert.False(HasUser("ha82"), "The user should not be in the database!")
	assert.False(HasRoom(10), "The room should not be in the database!")
}
