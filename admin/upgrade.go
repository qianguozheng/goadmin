package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"github.com/labstack/echo"
)

func HandleProjectUpgradeManage(c echo.Context) error {
	models := model.GetAllModels()
	firmware := model.GetAllFirmwares()
	selected := c.QueryParam("queryModel")

	//	fmt.Println("selected:", selected)
	//	fmt.Println("models:", models)
	//	fmt.Println("firmwares:", firmware)

	if selected == "" {
		fmt.Println("selected nil")
		selected = "3"
	}

	s, _ := strconv.Atoi(selected)

	path := RequestUrl(c)
	fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
		"Selected": s,
		"Models":   models,
		"Firmware": firmware,
		"Path":     path,
	})
}
func HandleProjectUpgradeAdd(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_add.html", map[string]interface{}{
		"Path": path,
	})
}

func HandleProjectUpgradeUpdate(c echo.Context) error {
	id := c.FormValue("id")

	name := c.FormValue("name")
	version := c.FormValue("v")
	hash := c.FormValue("hash")
	url := c.FormValue("url")
	modell := c.FormValue("model")
	state := c.FormValue("state")
	audit := c.FormValue("auditFlag")
	rcl := c.FormValue("rclFlag")
	remark := c.FormValue("remark")

	fmt.Println("id:", id)
	fmt.Println("name:", name)
	fmt.Println("version:", version)
	fmt.Println("hash:", hash)
	fmt.Println("url:", url)
	fmt.Println("model:", modell)
	fmt.Println("state:", state)
	fmt.Println("audit:", audit)
	fmt.Println("rcl:", rcl)
	fmt.Println("remark:", remark)

	var a, r bool
	if audit == "0" {
		a = false
	} else {
		a = true
	}

	if rcl == "0" {
		r = false
	} else {
		r = true
	}

	m, _ := strconv.Atoi(modell)
	idi, _ := strconv.Atoi(id)
	model.UpdateFirmware(remark, hash, url, state, name, version, r, a, m, idi)
	return c.Redirect(http.StatusFound, "/v3/project/firmware/v_list")
}
func HandleProjectUpgradeEdit(c echo.Context) error {
	id := c.QueryParam("id")
	mod := c.QueryParam("queryModel")
	idi, _ := strconv.Atoi(id)
	modi, _ := strconv.Atoi(mod)

	u := model.GetFirmwareById(idi, modi)
	fmt.Println(u)
	path := RequestUrl(c)
	fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_edit.html", map[string]interface{}{
		"Id":      id,
		"Upgrade": u,
		"Path":    path,
	})
}

func HandleProjectUpgradeDelete(c echo.Context) error {
	//return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{})
	id := c.QueryParam("ids")
	mod := c.QueryParam("queryModel")
	modeli, _ := strconv.Atoi(mod)

	fmt.Println("id, modeli", id, mod)
	if id == "" {
		model.DeleteFirmwareByModel(modeli)
	}

	idi, _ := strconv.Atoi(id)

	model.DeleteFirmwareById(idi, modeli)
	return c.Redirect(http.StatusFound, "/v3/project/firmware/v_list")
}

func HandleProjectUpgradeSave(c echo.Context) error {
	name := c.FormValue("name")
	version := c.FormValue("v")
	hash := c.FormValue("hash")
	url := c.FormValue("url")
	modell := c.FormValue("model")
	state := c.FormValue("state")
	audit := c.FormValue("auditFlag")
	rcl := c.FormValue("rclFlag")
	remark := c.FormValue("remark")

	//	fmt.Println("name:", name)
	//	fmt.Println("version:", version)
	//	fmt.Println("hash:", hash)
	//	fmt.Println("url:", url)
	//	fmt.Println("model:", modell)
	//	fmt.Println("state:", state)
	//	fmt.Println("audit:", audit)
	//	fmt.Println("rcl:", rcl)
	//	fmt.Println("remark:", remark)

	//func (remark, md5, url, status, name, ver string, rcl, audit bool, model int)
	var a, r bool
	if audit == "0" {
		a = false
	} else {
		a = true
	}

	if rcl == "0" {
		r = false
	} else {
		r = true
	}

	m, _ := strconv.Atoi(modell)
	model.InsertFirmware(remark, hash, url, state, name, version, r, a, m)

	return c.Redirect(http.StatusFound, "/v3/project/firmware/v_list")
}
