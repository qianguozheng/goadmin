package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"../model"
	"../rpc"
	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

type DeviceController struct{}

func (self DeviceController) RegisterRoute(g *echo.Group) {
	g.GET("/device/v_list", self.List)
	g.POST("/device/v_list", self.ListPost)
	g.GET("/device/v_edit_cfg", self.Edit)
	g.POST("/device/o_update", self.UpdateEdit)
	g.POST("/device/o_update_config", self.UpdateCloud)
	g.POST("/device/o_save_ssid", self.UpdateSSID)
	g.GET("/device/v_ajax_edit_ssid", self.EditSSID)
	g.GET("/device/v_ajax_restart", self.Restart)
	g.GET("/device/v_edit_delete", self.DeleteDevice)
	g.POST("/device/v_ajax_read_wan", self.GetWan)
	g.POST("/device/o_update_config_wan", self.UpdateWan)
	g.GET("/device/o_read_cfg", self.ReadConfig)
	g.GET("/device/v_ajax_update_cfg", self.VPN)

	//TODO: v_list page option select/restart/upgrade and son
	g.POST("/device/v_ajax_get_firmwares", self.List)
	g.POST("/device/v_restart", self.List)
	g.POST("/device/v_edit_fast", self.List)

	//Add device
	g.GET("/adddev/v_list", self.AddDev)
	g.POST("/adddev/o_save", self.Save)
}

func (self DeviceController) List(c echo.Context) error {
	path := RequestUrl(c)
	cookie, err := c.Cookie("_cookie_page_size")
	if err != nil {
		//return c.String(http.StatusNotFound, "No cookie_page_size found")
		fmt.Println("no cookie_page size")
	}
	var pageSize int
	if cookie != nil {
		pageSize, err = strconv.Atoi(cookie.Value)
	}
	if err != nil {
		pageSize = 20
	}

	devNum := model.GetTotalDeviceNum()
	pageNum := devNum / pageSize
	if devNum%pageSize > 0 {
		pageNum = pageNum + 1
	}

	devs := model.ListPageNoDevice(1, pageSize)

	for k, v := range devs {
		fmt.Println("modeltype:", v.ModelType)
		ssids := model.GetSsidByDeviceId(v.Id)
		for _, s := range ssids {
			devs[k].Ssid = append(devs[k].Ssid, s)
		}
	}

	//Update online status
	devStatus := model.ListPageNoDeviceStatus(1, pageSize)
	for k, v := range devStatus {
		diff := time.Now().Unix() - v.Heartbeat
		if diff > 300 {
			devStatus[k].Online = false
		} else {
			devStatus[k].Online = true
		}
	}

	prjs := model.GetProjects()
	models := model.GetAllModels()

	//fmt.Println("devices:", devs)
	return c.Render(http.StatusOK, "device_list.html", echo.Map{
		"Path":        path,
		"Devices":     devs,
		"PageNo":      1,
		"PageNum":     pageNum,
		"TotalDevice": devNum,
		"PageSize":    pageSize,
		"Projects":    prjs,
		"Models":      models,
		"DevStatus":   devStatus,
	})
}

func HandleProjectDeviceOffline(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "device_offline.html", echo.Map{
		"Path": path,
	})
}

func (self DeviceController) Edit(c echo.Context) error {
	name := c.QueryParam("modelName")
	// ids := c.QueryParam("id")
	// id, _ := strconv.Atoi(ids)
	id := goutils.MustInt(c.QueryParam("id"))
	dev := model.GetDeviceById(id)
	path := RequestUrl(c)
	// fmt.Println("path=", path)
	// fmt.Println("name=", name)
	// fmt.Println("dev=", dev)

	qos := model.GetQosByDeviceId(id)
	wanQos := model.GetWanQosByQosId(qos.Id)

	// fmt.Println("qos:", qos)
	// fmt.Println("wanQos:", wanQos)

	ssid := model.GetSsidByDeviceId(id)
	md5 := model.GetMd5ByDeviceId(id)
	prjs := model.GetProjects()
	// fmt.Println("ssid:", ssid)
	models := model.GetAllModels()
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":     path,
		"Name":     name,
		"Device":   dev,
		"Qos":      qos,
		"WanQos":   wanQos,
		"SSID":     ssid,
		"Projects": prjs,
		"Models":   models,
		"Md5":      md5,
		"Code":     5,
	})
}

