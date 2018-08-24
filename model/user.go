package model

import (
	"fmt"
)

type User struct {
	Name     string
	Password string
	Cookie   string
	Phone    string
	Email    string
	NickName string
}

func GetUser(name, pass string) bool {
	var user User
	DB.Where("name=? AND password=?", name, pass).Find(&user)
	if user.Name == name && user.Password == pass {
		return true
	}
	return false
}

func GetCookie(cookie string) bool {
	var user User
	DB.Where("cookie=?", cookie).Find(&user)
	fmt.Println("user=", user.Name, " pass=", user.Password, "cookie=", user.Cookie)
	if user.Name != "" && user.Password != "" {
		return true
	} else {
		return false
	}
}
func SetCookie(user, cookie string) {
	userName := &User{Name: user}
	DB.Model(userName).Update("cookie", cookie)
}

func Logout(cookie string) {
	//DB.Where("cookie = ?", cookie).Delete(&User{})
	DB.Model(&User{}).Where("cookie=?", cookie).Update("cookie", "")
}
