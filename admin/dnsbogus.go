package admin

import (
	"net/http"
	"strconv"

	"../model"
	"../rpc"
	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

//
// //Dns Bogus
// adminGrp.GET("/v3/project/dns_bogus/v_list", admin.HandleDnsBogusList)
// adminGrp.GET("/v3/project/dns_bogus/v_add", admin.HandleDnsBogusAdd)
// adminGrp.GET("/v3/project/dns_bogus/v_edit", admin.HandleDnsBogusEdit)
// adminGrp.POST("/v3/project/dns_bogus/o_save", admin.HandleDnsBogusSave)
// adminGrp.GET("/v3/project/dns_bogus/o_delete", admin.HandleDnsBogusDelete)
// adminGrp.POST("/v3/project/dns_bogus/o_update", admin.HandleDnsBogusUpdate)
// adminGrp.POST("/v3/project/dns_bogus/o_send", admin.HandleDnsBogusSend)
// adminGrp.GET("/v3/project/dns_bogus/v_ajax_edit", admin.HandleDnsBogusAJAX)
type DnsBogusController struct{}

func (self DnsBogusController) RegisterRoute(g *echo.Group) {
	g.Match([]string{"GET", "POST"}, "/dns_bogus/v_list", self.List)
	g.GET("/dns_bogus/v_add", self.Add)
	g.GET("/dns_bogus/v_edit", self.Edit)
	g.POST("/dns_bogus/o_save", self.Save)
	g.Match([]string{"GET", "POST"}, "/dns_bogus/o_delete", self.Delete)
	g.POST("/dns_bogus/o_update", self.Update)
	g.POST("/dns_bogus/o_send", self.Send)
	g.GET("/dns_bogus/v_ajax_edit", self.AJAX)
}

func (self DnsBogusController) List(c echo.Context) error {

	totalDnsBogus := model.GetTotalDnsBogus()
	page := GeneratePage(c, totalDnsBogus)
	//dns := model.GetAllDnsBogus()
	dns := model.ListPageNoDnsBogus(page.PageNo, page.PageSize)
	prjs := model.GetProjects()
	// fmt.Println("dns bogus:", dns)
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "dnsbogus_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
		"DnsBogus": dns,
		"Page":     page,
	})
}

func (self DnsBogusController) Add(c echo.Context) error {

	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "dnsbogus_add.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func (self DnsBogusController) Edit(c echo.Context) error {

	ids := c.QueryParam("id")
	id, _ := strconv.Atoi(ids)
	dnsBogus := model.GetDnsBogusById(id)

	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "dnsbogus_edit.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
		"DnsBogus": dnsBogus,
	})
}

func (self DnsBogusController) Save(c echo.Context) error {
	prj := c.FormValue("projectId")
	dns := c.FormValue("dns")
	ip := c.FormValue("distDns")
	comment := c.FormValue("remark")
	typeStr := c.FormValue("type")

	typeInt, _ := strconv.Atoi(typeStr)
	prjInt, _ := strconv.Atoi(prj)

	dnsBogus := model.DnsBogus{
		Type:         typeInt,
		ProjectRefer: prjInt,
		Domain:       dns,
		Ip:           ip,
		Comment:      comment,
		Status:       0,
	}

	if false == model.AddDnsBogus(dnsBogus) {
		return c.String(http.StatusCreated, "Already exist such dns bogus")
	}

	return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
}

func (self DnsBogusController) Update(c echo.Context) error {
	//	prj := c.FormValue("projectId")
	dns := c.FormValue("dns")
	ip := c.FormValue("distDns")
	comment := c.FormValue("remark")
	//	typeStr := c.FormValue("type")

	//	typeInt, _ := strconv.Atoi(typeStr)
	//	prjInt, _ := strconv.Atoi(prj)
	statusStr := c.FormValue("status")
	statusInt, _ := strconv.Atoi(statusStr)
	idStr := c.FormValue("id")
	id, _ := strconv.Atoi(idStr)

	dnsBogus := model.GetDnsBogusById(id)
	dnsBogus.Status = statusInt
	dnsBogus.Comment = comment
	dnsBogus.Domain = dns
	dnsBogus.Ip = ip

	model.UpdateDnsBogusById(dnsBogus)

	return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
}
func (self DnsBogusController) Delete(c echo.Context) error {
	if c.Request().Method == "GET" {
		id := goutils.MustInt(c.QueryParam("ids"))
		model.DeleteDnsBogusById(id)

	} else if c.Request().Method == "POST" {
		ids := getIDSFromParams(c)
		for _, v := range ids {
			model.DeleteDnsBogusById(v)
		}
	}
	return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
}

func (self DnsBogusController) Send(c echo.Context) error {
	idStr := c.FormValue("id")
	mac := c.FormValue("mac")

	id, _ := strconv.Atoi(idStr)

	data := rpc.DnsBogusParam{
		Id:  id,
		Mac: mac,
	}

	//Check Status
	dnsBogus := model.GetDnsBogusById(id)
	if dnsBogus.Status <= 0 {
		return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
	}

	if mac != "" {
		//Send to mac
		rpc.DnsBogusRequest(rpc.DeviceKind, data)
		return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
	}
	//Send
	rpc.DnsBogusRequest(rpc.ProjectKind, data)
	return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
}

type AjaxProject struct {
	Id          int    `json:"id"`
	ProjectName string `json:"projectName"`
	ProjectId   int    `json:"projectId"`
}

func (self DnsBogusController) AJAX(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, _ := strconv.Atoi(idStr)
	dns := model.GetDnsBogusById(id)
	prj := model.GetProjectById(dns.ProjectRefer)

	json := AjaxProject{
		Id:          id,
		ProjectName: prj.Name,
		ProjectId:   int(prj.ID),
	}
	return c.JSON(http.StatusOK, &json)
}
