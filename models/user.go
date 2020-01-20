package models

import (
	"errors"
)

// UserRole is the role of a user.
type UserRole int64

// UserRole enums.
const (
	NormalRole = iota
	AdminRole
	GuestRole
	RestrictedRole
)

// User is a smart home system user whom may interact with the system.
type User struct {
	Username  string `xorm:"pk"`
	Password  string
	FirstName string
	LastName  string
	Role      UserRole
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
	return u, nil
}

// GetUsers returns an array of all Users from the database.
func GetUsers() (users []User) {
	engine.Find(&users)
	return users
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
