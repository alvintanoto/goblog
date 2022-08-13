package model

type Post struct {
	ID       int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
	IsPublic bool   `gorm:"column:is_public"`
	Base
}

type PostUser struct {
	ID       int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
	IsPublic bool   `gorm:"column:is_public"`
	Username string
	Base
}
