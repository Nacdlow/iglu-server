package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFavRooms(t *testing.T) {
	assert := assert.New(t)
	engine := SetupTestEngine()
	defer engine.Close()
	// Add rooms
	assert.Nil(AddRoom(&Room{
		RoomName: "Main living",
		RoomType: LoungeRoom,
	}))
	assert.Nil(AddRoom(&Room{
		RoomName: "Main kitchen",
		RoomType: KitchenRoom,
	}))
	assert.Nil(AddRoom(&Room{
		RoomName: "Outdoor garage",
		RoomType: GarageRoom,
	}))
	// Add user
	assert.Nil(AddUser(&User{
		Username:     "az40",
		FirstName:    "Alakbar",
		LastName:     "Zeynalzade",
		Role:         AdminRole,
		FontSize:     "medium",
		FavRoomsList: []int64{1, 2},
	}))
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
