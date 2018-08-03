package main

import (
	"fmt"

	"../model"
)

func CountDevice() {
	var count int
	model.DB.Debug().Table("devices").Count(&count)
	fmt.Println("count device table:", count)
}

func LimitDevice(limit int) {
	var dev []model.Device
	model.DB.Debug().Limit(limit).Order("id asc").Find(&dev)
	fmt.Println("limit len dev:", len(dev))
	for k, v := range dev {
		fmt.Println("k=", k, ", id=", v.Id)
	}
}

func OffsetDevice(offset int) {
	var dev []model.Device

	model.DB.Debug().Order("id asc").Find(&dev).Limit(4).Offset(2).Find(&dev)
	fmt.Println("offset len dev:", len(dev))
	for k, v := range dev {
		fmt.Println("k=", k, ", id=", v.Id)
	}
}
func PageNoListDevice(pageNo, pageSize int) []model.Device {
	var dev []model.Device

	model.DB.Debug().Order("id asc").Find(&dev).Limit(pageSize).Offset((pageNo - 1) * pageSize).Find(&dev)
	fmt.Println("offset len dev:", len(dev))
	for k, v := range dev {
		fmt.Println("k=", k, ", id=", v.Id)
	}

	return dev
}

func main() {
	model.DB = model.InitDB()

	//CountDevice()
	//LimitDevice(4)
	//	OffsetDevice(2)
	PageNoListDevice(3, 3)
}
