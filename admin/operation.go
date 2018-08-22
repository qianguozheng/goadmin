package admin

import (
	//	"fmt"
	"net/http"
	//	"strconv"

	"fmt"
	"strconv"
	"strings"

	"../model"
	"github.com/labstack/echo"
)

// //Trust Ips
//
// adminGrp.GET("/v3/project/trust_ip/v_list", admin.HandleTrustIpsList)
// adminGrp.GET("/v3/project/trust_ip/v_add", admin.HandleTrustIpsAdd)
// adminGrp.GET("/v3/project/trust_ip/v_edit", admin.HandleTrustIpsEdit)
// adminGrp.POST("/v3/project/trust_ip/o_save", admin.HandleTrustIpsSave)
// adminGrp.GET("/v3/project/trust_ip/o_delete", admin.HandleTrustIpsDelete)
// adminGrp.POST("/v3/project/trust_ip/o_update", admin.HandleTrustIpsUpdate)
//
// //Trust Domain
// adminGrp.GET("/v3/project/trust_domain/v_list", admin.HandleTrustDomainsList)
// adminGrp.GET("/v3/project/trust_domain/v_add", admin.HandleTrustDomainsAdd)
// adminGrp.GET("/v3/project/trust_domain/v_edit", admin.HandleTrustDomainsEdit)
// adminGrp.POST("/v3/project/trust_domain/o_save_domain_strategy", admin.HandleTrustDomainsSave)
// adminGrp.GET("/v3/project/trust_domain/o_delete", admin.HandleTrustDomainsDelete)
// adminGrp.POST("/v3/project/trust_domain/o_update", admin.HandleTrustDomainsUpdate)

type TrustIpsController struct{}

func (self TrustIpsController) RegisterRoute(g *echo.Group) {
	g.Match([]string{"GET", "POST"}, "/trust_ip/v_list", self.List)
	g.GET("/trust_ip/v_add", self.Add)
	g.GET("/trust_ip/v_edit", self.Edit)
	g.POST("/trust_ip/o_save", self.Save)
	g.GET("/trust_ip/o_delete", self.Delete)
	g.POST("/trust_ip/o_update", self.Update)
}

type TrustDomainsController struct{}

func (self TrustDomainsController) RegisterRoute(g *echo.Group) {
	g.Match([]string{"GET", "POST"}, "/trust_domain/v_list", self.List)
	g.GET("/trust_domain/v_add", self.Add)
	g.GET("/trust_domain/v_edit", self.Edit)
	g.POST("/trust_domain/o_save_domain_strategy", self.Save)
	g.GET("/trust_domain/o_delete", self.Delete)
	g.POST("/trust_domain/o_update", self.Update)
}

