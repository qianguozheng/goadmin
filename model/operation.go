package model

type TrustIps struct {
	Id      int `gorm:"AUTO_INCREMENT"`
	Name    string
	IpRange string //192.168.0.0-192.168.3.4,172.15.0.0-172.17.0.255
	//	ProjectRefer int    //Foreign Key
	Comment string
	Global  int
}

type TrustDomains struct {
	Id     int `gorm:"AUTO_INCREMENT"`
	Name   string
	Domain string
	//	ProjectRefer int
	Comment string
	Global  int
}

type ProjectIps struct {
	Id            int
	TrustIpsRefer int
}

type ProjectDomains struct {
	Id                int
	TrustDomainsRefer int
}

func AddProjectIps(ips ProjectIps) {
	var i ProjectIps
	if DB.Debug().Where("id=? and trust_ips_refer=?", ips.Id, ips.TrustIpsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func AddProjectIps(ips ProjectDomains) {
	var i ProjectDomains
	if DB.Debug().Where("id=? and trust_domains_refer=?", ips.Id, ips.TrustDomainsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func GetProjectIpsByTrustIpsRefer(id int) []ProjectIps {
	var i []ProjectIps
	DB.Where("trust_ips_refer=?", id).Find(&i)
	return i
}

func GetProjectIpsByTrustDomainsRefer(id int) []ProjectDomains {
	var i []ProjectDomains
	DB.Where("trust_domains_refer=?", id).Find(&i)
	return i
}

func DeleteProjectIpsByTrustIpsRefer(refer int) {
	DB.Model(&ProjectIps{}).Where("trust_ips_refer=?", refer).Delete(&ProjectIps{})
}
func DeleteProjectIpsByTrustDomainsRefer(refer int) {
	DB.Model(&ProjectDomains{}).Where("trust_domains_refer=?", refer).Delete(&ProjectDomains{})
}

func AddTrustIps(ips TrustIps) {
	var u TrustIps
	if DB.Debug().Where("name=", ips.Name).Find(&u).RowsAffected > 1 {
		return
	}
	DB.Create(&ips)
}

func GetTrustIpsByProjectRefer(id int) []TrustIps {
	var ips []TrustIps
	DB.Where("project_refer=?", id).Find(&ips)
	return ips
}

func UpdateTrustIpsByName(ips TrustIps) {
	DB.Model(&TrustIps{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustIpsById(id int) {
	DB.Where("id=?", id).Delete(&TrustIps{})
}

func AddTrustDomains(domains TrustDomains) {
	var d TrustDomains
	if DB.Debug().Where("name=", domains.Name).Find(&d).RowsAffected > 1 {
		return
	}
	DB.Create(&domains)
}

func GetTrustDomainsByProjectRefer(id int) []TrustDomains {
	var ips []TrustDomains
	DB.Where("project_refer=?", id).Find(&ips)
	return ips
}

func UpdateTrustDomainsByName(ips TrustDomains) {
	DB.Model(&TrustDomains{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustDomainsById(id int) {
	DB.Where("id=?", id).Delete(&TrustDomains{})
}
