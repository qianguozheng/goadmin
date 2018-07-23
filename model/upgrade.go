package model

import "github.com/jinzhu/gorm"

type Upgrade struct{
	gorm.Model
	ModelType int
	Md5 string
	Url string
	Status string
	Rcl bool
	Audit bool
}

func InsertFirmware(md5, url, status string, rcl, audit bool, model int) {
	firmware := &Upgrade{
		ModelType: model,
		Md5: md5,
		Url:url,
		Status: status,
		Rcl:rcl,
		Audit:audit,
	}
	DB.Create(firmware)
}

func DeleteFirmware(id int) {
	DB.Where("id=?", id).Delete(&Upgrade{})
}
