package model

import (
	"fmt"

	"errors"

	"github.com/jinzhu/gorm"
	//"time"
)

//Device Group
type Project struct {
	gorm.Model
	//ID        int
	Name    string
	Comment string
	//CreatedAt time.Time
}

func GetProjects() []Project {
	var prj []Project
	DB.Find(&prj)
	return prj
}

func GetTotalProjectNum() int {
	var count int
	DB.Table("projects").Count(&count)
	return count
}

func ListPageNoProject(pageNo, pageSize int) []Project {
	var prj []Project

	DB.Order("id asc").Find(&prj).Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&prj)

	// fmt.Println("offset len dev:", len(prj))
	// for k, v := range prj {
	// 	fmt.Println("k=", k, ", id=", v.ID)
	// }

	return prj
}

func GetProjectById(id int) Project {
	var prj Project
	DB.Where("id=?", id).Find(&prj)
	return prj
}

func InsertProject(name, comment string) error {
	prj := &Project{
		Name:    name,
		Comment: comment,
	}
	var p Project
	if DB.Where("name=? and comment=?", name, comment).Find(&p).RowsAffected >= 1 {
		return errors.New("Alread exist name")
	}
	DB.Create(prj)
	return nil
}

func UpdateProjectBy(prj Project) {
	if DB.Model(&Project{}).Where("id=?", prj.ID).Update(&prj).RowsAffected >= 1 {
		fmt.Println("prject id updated")
		return
	}
	DB.Create(&prj)
}

func DeleteProjectById(id int) {
	DB.Where("id=?", id).Delete(&Project{})
}
