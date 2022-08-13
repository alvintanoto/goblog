package connection

import (
	"errors"
	"fmt"
	"time"

	"alvintanoto.id/blog/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgresql struct{}

var db *gorm.DB

var (
	ErrConflictData      = errors.New("violates_unique_constraint")
	ErrRecordNotFound    = errors.New("record_not_found")
	ErrInvalidCredential = errors.New("invalid_credentials")
)

func (p Postgresql) Init(dsn string) {
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		log.Get().ErrorLog.Fatal("Postgres Connection", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
}

func (g Postgresql) Get() *gorm.DB {
	return db
}
