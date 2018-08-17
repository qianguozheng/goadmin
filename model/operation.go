package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

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
	DomainRefer int //TrustDomains Id
}
type IPs struct {
	Ip      string
	IpRefer int //TrustIps Id
}

type ProjectIps struct {
	Id        int
	ProjectId int
	IpsRefer  int //TrustIps Id
}

type ProjectDomains struct {
	Id           int
	ProjectId    int
	DomainsRefer int //TrustDomains Id
}

func AddIps(ip IPs) {
	var i IPs
	if DB.Where("ip=? and ip_refer=?", ip.Ip, ip.IpRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ip)
}

func DeleteIps(refer int) {
	DB.Where("ip_refer=?", refer).Delete(&IPs{})
}
func DeleteDomains(refer int) {
	DB.Where("domain_refer=?", refer).Delete(&Domains{})
}

func AddDomains(domain Domains) bool {
	var i Domains
	if DB.Where("domain=? and domain_refer=?", domain.Domain, domain.DomainRefer).Find(&i).RowsAffected > 0 {
		return false
	}
	DB.Create(&domain)
	return true
}

func GetIpsByIpsRefer(refer int) []IPs {
	var i []IPs
	DB.Where("ip_refer=?", refer).Find(&i)
	return i
}

func GetAllIps() []IPs {
	var i []IPs
	DB.Model(&IPs{}).Find(&i)
	return i
}

func GetAllDomains() []Domains {
	var i []Domains
	DB.Model(&Domains{}).Find(&i)
	return i
}

func GetDomainsByDomainsRefer(refer int) []Domains {
	var i []Domains
	DB.Where("domain_refer=?", refer).Find(&i)
	return i
}

func AddProjectIps(ips ProjectIps) {
	var i ProjectIps
	if DB.Where("project_id=? and ips_refer=?", ips.ProjectId, ips.IpsRefer).Find(&i).RowsAffected > 0 {
		return
	}
	DB.Create(&ips)
}

func AddProjectDomains(ips ProjectDomains) {
	var i ProjectDomains
	if DB.Where("project_id=? and domains_refer=?", ips.ProjectId, ips.DomainsRefer).Find(&i).RowsAffected > 0 {
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
	if DB.Where("name=?", ips.Name).Find(&u).RowsAffected > 1 {
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

//Rcl need
func GetIpsByGlobal() []IPs {
	var i []TrustIps
	DB.Where("global=?", 0).Find(&i)

	var ips []IPs
	var ip []IPs
	for _, v := range i {
		DB.Where("ip_refer=?", v.Id).Find(&ip)
		ips = append(ips, ip...)
	}
	return ips
}

func GetDomainsByGloabl() []Domains {
	var i []TrustDomains
	DB.Where("global=?", 0).Find(&i)

	fmt.Println("i", i)
	var domains []Domains
	var domain []Domains
	for _, v := range i {
		DB.Where("domain_refer=?", v.Id).Find(&domain)
		domains = append(domains, domain...)
	}
	return domains
}

//according to ProjectId find trust ID
func GetIpsByProjectId(id int) []IPs {
	var i []ProjectIps
	DB.Where("project_id=?", id).Find(&i)

	var ips []IPs
	var ip []IPs
	for _, v := range i {
		trustId := v.IpsRefer
		DB.Where("ip_refer=?", trustId).Find(&ip)
		ips = append(ips, ip...)
	}
	return ips
}

func GetDomainsByProjectId(id int) []Domains {
	var i []ProjectDomains
	DB.Where("project_id=?", id).Find(&i)

	var domains []Domains
	var domain []Domains
	for _, v := range i {
		trustId := v.DomainsRefer
		DB.Where("domain_refer=?", trustId).Find(&domain)
		domains = append(domains, domain...)
	}
	return domains
}

// end rcl

func GetTrustIpsByName(name string) TrustIps {
	var i TrustIps
	DB.Model(&TrustIps{}).Where("name=?", name).Find(&i)
	return i
}

func GetTrustDomainsByName(name string) TrustDomains {
	var i TrustDomains
	DB.Where("name=?", name).Find(&i)
	return i
}

func GetAllTrustIps() []TrustIps {
	var i []TrustIps
	DB.Model(&TrustIps{}).Find(&i)
	return i
}

func UpdateTrustIpsById(ips TrustIps) {
	DB.Model(&TrustIps{}).Where("id=?", ips.Id).Update(&ips)
	DB.Model(&TrustIps{}).Where("id=?", ips.Id).Update("global", ips.Global)
}

func UpdateTrustIpsByName(ips TrustIps) {
	DB.Model(&TrustIps{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustIpsById(id int) {
	DB.Where("id=?", id).Delete(&TrustIps{})
}

func AddTrustDomains(domains TrustDomains) bool {
	var d TrustDomains
	if DB.Where("name=?", domains.Name).Find(&d).RowsAffected > 0 {
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
	DB.Model(&TrustDomains{}).Where("id=?", ips.Id).Update(&ips)
	DB.Model(&TrustDomains{}).Where("id=?", ips.Id).Update("global", ips.Global)
}

func UpdateTrustDomainsByName(ips TrustDomains) {
	DB.Model(&TrustDomains{}).Where("name=?", ips.Name).Update(&ips)
}

func DeleteTrustDomainsById(id int) {
	DB.Where("id=?", id).Delete(&TrustDomains{})
}

//Update after change trustips or trustdomains
func (ips *TrustIps) AfterUpdate(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectIpsByIpsRefer(ips.Id)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (domain *TrustDomains) AfterUpdate(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectDomainsByDomainsRefer(domain.Id)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (domain *Domains) AfterUpdate(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectDomainsByDomainsRefer(domain.DomainRefer)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (ips *IPs) AfterUpdate(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectIpsByIpsRefer(ips.IpRefer)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (prj *ProjectDomains) AfterUpdate(tx *gorm.DB) (err error) {

	//Find device Id by Project Id
	devs := GetDeviceByProjectId(prj.ProjectId)
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

func (prj *ProjectIps) AfterUpdate(tx *gorm.DB) (err error) {

	//Find device Id by Project Id
	devs := GetDeviceByProjectId(prj.ProjectId)
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

func (domain *Domains) AfterCreate(tx *gorm.DB) (err error) {

	fmt.Println("after create")
	//Get Project by Id
	prjs := GetProjectDomainsByDomainsRefer(domain.DomainRefer)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (ips *IPs) AfterCreate(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectIpsByIpsRefer(ips.IpRefer)
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (domain *Domains) AfterDelete(tx *gorm.DB) (err error) {

	fmt.Println("after create")
	//Get Project by Id
	prjs := GetProjectDomainsByDomainsRefer(domain.DomainRefer)
	//[]ProjectIps
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}

func (ips *IPs) AfterDelete(tx *gorm.DB) (err error) {

	//Get Project by Id
	prjs := GetProjectIpsByIpsRefer(ips.IpRefer)
	for _, v := range prjs {
		//Find device Id by Project Id
		devs := GetDeviceByProjectId(v.ProjectId)
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
	}

	return nil
}
