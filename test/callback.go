package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	// db, err = gorm.Open("postgres", "user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable")
	// db, err = gorm.Open("mysql", "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True")
	// db, err = gorm.Open("mssql", "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})
}
type User struct {
	Name string
	Sex string
}

//func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
//	tx.Debug().Model(u).Update("name", "fix")
//	return
//}


func (u *User) AfterUpdate(scope *gorm.Scope) (err error) {
	scope.DB().Debug().Model(u).Update("name", "fix")
	return
}

func main() {
	Init()
	db.Model(&User{}).Update("name", "rich")
}