package models

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
