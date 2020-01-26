package forms

// AddRoomForm represents a form to add a new Room to the house.
type AddRoomForm struct {
	RoomName     string `form:"room_name" binding:"Required"`
	Description  string `form:"description" binding:"Required"`
	RoomType     int64  `form:"room_type"` // We don't use required because that would require number >0
	WindowsCount int64  `form:"windows_count"`
	PartOfRoom   int64  `form:"part_of_room"`
}

// AddDeviceForm represents a form to add a new Device to a room.
type AddDeviceForm struct {
	RoomID      int64  `form:"room_id"`
	DeviceType  int64  `form:"device_type"`
	Description string `form:"description" binding:"Required"`
	IsMainLight bool   `form:"is_main_light_source"`
}
