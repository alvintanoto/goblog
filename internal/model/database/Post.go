package model

type Post struct {
	ID       int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
	IsPublic bool   `gorm:"column:is_public"`
	IsEdited bool   `gorm:"column:is_edited"`
	Base
}

type PostUser struct {
	ID       int    `gorm:"column:id"`
	Title    string `gorm:"column:title"`
	Content  string `gorm:"column:content"`
	IsPublic bool   `gorm:"column:is_public"`
	IsEdited bool   `gorm:"column:is_edited"`
	Username string
	Base
}
