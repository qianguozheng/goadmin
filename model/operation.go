package model

type TrustIps struct {
	Id           int `gorm:"AUTO_INCREMENT"`
	Name         string
	IpRange      string //192.168.0.0-192.168.3.4,172.15.0.0-172.17.0.255
	ProjectRefer int    //Foreign Key
	Comment      string
	Global       int
}

type TrustDomains struct {
	Id           int `gorm:"AUTO_INCREMENT"`
	Name         string
	Domain       string
	ProjectRefer int
	Comment      string
	Global       int
}
