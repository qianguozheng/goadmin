package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//Test usage, standalone db goauth.db
type AuthUser struct {
	Id       int `gorm:"AUTO_INCREMENT"`
	Name     string
	Password string
	Mac      string
	PhoneNo  string
	Valid    int
	LoginAt  time.Time
	Cookie   string
	//TODO: low priority -Account, money, buy time, to be implement
}

var DB2 *gorm.DB

func InitDB2() *gorm.DB {
	db, err := gorm.Open("sqlite3", "goauth.db")
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&AuthUser{})

	user := new(User)
	db.First(user, "name=?", "admin")
	if user.Password != "pass" {
		db.Create(&User{Name: "admin", Password: "pass"})
	}

	return db
}

func AddAuthUser(name string) {
	au := AuthUser{
		Name:  name,
		Valid: 1440,
	}
	if DB2.Debug().Where("name=?", name).Find(&au).RowsAffected > 0 {
		fmt.Println("alread inserted:", name)
		return
	}
	DB2.Debug().Create(&au)
}

func GetAuthUserByName(name string) AuthUser {
	var au AuthUser
	DB2.Debug().Where("name=?", name).Find(&au)
	return au
}

func GetAuthUserById(id int) AuthUser {
	var au AuthUser
	DB2.Debug().Where("id=?", id).Find(&au)
	return au
}

func DeleteUserByName(name string) {
	DB2.Debug().Model(&AuthUser{}).Where("name=?", name).Delete(&AuthUser{})
}

func DeleteUserById(id int) {
	DB2.Debug().Model(&AuthUser{}).Where("id=?", id).Delete(&AuthUser{})
}

func UpdateUser(au AuthUser) {
	DB2.Debug().Model(&AuthUser{}).Update(&au)
}

func AuthUserTest() {
	AddAuthUser("hello")
	AddAuthUser("world")
}
