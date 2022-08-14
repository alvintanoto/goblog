package database

import (
	"time"

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

func (pdb PostDB) GetHomePosts(userID *int) (*[]model.PostUser, error) {
	db := new(connection.Postgresql).Get()

	postUser := &[]model.PostUser{}
	s := []string{
		"post.id", "post.title", "post.content", "post.is_public", "post.created_at", "post.updated_at", "users.username",
	}

	if userID != nil {
		db.Table("post").
			Where("created_by = ?", userID).
			Or("is_public = ?", true).
			Select(s).
			Order("post.updated_at desc").
			Limit(50).
			Joins("left join users on post.created_by = users.id").Scan(&postUser)
	} else {
		db.Table("post").
			Where("is_public = ?", true).
			Select(s).
			Order("post.updated_at desc").
			Limit(50).
			Joins("left join users on post.created_by = users.id").Scan(&postUser)
	}

	return postUser, nil
}

func (pdb PostDB) Get(id int) (*model.PostUser, error) {
	db := new(connection.Postgresql).Get()

	postUser := &model.PostUser{}
	s := []string{
		"post.id", "post.title", "post.content", "post.is_public", "post.created_at", "post.updated_at", "users.username", "post.created_by",
	}

	db.Table("post").
		Where("post.id = ?", id).
		Select(s).
		Order("post.updated_at desc").
		Joins("left join users on post.created_by = users.id").Scan(&postUser)

	return postUser, nil
}

func (pdb PostDB) Update(title string, content string, isPublic bool, postID int, userID int) (int, error) {
	db := new(connection.Postgresql).Get()

	post := &model.Post{
		Title:    title,
		Content:  content,
		IsPublic: isPublic,
		IsEdited: true,
		Base: model.Base{
			UpdatedBy: userID,
			UpdatedAt: time.Now(),
		},
	}

	result := db.Table("post").Where("id = ?", postID).Updates(post)

	if result.Error != nil {
		log.Get().ErrorLog.Println(result.Error)
		return 0, result.Error
	}

	return post.ID, nil
}