//post list_edit
func (self DeviceController) UpdateEdit(c echo.Context) error {
	// ids := c.FormValue("id")
	name := c.FormValue("name")
	// state := c.FormValue("status")
	// mode := c.FormValue("mode")
	// prjIds := c.FormValue("project")
	//
	// //save data to database
	// id, _ := strconv.Atoi(ids)
	// modei, _ := strconv.Atoi(mode)
	// status, _ := strconv.Atoi(state)
	// prjId, _ := strconv.Atoi(prjIds)

	id := goutils.MustInt(c.FormValue("id"))
	modei := goutils.MustInt(c.FormValue("mode"))
	status := goutils.MustInt(c.FormValue("status"))
	prjId := goutils.MustInt(c.FormValue("project"))

	dev := model.GetDeviceById(id)
	dev.Name = name
	dev.Status = status
	dev.Mode = modei
	dev.ProjectRefer = prjId

	// fmt.Println("id=", id)
	// fmt.Println("name=", name)
	// fmt.Println("state=", state)
	// fmt.Println("mode=", mode)
	// fmt.Println("dev=", dev)
	// fmt.Println("prjId=", prjId)

	model.UpdateDevice(dev)
	md5 := model.GetMd5ByDeviceId(id)
	models := model.GetAllModels()
	prjs := model.GetProjects()
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":     path,
		"Name":     "list_edit",
		"Device":   dev,
		"Projects": prjs,
		"Models":   models,
		"Md5":      md5,
		"Code":     5,
	})
}

//post list_cloud
func (self DeviceController) UpdateCloud(c echo.Context) error {
	//list_cloud
	var page_name string
	page_name = c.FormValue("modelName")

	token := c.FormValue("ccToken")
	port := c.FormValue("ccPort")
	host := c.FormValue("ccHost")
	// ids := c.FormValue("deviceId")
	// id, _ := strconv.Atoi(ids)
	id := goutils.MustInt(c.FormValue("deviceId"))
	dev := model.GetDeviceById(id)

	if token != "" &&
		port != "" &&
		host != "" {
		page_name = "list_cloud"
		fmt.Println("list cloud")
		porti, _ := strconv.Atoi(port)
		dev.CloudHost = host
		dev.CloudPort = porti
		dev.CloudToken = token
	}

	//list_lan
	lanIp := c.FormValue("lanIp")
	lanMask := c.FormValue("lanMask")

	if lanIp != "" &&
		lanMask != "" {
		page_name = "list_lan"
		// fmt.Println("list lan")
		dev.LanIp = lanIp
		dev.LanMask = lanMask
	}

	//list rf

	rfMode := c.FormValue("rfWifiMode")
	rfFreq := c.FormValue("rfFrequency")
	rfPower := c.FormValue("rfPower")

	if rfMode != "" &&
		rfFreq != "" &&
		rfPower != "" {
		page_name = "list_rf"

		rm, _ := strconv.Atoi(rfMode)
		rf, _ := strconv.Atoi(rfFreq)
		rp, _ := strconv.Atoi(rfPower)
		dev.RfMode = rm
		dev.RfFreq = rf
		dev.RfPower = rp
	}
	//list dhcp
	startip := c.FormValue("dhcpStartIp")
	endip := c.FormValue("dhcpEndIp")
	lease := c.FormValue("dhcpLeaseTime")

	if startip != "" &&
		endip != "" &&
		lease != "" {
		page_name = "list_dhcp"
		leasetime, _ := strconv.Atoi(lease)
		dev.DhcpStartIp = startip
		dev.DhcpEndIp = endip
		dev.DhcpLeaseTime = leasetime
	}

	//Qos info
	/*
		DeviceId -> find Qos, WanQos

	*/
	var qos model.Qos
	var wanQos []model.WanQos
	qosDownRate := c.FormValue("qosDownRate")
	qosUpRate := c.FormValue("qosUpRate")
	qosTcpLimit := c.FormValue("qosTcpLimit")
	qosUdpLimit := c.FormValue("qosUdpLimit")
	if qosDownRate != "" &&
		qosUpRate != "" &&
		qosTcpLimit != "" &&
		qosUdpLimit != "" {
		page_name = "list_qos"

		up, _ := strconv.Atoi(qosUpRate)
		down, _ := strconv.Atoi(qosDownRate)
		tcp, _ := strconv.Atoi(qosTcpLimit)
		udp, _ := strconv.Atoi(qosUdpLimit)
		//Qos table
		qos = model.GetQosByDeviceId(id)
		qos.UpRate = up
		qos.DownRate = down
		qos.TcpLimit = tcp
		qos.UdpLimit = udp

		model.UpdateQos(qos)
		//WanQos table
		printFormParams(c)
		wanQos = getWanQosFormParams(c, qos.Id)
		model.UpdateWanQoss(wanQos)

		// fmt.Println("qos:", qos)
		// fmt.Println("wanQos:", wanQos)

	}

	///////////////////////////  Update device info   ////////////////////////
	model.UpdateDevice(dev)

	// fmt.Println(dev)
	path := RequestUrl(c)
	prjs := model.GetProjects()
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":     path,
		"Name":     page_name,
		"Device":   dev,
		"Qos":      qos,
		"WanQos":   wanQos,
		"Projects": prjs,
		"Code":     5,
	})
}

