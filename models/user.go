package models

import (
	"errors"
)

// UserRole is the role of a user.
type UserRole int64

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

func GetUsers() (users []User) {
	engine.Find(&users)
	return users
}

func AddUser(u *User) (err error) {
	_, err = engine.Insert(u)
	return err
}

func HasUser(user string) (has bool) {
	has, _ = engine.Get(&User{Username: user})
	return has
}

func UpdateUser(u *User) (err error) {
	_, err = engine.Id(u.Username).Update(u)
	return
}

func UpdateUserCols(u *User, cols ...string) error {
	_, err := engine.Id(u.Username).Cols(cols...).Update(u)
	return err
}
