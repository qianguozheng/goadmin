package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	//"time"
)

type Upgrade struct {
	gorm.Model
	//ID        int `gorm:"AUTO_INCREMENT"`
	ModelType int
	ModelName string
	Md5       string
	Url       string
	Status    string
	Rcl       bool
	Audit     bool
	Version   string
	Remark    string
	//CreatedAt time.Time
}

type Model struct {
	Id        int
	ModelName string
	ModelType int
}

func GetAllModels() []Model {
	var models []Model
	DB.Find(&models)
	return models
}

func InitModels() {
	model := Model{
		ModelName: "WR_9344_A",
		ModelType: 3,
	}

	if DB.Find(&model).RowsAffected > 0 {
		fmt.Println("already in database ar9344")
		return
	}

	DB.Create(&model)
	DB.Create(&Model{ModelName: "WR_9344_B", ModelType: 4})
	DB.Create(&Model{ModelName: "WR_9341_A", ModelType: 5})
	DB.Create(&Model{ModelName: "GW500", ModelType: 6})
	DB.Create(&Model{ModelName: "WR_7620_A", ModelType: 2})
	DB.Create(&Model{ModelName: "WR_9531_A", ModelType: 7})
	DB.Create(&Model{ModelName: "WR_9531_C", ModelType: 8})
}

func InitUpgrade() {

	firmware := Upgrade{
		ModelType: 3,
	}
	//TODO: remove test data
	if DB.Find(&firmware).RowsAffected > 0 {
		fmt.Println("Test firmware data inserted")
		return
	}
	InsertFirmware("", "123", "xxxx", "on", "ar9344", "1.1", false, true, 3)
	InsertFirmware("", "123", "xxxx", "on", "ar9344-b", "1.1", false, true, 4)
	InsertFirmware("", "123", "xxxx", "on", "ar9341-a", "1.1", false, true, 5)
	InsertFirmware("", "123", "xxxx", "on", "gw500", "1.1", true, false, 6)
	InsertFirmware("", "123", "xxxx", "on", "7620a", "1.1", true, false, 2)
	InsertFirmware("", "123", "xxxx", "on", "9531-a", "1.1", true, false, 7)
	InsertFirmware("", "123", "xxxx", "on", "9531-c", "1.1", true, false, 8)
}
func InsertFirmware(remark, md5, url, status, name, ver string, rcl, audit bool, model int) {
	firmware := &Upgrade{
		ModelType: model,
		ModelName: name,
		Md5:       md5,
		Version:   ver,
		Url:       url,
		Status:    status,
		Rcl:       rcl,
		Audit:     audit,
		Remark:    remark,
	}
	DB.Create(firmware)
}

func UpdateFirmware(remark, md5, url, status, name, ver string, rcl, audit bool, model, id int) {
	firmware := &Upgrade{
		ModelType: model,
		ModelName: name,
		Md5:       md5,
		Version:   ver,
		Url:       url,
		Status:    status,
		Rcl:       rcl,
		Audit:     audit,
		Remark:    remark,
	}
	DB.Model(&Upgrade{}).Where("model_type=? and id=?", model, id).Update(firmware)
}

func GetAllFirmwares() []Upgrade {
	var firmware []Upgrade
	DB.Find(&firmware)
	return firmware
}

func GetFirmwareById(id, model int) Upgrade {
	var u Upgrade
	DB.Where("id=? and model_type =?", id, model).Find(&u)
	return u
}

func DeleteFirmwareById(id, model int) {
	DB.Where("id=? and model_type =?", id, model).Delete(&Upgrade{})
}
func DeleteFirmwareByModel(model int) {
	DB.Where("model_type=?", model).Delete(&Upgrade{})
}
