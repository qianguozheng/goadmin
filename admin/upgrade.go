package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/qianguozheng/goadmin/model"
)

type UpgradeController struct{}

// adminGrp.GET("/v3/project/firmware/v_list", admin.HandleProjectUpgradeManage)
// adminGrp.GET("/v3/project/firmware/v_add", admin.HandleProjectUpgradeAdd)
// adminGrp.POST("/v3/project/firmware/o_save", admin.HandleProjectUpgradeSave)
// adminGrp.GET("/v3/project/firmware/v_edit", admin.HandleProjectUpgradeEdit)
// adminGrp.POST("/v3/project/firmware/v_delete", admin.HandleProjectUpgradeDelete)
// adminGrp.GET("/v3/project/firmware/o_delete", admin.HandleProjectUpgradeDelete)
// adminGrp.POST("/v3/project/firmware/o_delete", admin.HandleProjectUpgradeDelete)
// adminGrp.POST("/v3/project/firmware/o_update", admin.HandleProjectUpgradeUpdate)
// 注册路由
func (self UpgradeController) RegisterRoute(g *echo.Group) {
	g.GET("/firmware/v_list", self.Manage)
	g.GET("/firmware/v_add", self.Add)
	g.GET("/firmware/v_edit", self.Edit)
	g.POST("/firmware/o_save", self.Save)
	g.POST("/firmware/v_delete", self.Delete)
	g.POST("/firmware/o_update", self.Update)
	g.Match([]string{"GET", "POST"}, "/firmware/o_delete", self.Delete)
}

func (self UpgradeController) Manage(c echo.Context) error {
	models := model.GetAllModels()
	firmware := model.GetAllFirmwares()
	selected := c.QueryParam("queryModel")

	if selected == "" {
		fmt.Println("selected nil")
		selected = "3"
	}

	s, _ := strconv.Atoi(selected)

	path := RequestUrl(c)
	user := c.Get("user")
	// fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
		"Selected": s,
		"Models":   models,
		"Firmware": firmware,
		"Path":     path,
		"User":     user,
	})
}
func (self UpgradeController) Add(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_add.html", map[string]interface{}{
		"Path": path,
	})
}

func (self UpgradeController) Update(c echo.Context) error {
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
func (self UpgradeController) Edit(c echo.Context) error {
	id := c.QueryParam("id")
	mod := c.QueryParam("queryModel")
	idi, _ := strconv.Atoi(id)
	modi, _ := strconv.Atoi(mod)

	u := model.GetFirmwareById(idi, modi)
	fmt.Println(u)
	path := RequestUrl(c)
	user := c.Get("user")
	// fmt.Println("path=", path)
	return c.Render(http.StatusOK, "firmware_edit.html", map[string]interface{}{
		"Id":      id,
		"Upgrade": u,
		"Path":    path,
		"User":    user,
	})
}

func (self UpgradeController) Delete(c echo.Context) error {
	//return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{})

	if c.Request().Method == "GET" {
		id := c.QueryParam("ids")
		mod := c.QueryParam("queryModel")
		modeli, _ := strconv.Atoi(mod)

		idi, _ := strconv.Atoi(id)
		model.DeleteFirmwareById(idi, modeli)

	} else if c.Request().Method == "POST" {
		ids := getIDSFromParams(c)
		mod := c.QueryParam("queryModel")
		modeli, _ := strconv.Atoi(mod)

		for _, v := range ids {
			model.DeleteFirmwareById(v, modeli)
		}

	}

	return c.Redirect(http.StatusFound, "/v3/project/firmware/v_list")
}

func (self UpgradeController) Save(c echo.Context) error {
	name := c.FormValue("name")
	version := c.FormValue("v")
	hash := c.FormValue("hash")
	url := c.FormValue("url")
	modell := c.FormValue("model")
	state := c.FormValue("state")
	audit := c.FormValue("auditFlag")
	rcl := c.FormValue("rclFlag")
	remark := c.FormValue("remark")

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
