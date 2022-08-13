package database

import (
	"strings"

	"alvintanoto.id/blog/internal/database/connection"
	model "alvintanoto.id/blog/internal/model/database"
	"alvintanoto.id/blog/pkg/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDB struct {
}

func (udb UserDB) Insert(username string, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	db := new(connection.Postgresql).Get()

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	result := db.Table("user").Create(&user)

	if result.Error != nil {
		log.Get().ErrorLog.Println(result.Error)

		if strings.HasPrefix(result.Error.Error(), "ERROR: duplicate key value violates unique constraint") {
			return 0, connection.ErrConflictData
		}

		return 0, result.Error
	}

	return user.ID, nil
}

func (udb UserDB) Authenticate(username string, password string) (int, error) {
	db := new(connection.Postgresql).Get()

	user := &model.User{}

	result := db.Table("user").Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, connection.ErrInvalidCredential
		}

		return 0, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (udb UserDB) Get(userId int) (*model.User, error) {
	db := new(connection.Postgresql).Get()

	user := &model.User{}

	result := db.Table("user").Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
