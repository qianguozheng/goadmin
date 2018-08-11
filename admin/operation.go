package admin

import (
	//	"fmt"
	"net/http"
	//	"strconv"

	"../model"
	"github.com/labstack/echo"
	"fmt"
	"strings"
	"strconv"
)

func HandleTrustIpsList(c echo.Context) error {
	path := RequestUrl(c)
	trust := model.GetAllTrustIps()
	ips := model.GetAllIps()
	prjs := model.GetProjects()

	fmt.Println("trust:", trust)
	fmt.Println("ips:", ips)
	fmt.Println("Projects:", prjs)

	return c.Render(http.StatusOK, "trustip_list.html", echo.Map{
		"Path": path,
		"IPs":ips,
		"TrustIps":trust,
		"Projects": prjs,
	})
}
func HandleTrustIpsAdd(c echo.Context) error {
	path := RequestUrl(c)
	prjs := model.GetProjects()
	//fmt.Println("Projects:", prjs)
	return c.Render(http.StatusOK, "trustip_add.html", echo.Map{
		"Path": path,
		"Projects": prjs,
	})
}

func HandleTrustIpsSave(c echo.Context) error {

	//printFormParams(c)

	content := c.FormValue("content")
	prjs := c.FormValue("projectIds")
	prj := strings.Split(prjs, ",")

	//tp := c.FormValue("type")

	name:=c.FormValue("name")
	remark:= c.FormValue("remark")
	status := c.FormValue("status") // 0 - global, 1 - specified

	statId,_:=strconv.Atoi(status)
	trustIp := model.TrustIps{
		Name    :name,
		Comment :remark,
		Global  :statId,
		IpNum: len(prj),
	}
	if false == model.AddTrustIps(trustIp) {
		return c.String(http.StatusOK,"Already have such projects")
	}
	trustIp = model.GetTrustIpsByName(name)
	fmt.Println("trustIps-id=",trustIp.Id)

	for _,v:= range prj{
		fmt.Println(v)
		id,_:= strconv.Atoi(v)
		pip := model.ProjectIps{
			IpsRefer: trustIp.Id,
			Id:id,
		}
		fmt.Println("projectIps:", pip)
		model.AddProjectIps(pip)
	}

	ips := strings.Split(content, ",")
	for _,v := range ips {
		ip := model.IPs{
			Ip:v,
			IpRefer:trustIp.Id,
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

func HandleTrustIpsDelete(c echo.Context) error {
	ids := c.QueryParam("ids")
	id, _ :=strconv.Atoi(ids)
	fmt.Println("ids=", id)

	model.DeleteTrustIpsById(id)
	model.DeleteProjectIpsByIpsRefer(id)
	model.DeleteIps(id)
	return c.Redirect(http.StatusFound, "/v3/project/trust_ip/v_list")
}