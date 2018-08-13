package model

type TrustIps struct {
	Id    int `gorm:"AUTO_INCREMENT"`
	Name  string
	IpNum int //192.168.0.0-192.168.3.4,172.15.0.0-172.17.0.255
	//	ProjectRefer int    //Foreign Key
	Comment string
	Global  int
}

type TrustDomains struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	Name      string
	DomainNum int
	Comment   string
	Global    int
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
	Id        int
	ProjectId int
	IpsRefer  int
}

type ProjectDomains struct {
	Id           int
	ProjectId    int
	DomainsRefer int
}

func AddIps(ip IPs) {
	var i IPs
	if DB.Debug().Where("ip=? and ip_refer=?", ip.Ip, ip.IpRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Debug().Create(&ip)
}

func DeleteIps(refer int) {
	DB.Debug().Where("ip_refer=?", refer).Delete(&IPs{})
}
func DeleteDomains(refer int) {
	DB.Debug().Where("domain_refer=?", refer).Delete(&Domains{})
}

func AddDomains(domain Domains) bool {
	var i Domains
	if DB.Debug().Where("domain=? and domain_refer=?", domain.Domain, domain.DomainRefer).Find(&i).RowsAffected > 0 {
		return false
	}
	DB.Create(&domain)
	return true
}

func GetIpsByIpsRefer(refer int) []IPs {
	var i []IPs
	DB.Debug().Where("ip_refer=?", refer).Find(&i)
	return i
}

func GetAllIps() []IPs {
	var i []IPs
	DB.Debug().Model(&IPs{}).Find(&i)
	return i
}

func GetAllDomains() []Domains {
	var i []Domains
	DB.Debug().Model(&Domains{}).Find(&i)
	return i
}

func GetDomainsByDomainsRefer(refer int) []Domains {
	var i []Domains
	DB.Debug().Where("domain_refer=?", refer).Find(&i)
	return i
}

func AddProjectIps(ips ProjectIps) {
	var i ProjectIps
	if DB.Debug().Where("project_id=? and ips_refer=?", ips.ProjectId, ips.IpsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Debug().Create(&ips)
}

func AddProjectDomains(ips ProjectDomains) {
	var i ProjectDomains
	if DB.Debug().Where("project_id=? and domains_refer=?", ips.ProjectId, ips.DomainsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func GetProjectIpsByIpsRefer(id int) []ProjectIps {
	var i []ProjectIps
	DB.Where("ips_refer=?", id).Find(&i)
	return i
}

func GetProjectDomainsByDomainsRefer(id int) []ProjectDomains {
	var i []ProjectDomains
	DB.Where("domains_refer=?", id).Find(&i)
	return i
}

func DeleteProjectIpsByIpsRefer(refer int) {
	DB.Model(&ProjectIps{}).Where("ips_refer=?", refer).Delete(&ProjectIps{})
}
func DeleteProjectDomainsByDomainsRefer(refer int) {
	DB.Model(&ProjectDomains{}).Where("domains_refer=?", refer).Delete(&ProjectDomains{})
}

func AddTrustIps(ips TrustIps) bool {
	var u TrustIps
	if DB.Debug().Where("name=", ips.Name).Find(&u).RowsAffected > 1 {
		return false
	}
	DB.Create(&ips)
	return true
}

func GetTrustIpsById(id int) TrustIps {
	var ips TrustIps
	DB.Where("id=?", id).Find(&ips)
	return ips
}
func GetTrustIpsByName(name string) TrustIps {
	var i TrustIps
	DB.Model(&TrustIps{}).Where("name=?", name).Find(&i)
	return i
}

func GetTrustDomainsByName(name string) TrustDomains {
	var i TrustDomains
	DB.Model(&TrustDomains{}).Where("name=?", name).Find(&i)
	return i
}

func GetAllTrustIps() []TrustIps {
	var i []TrustIps
	DB.Model(&TrustIps{}).Find(&i)
	return i
}

func UpdateTrustIpsById(ips TrustIps) {
	DB.Model(&TrustIps{}).Debug().Where("id=?", ips.Id).Update(&ips)
	DB.Model(&TrustIps{}).Debug().Where("id=?", ips.Id).Update("global", ips.Global)
}

func UpdateTrustIpsByName(ips TrustIps) {
	DB.Model(&TrustIps{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustIpsById(id int) {
	DB.Where("id=?", id).Delete(&TrustIps{})
}

func AddTrustDomains(domains TrustDomains) bool {
	var d TrustDomains
	if DB.Debug().Where("name=", domains.Name).Find(&d).RowsAffected > 1 {
		return false
	}
	DB.Create(&domains)
	return true
}

func GetAllTrustDomains() []TrustDomains {
	var i []TrustDomains
	DB.Model(&TrustDomains{}).Find(&i)
	return i
}
func GetTrustDomainsById(id int) TrustDomains {
	var ips TrustDomains
	DB.Where("id=?", id).Find(&ips)
	return ips
}

func UpdateTrustDomainsById(ips TrustDomains) {
	DB.Model(&TrustDomains{}).Debug().Where("id=?", ips.Id).Update(&ips)
	DB.Model(&TrustDomains{}).Debug().Where("id=?", ips.Id).Update("global", ips.Global)
}

func UpdateTrustDomainsByName(ips TrustDomains) {
	DB.Model(&TrustDomains{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustDomainsById(id int) {
	DB.Where("id=?", id).Delete(&TrustDomains{})
}
