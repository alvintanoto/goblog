package database

import (
	"alvintanoto.id/blog/internal/database/connection"
	model "alvintanoto.id/blog/internal/model/database"
	"alvintanoto.id/blog/pkg/log"
)

type PostDB struct {
}

func (pdb PostDB) Insert(title string, content string, isPublic bool, userID int) (int, error) {
	db := new(connection.Postgresql).Get()

	post := &model.Post{
		Title:    title,
		Content:  content,
		IsPublic: isPublic,
		Base: model.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}

	result := db.Table("post").Create(&post)

	if result.Error != nil {
		log.Get().ErrorLog.Println(result.Error)
		return 0, result.Error
	}

	return post.ID, nil
}
