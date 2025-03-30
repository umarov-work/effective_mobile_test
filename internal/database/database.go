package database

import (
	"effective_mobile_test/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB(host string, port string, user string, password string, dbname string) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	dsn := fmt.Sprintf("host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" port=" + port +
		" sslmode=disable")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Person{})
	if err != nil {
		return err
	}

	return nil
}
