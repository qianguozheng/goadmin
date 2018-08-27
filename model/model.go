package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qianguozheng/goadmin/config"
)

var DB *gorm.DB

func Mysql() string {
	mysqlConfig, err := config.ConfigFile.GetSection("mysql")
	if err != nil {
		fmt.Println("get mysql config error:", err)
		return err.Error()
	}

	tmpDns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])
	// db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	// defer db.Close()
	return tmpDns
}

func InitDB(path string, dbType string) *gorm.DB {

	db, err := gorm.Open(dbType, path)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect to database")
	}

	db.AutoMigrate(&User{}, &Upgrade{}, &Model{},
		&Device{}, &Qos{}, &Wan{}, &WanQos{}, &Ssid{}, &Project{},
		&Md5{}, &TermFree{}, &TrustIps{}, &TrustDomains{}, &DnsBogus{},
		&IPs{}, &Domains{}, &ProjectIps{}, &ProjectDomains{}, &DeviceStatus{})

	user := new(User)
	db.First(user, "name=?", "admin")
	if user.Password == "" {
		db.Create(&User{Name: "admin", Password: "pass"})
	}

	return db
}
