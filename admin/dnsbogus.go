package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"../rpc"
	"github.com/labstack/echo"
)

func HandleDnsBogusList(c echo.Context) error {

	dns := model.GetAllDnsBogus()
	prjs := model.GetProjects()
	fmt.Println("dns bogus:", dns)
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "dnsbogus_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
		"DnsBogus": dns,
	})
}

func HandleDnsBogusAdd(c echo.Context) error {

	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "dnsbogus_add.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func HandleDnsBogusEdit(c echo.Context) error {

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

func HandleDnsBogusSave(c echo.Context) error {
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

func HandleDnsBogusUpdate(c echo.Context) error {
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
func HandleDnsBogusDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	id, _ := strconv.Atoi(ids)

	model.DeleteDnsBogusById(id)

	return c.Redirect(http.StatusFound, "/v3/project/dns_bogus/v_list")
}

func HandleDnsBogusSend(c echo.Context) error {
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

func HandleDnsBogusAJAX(c echo.Context) error {
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
