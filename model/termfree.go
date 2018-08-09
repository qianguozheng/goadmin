package model

type TermFree struct {
	Id      int `gorm:"AUTO_INCREMENT"`
	Mac     string
	Comment string
}

func AddTermFree(mac, comment string) {
	term := TermFree{
		Mac:     mac,
		Comment: comment,
	}
	if DB.Model(&TermFree{}).Where("mac=? and comment=?", mac, comment).Find(&term).RowsAffected > 0 {
		return
	}
	DB.Debug().Create(&term)
}

func DeleteTermFreeById(id int) {
	DB.Debug().Where("id=?", id).Delete(&TermFree{})
}
func DeleteTermFreeByMac(mac string) {
	DB.Debug().Where("mac=?", mac).Delete(&TermFree{})
}

func CheckTermFreeByMac(mac string) bool {
	var term TermFree
	if DB.Model(&TermFree{}).Where("mac=?", mac).Find(&term).RowsAffected > 0 {
		return true
	}
	return false
}

func GetAllTermFree() []TermFree {
	var tf []TermFree
	DB.Debug().Find(&tf)
	return tf
}

func GetTotalTermFreeNum() int {
	var count int
	DB.Debug().Table("term_frees").Count(&count)
	return count
}

func ListPageNoDeviceTermFree(pageNo, pageSize int) []TermFree {
	var dev []TermFree

	DB.Debug().Order("id asc").Find(&dev).Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&dev)

	//	fmt.Println("offset len dev:", len(dev))
	//	for k, v := range dev {
	//		fmt.Println("k=", k, ", id=", v.Id)
	//	}

	return dev
}
