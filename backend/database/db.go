package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/postgres"
)

var db *gorm.DB

// Connect to the database
func ConnectDB() (*gorm.DB, error) {
	var err error
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
