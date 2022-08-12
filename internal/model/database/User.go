package model

type User struct {
	ID       int
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}
