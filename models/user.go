package models

type UserRole int64

const (
	NormalRole = iota
	AdminRole
	GuestRole
	RestrictedRole
)

type User struct {
	Username  string `xorm:"pk"`
	Password  string
	FirstName string
	LastName  string
	Role      UserRole
}