func (self TrustIpsController) List(c echo.Context) error {
	path := RequestUrl(c)

	//trust := model.GetAllTrustIps()
	ipsNum := model.GetTotalTrustIpsNum()
	page := GeneratePage(c, ipsNum)
	trust := model.ListPageNoTrustIps(page.PageNo, page.PageSize)

	//ips := model.GetAllIps()
	var ips []model.IPs
	for _, v := range trust {
		ip := model.GetIpsByIpsRefer(v.Id)
		ips = append(ips, ip...)
	}

	prjs := model.GetProjects()

	// fmt.Println("trust:", trust)
	// fmt.Println("ips:", ips)
	// fmt.Println("Projects:", prjs)

	return c.Render(http.StatusOK, "trustip_list.html", echo.Map{
		"Path":     path,
		"IPs":      ips,
		"TrustIps": trust,
		"Projects": prjs,
		"Page":     page,
	})
}
func (self TrustIpsController) Add(c echo.Context) error {
	path := RequestUrl(c)
	prjs := model.GetProjects()
	//fmt.Println("Projects:", prjs)
	return c.Render(http.StatusOK, "trustip_add.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func (self TrustIpsController) Save(c echo.Context) error {

	//printFormParams(c)

	content := c.FormValue("content")
	prjs := c.FormValue("projectIds")
	prj := strings.Split(prjs, ",")

	//tp := c.FormValue("type")

	name := c.FormValue("name")
	remark := c.FormValue("remark")
	status := c.FormValue("status") // 0 - global, 1 - specified

	statId, _ := strconv.Atoi(status)
	trustIp := model.TrustIps{
		Name:    name,
		Comment: remark,
		Global:  statId,
		IpNum:   len(prj),
	}
	if false == model.AddTrustIps(trustIp) {
		return c.String(http.StatusOK, "Already have such projects")
	}
	trustIp = model.GetTrustIpsByName(name)
	fmt.Println("trustIps-id=", trustIp.Id)

	for _, v := range prj {
		fmt.Println("project:", v)
		id, _ := strconv.Atoi(v)
		pip := model.ProjectIps{
			IpsRefer:  trustIp.Id,
			ProjectId: id,
		}
		fmt.Println("projectIps:", pip.IpsRefer)
		model.AddProjectIps(pip)
	}

	ips := strings.Split(content, ",")
	for _, v := range ips {
		ip := model.IPs{
			Ip:      v,
			IpRefer: trustIp.Id,
		}
		fmt.Println("ips:", ip)
		model.AddIps(ip)
	}

	//fmt.Println("content:", content)
	//fmt.Println("remark:", remark)
	//fmt.Println("status:", status)
	//fmt.Println("name:", name)
	//fmt.Println("projects:", prjs)

	return c.Redirect(http.StatusFound, "/v3/project/trust_ip/v_list")
}

func (self TrustIpsController) Delete(c echo.Context) error {
	ids := c.QueryParam("ids")
	id, _ := strconv.Atoi(ids)
	fmt.Println("ids=", id)

	model.DeleteTrustIpsById(id)
	model.DeleteProjectIpsByIpsRefer(id)
	model.DeleteIps(id)
	return c.Redirect(http.StatusFound, "/v3/project/trust_ip/v_list")
}

func (self TrustIpsController) Edit(c echo.Context) error {
	//return c.Redirect(http.StatusFound, "/v3/project/trust_ip/v_list")
	ids := c.QueryParam("id")
	id, _ := strconv.Atoi(ids)
	prjIp := model.GetProjectIpsByIpsRefer(id)
	ips := model.GetIpsByIpsRefer(id)
	trustIps := model.GetTrustIpsById(id)

	prjs := model.GetProjects()
	path := RequestUrl(c)

	for _, v := range prjs {
		fmt.Println("id:", v.ID, "Name:", v.Name)
	}
	for _, v := range prjIp {
		fmt.Println("prjectId:", v.ProjectId)
	}

	return c.Render(http.StatusOK, "trustip_edit.html", echo.Map{
		"Path":      path,
		"Projects":  prjs,
		"ProjectIp": prjIp,
		"IPS":       ips,
		"TrustIps":  trustIps,
	})
}

func (self TrustIpsController) Update(c echo.Context) error {

	content := c.FormValue("content")
	prjs := c.FormValue("projectIds")
	prj := strings.Split(prjs, ",")
	ids := c.FormValue("id")
	id, _ := strconv.Atoi(ids)
	name := c.FormValue("name")
	remark := c.FormValue("remark")
	status := c.FormValue("status") // 0 - global, 1 - specified

	statId, _ := strconv.Atoi(status)
	trustIp := model.TrustIps{
		Id:      id,
		Name:    name,
		Comment: remark,
		Global:  statId,
		IpNum:   len(prj),
	}
	//Update trustIps
	//model.DeleteTrustIpsById(id)
	model.UpdateTrustIpsById(trustIp)

	//Update Ips
	model.DeleteProjectIpsByIpsRefer(trustIp.Id)
	for _, v := range prj {
		fmt.Println(v)
		id, _ := strconv.Atoi(v)
		pip := model.ProjectIps{
			IpsRefer:  trustIp.Id,
			ProjectId: id,
		}
		fmt.Println("projectIps:", pip)
		model.AddProjectIps(pip)
	}

	//Update ProjectIps
	model.DeleteIps(trustIp.Id)
	ips := strings.Split(content, ",")
	for _, v := range ips {
		ip := model.IPs{
			Ip:      v,
			IpRefer: trustIp.Id,
		}
		fmt.Println("ips:", ip)
		model.AddIps(ip)
	}

	return c.Redirect(http.StatusFound, "/v3/project/trust_ip/v_list")
}

func (self TrustDomainsController) Add(c echo.Context) error {
	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "trustdomain_add.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func (self TrustDomainsController) List(c echo.Context) error {

	//trust := model.GetAllTrustDomains()
	domainsNum := model.GetTotalTrustDomainsNum()
	page := GeneratePage(c, domainsNum)
	trust := model.ListPageNoTrustDomains(page.PageNo, page.PageSize)

	// domains := model.GetAllDomains()
	var domains []model.Domains
	for _, v := range trust {
		domain := model.GetDomainsByDomainsRefer(v.Id)
		domains = append(domains, domain...)
	}

	// for _, v := range trust {
	// 	fmt.Println("  name:", v.Name)
	// 	fmt.Println("gloabl:", v.Global)
	// }

	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "trustdomain_list.html", echo.Map{
		"Path":         path,
		"Projects":     prjs,
		"TrustDomains": trust,
		"Domains":      domains,
		"Page":         page,
	})
}

