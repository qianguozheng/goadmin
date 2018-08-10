package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func InitDB(path string) *gorm.DB {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&User{}, &Upgrade{}, &Model{},
		&Device{}, &Qos{}, &Wan{}, &WanQos{}, &Ssid{}, &Project{},
		&Md5{}, &TermFree{}, &TrustIps{}, &TrustDomains{})

	user := new(User)
	db.First(user, "name=?", "admin")
	if user.Password != "pass" {
		db.Create(&User{Name: "admin", Password: "pass"})
	}

	return db
}
