package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "goadmin.db")
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&User{}, &Upgrade{}, &Model{},
		&Device{}, &Qos{}, &Wan{}, &WanQos{}, &Ssid{})

	user := new(User)
	db.First(user, "name=?", "admin")
	if user.Password != "pass" {
		db.Create(&User{Name: "admin", Password: "pass"})
	}

	return db
}