func (self TrustDomainsController) Edit(c echo.Context) error {

	ids := c.QueryParam("id")
	id, _ := strconv.Atoi(ids)
	prjDomain := model.GetProjectDomainsByDomainsRefer(id)
	domains := model.GetDomainsByDomainsRefer(id)
	trustDomains := model.GetTrustDomainsById(id)

	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "trustdomain_edit.html", echo.Map{
		"Path":          path,
		"Projects":      prjs,
		"ProjectDomain": prjDomain,
		"TrustDomain":   trustDomains,
		"Domains":       domains,
	})
}

func (self TrustDomainsController) Save(c echo.Context) error {

	printFormParams(c)
	content := c.FormValue("content")
	prjs := c.FormValue("projectIds")
	prj := strings.Split(prjs, ",")

	//tp := c.FormValue("type")

	name := c.FormValue("name")
	remark := c.FormValue("remark")
	status := c.FormValue("status") // 0 - global, 1 - specified

	statId, _ := strconv.Atoi(status)
	//	trustIp := model.TrustIps{
	//		Name:    name,
	//		Comment: remark,
	//		Global:  statId,
	//		IpNum:   len(prj),
	//	}
	trustDomain := model.TrustDomains{
		Name:      name,
		Comment:   remark,
		Global:    statId,
		DomainNum: len(prj),
	}

	fmt.Println("domain num:", len(prj))

	if false == model.AddTrustDomains(trustDomain) {
		return c.String(http.StatusOK, "Already have such projects")
	}
	trustDomain = model.GetTrustDomainsByName(name)
	fmt.Println("trustDomain-id=", trustDomain.Id)

	for _, v := range prj {
		fmt.Println("prj:", v)
		id, _ := strconv.Atoi(v)
		pip := model.ProjectDomains{
			DomainsRefer: trustDomain.Id,
			ProjectId:    id,
		}
		fmt.Println("projectdomains:", pip)
		model.AddProjectDomains(pip)
	}

	ips := strings.Split(content, ",")
	for _, v := range ips {
		d := model.Domains{
			Domain:      v,
			DomainRefer: trustDomain.Id,
		}
		fmt.Println("domain:", d)
		model.AddDomains(d)
	}

	//fmt.Println("content:", content)
	//fmt.Println("remark:", remark)
	//fmt.Println("status:", status)
	//fmt.Println("name:", name)
	//fmt.Println("projects:", prjs)

	return c.Redirect(http.StatusFound, "/v3/project/trust_domain/v_list")
}
func (self TrustDomainsController) Delete(c echo.Context) error {
	ids := c.QueryParam("ids")
	id, _ := strconv.Atoi(ids)
	fmt.Println("ids=", id)

	model.DeleteTrustDomainsById(id)
	model.DeleteProjectDomainsByDomainsRefer(id)
	model.DeleteDomains(id)
	return c.Redirect(http.StatusFound, "/v3/project/trust_domain/v_list")
}
func (self TrustDomainsController) Update(c echo.Context) error {
	content := c.FormValue("content")
	prjs := c.FormValue("projectIds")
	prj := strings.Split(prjs, ",")
	ids := c.FormValue("id")
	id, _ := strconv.Atoi(ids)
	name := c.FormValue("name")
	remark := c.FormValue("remark")
	status := c.FormValue("status") // 0 - global, 1 - specified

	statId, _ := strconv.Atoi(status)
	trustDomain := model.TrustDomains{
		Id:        id,
		Name:      name,
		Comment:   remark,
		Global:    statId,
		DomainNum: len(prj),
	}
	fmt.Println("global:", status)

	//Update trustIps
	//model.DeleteTrustIpsById(id)
	model.UpdateTrustDomainsById(trustDomain)

	//Update Ips
	model.DeleteProjectDomainsByDomainsRefer(trustDomain.Id)
	for _, v := range prj {
		fmt.Println(v)
		id, _ := strconv.Atoi(v)
		pip := model.ProjectDomains{
			DomainsRefer: trustDomain.Id,
			ProjectId:    id,
		}
		fmt.Println("projectdomains:", pip)
		model.AddProjectDomains(pip)
	}

	//Update ProjectDomains
	model.DeleteDomains(trustDomain.Id)
	dms := strings.Split(content, ",")
	for _, v := range dms {
		domain := model.Domains{
			Domain:      v,
			DomainRefer: trustDomain.Id,
		}
		fmt.Println("domain:", domain)
		model.AddDomains(domain)
	}

	return c.Redirect(http.StatusFound, "/v3/project/trust_domain/v_list")
}
