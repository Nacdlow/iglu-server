package models

type UserType int64

const (
    Normal = iota
    Admin
    Guest
    Restricted
)

type User struct {
    Username  string `xorm:"pk"`
    Password  string
    FirstName string
    LastName  string
    Type      UserType
}
