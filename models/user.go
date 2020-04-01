package models

import (
	"errors"
	"log"

	"github.com/brianvoe/gofakeit/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserRole is the role of a user.
type UserRole int64

// UserRole enums.
const (
	NormalRole     = iota // 0
	AdminRole             // 1
	GuestRole             // 2
	RestrictedRole        // 3
)

// User is a smart home system user whom may interact with the system.
type User struct {
	Username      string            `xorm:"pk" fake:"{animal.petname}###" json:"username" xml:"username,attr"`
	Password      string            `fake:"skip" json:"password" xml:"password"`
	FirstName     string            `fake:"{person.first}" json:"firstName" xml:"name>first"`
	LastName      string            `fake:"{person.last}" json:"lastName" xml:"name>last"`
	Role          UserRole          `fake:"skip" json:"role" xml:"role"`
	FavRoomsList  []int64           `fake:"skip" json:"favRoomsList" xml:"fav_rooms_list"`
	FavRooms      []*Room           `xorm:"-" fake:"skip" json:"-" xml:"-"` // This means do not store this in the DB.
	CreatedUnix   int64             `xorm:"created" json:"createdUnix" xml:"timestamps>created_unix"`
	UpdatedUnix   int64             `xorm:"updated" json:"updatedUnix" xml:"timestamps>updated_unix"`
	FontSize      string            `json:"fontSize" xml:"font_size"`
	PluginKVStore map[string]string `xorm:"null" json:"pluginKVStore" xml:"plugin_key_value_store"` // must check if null first
	Avatar        string            `json:"avatar" xml:"avatar"`
}

// GetFakeUser returns a new randomly generated User. This is used for testing
// purposes.
func GetFakeUser() (u *User) {
	u = new(User)
	gofakeit.Struct(u)
	newPass := gofakeit.Password(true, true, true, true, false, 12)
	pass, err := bcrypt.GenerateFromPassword([]byte(newPass), 4)
	if err != nil {
		panic(err)
	}
	u.Password = string(pass)
	log.Printf("Fake user %s has password of: %s\n", u.Username, newPass)
	u.Role = UserRole(gofakeit.Number(0, 3)) // This must match the number of enums!
	for i := 0; i < gofakeit.Number(0, 4); i++ {
		u.FavRoomsList = append(u.FavRoomsList, int64(gofakeit.Number(0, 9)))
	}
	return
}

// GetUser gets a User based on its ID from the database.
func GetUser(user string) (*User, error) {
	u := new(User)
	has, err := engine.ID(user).Get(u)
	if err != nil {
		return u, err
	} else if !has {
		return u, errors.New("User does not exist")
	}
	// Load favourite rooms
	for _, i := range u.FavRoomsList {
		room, err := GetRoom(i)
		if err == nil {
			u.FavRooms = append(u.FavRooms, room)
		}
	}
	return u, nil
}

// GetUsers returns an array of all Users from the database.
func GetUsers() (users []User, err error) {
	err = engine.Find(&users)
	return
}

// AddUser adds a User in the database.
func AddUser(u *User) (err error) {

	_, err = engine.Insert(u)
	return err
}

// HasUser returns whether a User is in the database or not.
func HasUser(user string) (has bool) {
	has, _ = engine.Get(&User{Username: user})
	return has
}

// UpdateUser updates an User in the database.
func UpdateUser(u *User) (err error) {
	_, err = engine.Id(u.Username).Update(u)
	return
}

// UpdateUserCols will update the columns of an item even if they are empty.
func UpdateUserCols(u *User, cols ...string) error {
	_, err := engine.Id(u.Username).Cols(cols...).Update(u)
	return err
}

// DeleteUser deletes a User from the database.
func DeleteUser(user string) (err error) {
	_, err = engine.ID(user).Delete(&User{})
	return
}
