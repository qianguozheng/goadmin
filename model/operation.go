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

type Domains struct {
	Domain      string
	DomainRefer int
}
type IPs struct {
	Ip      string
	IpRefer int
}

type ProjectIps struct {
	Id       int
	IpsRefer int
}

type ProjectDomains struct {
	Id           int
	DomainsRefer int
}

func AddIps(ip IPs) {
	var i IPs
	if DB.Debug().Where("ip=? and ip_refer=?", ip.Ip, ip.IpRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ip)
}

func AddDomains(domain Domains) {
	var i Domains
	if DB.Debug().Where("domain=? and domain_refer=?", domain.Domain, domain.DomainRefer).RowsAffected > 0 {
		return
	}
	DB.Create(&domain)
}

func GetIpsByIpsRefer() []IPs {

}

func GetDomainsByDomainsRefer() []Domains {

}

func AddProjectIps(ips ProjectIps) {
	var i ProjectIps
	if DB.Debug().Where("id=? and ips_refer=?", ips.Id, ips.IpsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func AddProjectIps(ips ProjectDomains) {
	var i ProjectDomains
	if DB.Debug().Where("id=? and domains_refer=?", ips.Id, ips.DomainsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func GetProjectIpsByIpsRefer(id int) []ProjectIps {
	var i []ProjectIps
	DB.Where("ips_refer=?", id).Find(&i)
	return i
}

func GetProjectIpsByDomainsRefer(id int) []ProjectDomains {
	var i []ProjectDomains
	DB.Where("domains_refer=?", id).Find(&i)
	return i
}

func DeleteProjectIpsByIpsRefer(refer int) {
	DB.Model(&ProjectIps{}).Where("ips_refer=?", refer).Delete(&ProjectIps{})
}
func DeleteProjectIpsByDomainsRefer(refer int) {
	DB.Model(&ProjectDomains{}).Where("domains_refer=?", refer).Delete(&ProjectDomains{})
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