func (self DeviceController) UpdateSSID(c echo.Context) error {
	// printFormParams(c)

	//page_name := c.QueryParam("modelName")
	// ports := c.FormValue("port")
	// port, _ := strconv.Atoi(ports)
	port := goutils.MustInt(c.FormValue("port"))

	name := c.FormValue("name")
	name5g := c.FormValue("name5g")
	url := c.FormValue("url")
	devIds := c.FormValue("deviceId")
	devId, _ := strconv.Atoi(devIds)

	s := model.Ssid{
		Port:        port,
		Name:        name,
		Name5:       name5g,
		Url:         url,
		DeviceRefer: devId,
	}
	// fmt.Println("set ssid:", s)
	model.UpdateSsid(s)

	//pass param to template
	ssid := model.GetSsidByDeviceId(devId)

	dev := model.GetDeviceById(devId)
	path := RequestUrl(c)

	// fmt.Println("path:", path)
	// //fmt.Println("name:", page_name)
	// fmt.Println("dev:", dev)
	// fmt.Println("ssid:", ssid)

	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   "list_ssid",
		"Device": dev,
		"SSID":   ssid,
		"Code":   5,
	})

}

//get ssid for editing
type SsidJson struct {
	DeviceId string `json:"deviceId"`
	Name     string `json:"name"`
	Name5g   string `json:"name5g"`
	Port1    int    `json:"port1"`
	Port2    int    `json:"port2"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

func (self DeviceController) EditSSID(c echo.Context) error {
	// ids := c.FormValue("deviceId")
	// ports := c.FormValue("port")
	// port, _ := strconv.Atoi(ports)
	// id, _ := strconv.Atoi(ids)
	id := goutils.MustInt(c.FormValue("deviceId"))
	port := goutils.MustInt(c.FormValue("port"))

	//dev := model.GetDeviceById(id)
	ssid := model.GetSsidByDeviceIdPort(id, port)

	s := &SsidJson{
		Port1:    port,
		Port2:    port,
		Name:     ssid.Name,
		Name5g:   ssid.Name5,
		Password: ssid.Password,
		Url:      ssid.Url,
	}
	// fmt.Println("ssid:", ssid)
	// fmt.Println("ssidJson:", s)
	return c.JSON(http.StatusOK, s)
}

//ajax request
type WanJson struct {
	Port             int    `json:"port"`
	TxRate           int    `json:"txRate"`
	RxRate           int    `json:"rxRate"`
	WanMode          int    `json:"wanMode"`
	WanIp            string `json:"wanIp"`
	WanMask          string `json:"wanMask"`
	WanPppoeAccount  string `json:"wanPppoeAccount"`
	WanPppoePassword string `json:"wanPppoePassword"`
	WanPrimaryDns    string `json:"wanPrimaryDns"`
	WanSecondaryDns  string `json:"wanSecondaryDns"`
	Success          string `json:"success"`
}

func (self DeviceController) GetWan(c echo.Context) error {
	// ids := c.FormValue("id")
	// ports := c.FormValue("port")
	// id, _ := strconv.Atoi(ids)
	// port, _ := strconv.Atoi(ports)
	id := goutils.MustInt(c.FormValue("id"))
	port := goutils.MustInt(c.FormValue("port"))

	wan := model.Wan{
		Port:        port,
		DeviceRefer: id,
	}
	wancfg := model.GetWanByDeviceIdPort(wan)
	// fmt.Println("wan:", wancfg)
	mode, _ := strconv.Atoi(wancfg.Mode)
	wanJson := WanJson{
		Port:             wancfg.Port,
		WanMode:          mode,
		WanIp:            wancfg.FixIp,
		WanMask:          wancfg.FixMask,
		WanPppoeAccount:  wancfg.PPPoEAccount,
		WanPppoePassword: wancfg.PPPoEPassword,
		WanPrimaryDns:    wancfg.PrimaryDns,
		WanSecondaryDns:  wancfg.SecondaryDns,
		Success:          "true",
	}

	// fmt.Println("wanJson:", wanJson)

	return c.JSON(http.StatusOK, wanJson)
}
func (self DeviceController) UpdateWan(c echo.Context) error {
	name := c.FormValue("modelName")
	// ids := c.FormValue("id")
	// ports := c.FormValue("port")
	//
	// id, _ := strconv.Atoi(ids)
	// port, _ := strconv.Atoi(ports)

	id := goutils.MustInt(c.FormValue("id"))
	port := goutils.MustInt(c.FormValue("port"))

	mode := c.FormValue("wanMode")
	ip := c.FormValue("wanIp")
	wm := c.FormValue("wanMask")
	pa := c.FormValue("wanPppoeAccount")
	pp := c.FormValue("wanPppoePassword")
	sec := c.FormValue("wanSecondaryDns")
	pri := c.FormValue("wanPrimaryDns")
	gw := c.FormValue("wanGateway")

	//md, _ := strconv.Atoi(mode)

	// fmt.Println("port:", port)
	// fmt.Println("ip:", ip)
	// fmt.Println("mask:", wm)
	// fmt.Println("pa:", pa)
	// fmt.Println("pp:", pp)
	// fmt.Println("sec:", sec)
	// fmt.Println("pri:", pri)
	// fmt.Println("mode:", mode)
	// fmt.Println("gw:", gw)

	wan := model.Wan{
		Port:          port,
		Mode:          mode,
		FixIp:         ip,
		FixMask:       wm,
		FixGateway:    gw,
		PPPoEAccount:  pa,
		PPPoEPassword: pp,
		PrimaryDns:    pri,
		SecondaryDns:  sec,
		DeviceRefer:   id,
	}
	model.UpdateWan(wan)

	dev := model.GetDeviceById(id)
	path := RequestUrl(c)

	// fmt.Println("path=", path)
	// fmt.Println("name=", name)
	// fmt.Println("dev=", dev)

	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   name,
		"Device": dev,
		"Code":   5,
	})
}

func (self DeviceController) AddDev(c echo.Context) error {
	path := RequestUrl(c)
	prjs := model.GetProjects()
	return c.Render(http.StatusOK, "device_add.html", echo.Map{
		"Path":     path,
		"Code":     1,
		"Message":  "nothing",
		"Projects": prjs,
	})
}

func (self DeviceController) Save(c echo.Context) error {
	mac := c.FormValue("mac")
	//modelType := c.FormValue("model")
	name := c.FormValue("name")
	//prjId := c.FormValue("project")

	// fmt.Println("mac:", mac)
	// fmt.Println("model:", modelType)
	// fmt.Println("name:", name)
	// fmt.Println("project:", prjId)

	// pid, _ := strconv.Atoi(prjId)
	// mm, _ := strconv.Atoi(modelType)

	mm := goutils.MustInt(c.FormValue("model"))
	pid := goutils.MustInt(c.FormValue("project"))

	dev := model.Device{
		Mac:           mac,
		Name:          name,
		Online:        true,
		Heartbeat:     1231242,
		Version:       "3.16",
		LanIp:         "192.168.48.1",
		LanMask:       "255.255.255.0",
		DhcpEnable:    1,
		DhcpStartIp:   "192.168.48.2",
		DhcpEndIp:     "192.168.48.254",
		DhcpLeaseTime: 3600,
		MultiSsid:     false,
		RfType:        "2.4",
		RfMode:        1,
		RfFreq:        6,
		RfPower:       23,
		RfType5:       "5",
		RfMode5:       5,
		RfFreq5:       0,
		RfPower5:      23,
		ModelType:     mm,
		CloudHost:     HOST,
		CloudToken:    "helloworld",
		ProjectRefer:  pid,
	}

	err := model.AddDevice(dev)

	path := RequestUrl(c)
	prjs := model.GetProjects()
	if err != nil {
		return c.Render(http.StatusOK, "device_add.html", echo.Map{
			"Path":     path,
			"Code":     -1,
			"Message":  "failed",
			"Projects": prjs,
		})
	}

	return c.Render(http.StatusOK, "device_add.html", echo.Map{
		"Path":     path,
		"Code":     0,
		"Message":  "success",
		"Projects": prjs,
	})
}

func (self DeviceController) DeleteDevice(c echo.Context) error {
	id := goutils.MustInt(c.QueryParam("ids"))
	//modelName := c.QueryParam("modelName")
	// fmt.Println("id:", id)
	// fmt.Println("modelName:", modelName)

	// printFormParams(c)

	model.DeleteDeviceById(id)

	// path := RequestUrl(c)
	// fmt.Println("path=", path)
	devs := model.GetDevices()

	for k, v := range devs {
		ssids := model.GetSsidByDeviceId(v.Id)
		for _, s := range ssids {
			//			fmt.Println("ssid=", s)
			devs[k].Ssid = append(devs[k].Ssid, s)
		}
	}

	// fmt.Println("devices:", devs)
	return c.Redirect(http.StatusFound, "/v3/project/device/v_list")
}

func (self DeviceController) ListPost(c echo.Context) error {
	// printFormParams(c)

	cookie, err := c.Cookie("_cookie_page_size")
	if err != nil {
		//return c.String(http.StatusNotFound, "No cookie_page_size found")
		fmt.Println("no cookie_page_size")
	}
	var pageSize, pageNo int
	pageSize, err = strconv.Atoi(cookie.Value)
	if err != nil {
		pageSize = 20
	}
	pageNo, err = strconv.Atoi(c.FormValue("pageNo"))
	if err != nil {
		pageNo = 1
	}

	//fmt.Println("cookie_page_size:", cookie.Value)
	// fmt.Println("pageNo:", pageNo)

	devNum := model.GetTotalDeviceNum()
	pageNum := devNum / pageSize
	if devNum%pageSize > 0 {
		pageNum = pageNum + 1
	}
	// fmt.Println("pageNum:", pageNum)

	path := RequestUrl(c)
	// fmt.Println("path=", path)
	//devs := model.GetDevices()
	devs := model.ListPageNoDevice(pageNo, pageSize)

	for k, v := range devs {
		ssids := model.GetSsidByDeviceId(v.Id)
		for _, s := range ssids {
			//			fmt.Println("ssid=", s)
			devs[k].Ssid = append(devs[k].Ssid, s)
		}
	}

	models := model.GetAllModels()
	//fmt.Println("devices:", devs)
	return c.Render(http.StatusOK, "device_list.html", echo.Map{
		"Path":        path,
		"Devices":     devs,
		"PageNum":     pageNum,
		"PageNo":      pageNo,
		"TotalDevice": devNum,
		"PageSize":    pageSize,
		"Models":      models,
	})
}

func (self DeviceController) ReadConfig(c echo.Context) error {
	name := c.QueryParam("modelName")
	id := goutils.MustInt(c.QueryParam("id"))
	dev := model.GetDeviceById(id)
	path := RequestUrl(c)
	// fmt.Println("path=", path)
	// fmt.Println("name=", name)
	// fmt.Println("dev=", dev)

	//send config_read command to client
	code := rpc.ReadConfig(dev.Mac)

	qos := model.GetQosByDeviceId(id)
	wanQos := model.GetWanQosByQosId(qos.Id)
	//
	// fmt.Println("qos:", qos)
	// fmt.Println("wanQos:", wanQos)

	ssid := model.GetSsidByDeviceId(id)
	md5 := model.GetMd5ByDeviceId(id)
	prjs := model.GetProjects()
	// fmt.Println("ssid:", ssid)
	models := model.GetAllModels()
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":     path,
		"Name":     name,
		"Device":   dev,
		"Qos":      qos,
		"WanQos":   wanQos,
		"SSID":     ssid,
		"Projects": prjs,
		"Models":   models,
		"Md5":      md5,
		"Code":     code,
	})
}

type Restart struct {
	Success bool `json:"success"`
}

func (self DeviceController) Restart(c echo.Context) error {
	id := goutils.MustInt(c.QueryParam("id"))
	dev := model.GetDeviceById(id)

	code := rpc.Restart(dev.Mac)

	restart := Restart{
		Success: code,
	}
	return c.JSON(http.StatusOK, restart)
}

func (self DeviceController) VPN(c echo.Context) error {
	id := goutils.MustInt(c.QueryParam("id"))
	dev := model.GetDeviceById(id)

	code := rpc.VPNNgrok(dev.Mac)

	vpn := Restart{
		Success: code,
	}
	return c.JSON(http.StatusOK, vpn)
}
