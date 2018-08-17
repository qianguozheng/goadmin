package model

import "github.com/jinzhu/gorm"

type DnsBogus struct {
	Id           int
	Type         int //1-project, 2-global
	Domain       string
	Ip           string
	Comment      string
	Status       int
	ProjectRefer int
}

func AddDnsBogus(db DnsBogus) bool {
	var i DnsBogus
	if DB.Debug().Where("type=? and domain=? and ip=? and project_refer=?",
		db.Type, db.Domain, db.Ip, db.ProjectRefer).Find(&i).RowsAffected > 0 {
		return false
	}
	DB.Create(&db)
	return true
}

func GetDnsBogusById(id int) DnsBogus {
	var i DnsBogus
	DB.Where("id=?", id).Find(&i)
	return i
}

func GetDnsBogusByProjectId(id int) []DnsBogus {
	var i []DnsBogus
	DB.Where("project_refer=?", id).Find(&i)
	return i
}

func DeleteDnsBogusById(id int) {
	DB.Model(&DnsBogus{}).Where("id=?", id).Delete(&DnsBogus{})
}

func UpdateDnsBogusById(dns DnsBogus) {
	DB.Debug().Model(&DnsBogus{}).Where("id=?", dns.Id).Update(&dns)
}

func GetAllDnsBogus() []DnsBogus {
	var i []DnsBogus
	DB.Model(&DnsBogus{}).Find(&i)
	return i
}

func (bogus *DnsBogus) AfterCreate(tx *gorm.DB) (err error) {
	devs := GetDeviceByProjectId(bogus.ProjectRefer)
	for _, dev := range devs {
		md5 := &Md5{
			DeviceRefer: dev.Id,
			Md5:         Md5sum(),
		}
		if tx.Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
			continue
		} else {
			tx.Create(md5)
		}
	}
	return nil
}

func (bogus *DnsBogus) AfterUpdate(tx *gorm.DB) (err error) {
	devs := GetDeviceByProjectId(bogus.ProjectRefer)
	for _, dev := range devs {
		md5 := &Md5{
			DeviceRefer: dev.Id,
			Md5:         Md5sum(),
		}
		if tx.Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
			continue
		} else {
			tx.Create(md5)
		}
	}
	return nil
}

func (bogus *DnsBogus) AfterDelete(tx *gorm.DB) (err error) {
	devs := GetDeviceByProjectId(bogus.ProjectRefer)
	for _, dev := range devs {
		md5 := &Md5{
			DeviceRefer: dev.Id,
			Md5:         Md5sum(),
		}
		if tx.Model(md5).Update("md5=? and device_refer=?", md5.Md5, md5.DeviceRefer).RowsAffected >= 1 {
			continue
		} else {
			tx.Create(md5)
		}
	}
	return nil
}
