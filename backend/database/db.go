package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Connect to the database
func ConnectDB() (*gorm.DB, error) {
	var err error
	dsn := "host=localhost user=postgres password=Gurung67 dbname=Portfolio port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}

	return db, nil
}
