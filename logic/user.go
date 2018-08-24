package logic

import (
	"github.com/qianguozheng/goadmin/model"
)

type UserLogic struct{}

var DefaultUser = UserLogic{}

func (u UserLogic) FindCurrentUser(name string) interface{} {
	var user model.User
	model.DB.Where("name=?", name).Find(&user)
	return user
}

func (u UserLogic) UpdateUser(iuser interface{}) {
	user := iuser.(model.User)
	model.DB.Debug().Model(&model.User{}).Where("name=?", user.Name).Update(&user)
}
