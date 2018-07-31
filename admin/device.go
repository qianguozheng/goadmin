package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"github.com/labstack/echo"
)

func HandleProjectDeviceList(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)
	devs := model.GetDevices()

	for k, v := range devs {
		ssids := model.GetSsidByDeviceId(v.Id)
		for _, s := range ssids {
			fmt.Println("ssid=", s)
			//v.Ssid = make([]model.Ssid, 0)
			//v.Ssid = append(v.Ssid, s)
			devs[k].Ssid = append(devs[k].Ssid, s)
		}
	}

	fmt.Println("devices:", devs)
	return c.Render(http.StatusOK, "device_list.html", echo.Map{
		"Path":    path,
		"Devices": devs,
	})
}

func HandleProjectDeviceOffline(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "device_offline.html", echo.Map{
		"Path": path,
	})
}

func HandleProjectDeviceEdit(c echo.Context) error {
	name := c.QueryParam("modelName")
	ids := c.QueryParam("id")
	id, _ := strconv.Atoi(ids)
	dev := model.GetDeviceById(id)
	path := RequestUrl(c)
	fmt.Println("path=", path)
	fmt.Println("name=", name)
	fmt.Println("dev=", dev)

	qos := model.QueryQos(id)
	wanQos := model.QueryWanQos(qos.Id)

	fmt.Println("qos:", qos)
	fmt.Println("wanQos:", wanQos)

	ssid := model.GetSsidByDeviceId(id)
	fmt.Println("ssid:", ssid)
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   name,
		"Device": dev,
		"Qos":    qos,
		"WanQos": wanQos,
		"SSID":   ssid,
	})
}

//post list_edit
func HandleProjectDeviceUpdateEdit(c echo.Context) error {
	ids := c.FormValue("id")
	name := c.FormValue("name")
	state := c.FormValue("status")
	mode := c.FormValue("mode")

	//save data to database
	id, _ := strconv.Atoi(ids)
	modei, _ := strconv.Atoi(mode)
	status, _ := strconv.Atoi(state)
	dev := model.GetDeviceById(id)
	dev.Name = name
	dev.Status = status
	dev.Mode = modei

	fmt.Println("id=", id)
	fmt.Println("name=", name)
	fmt.Println("state=", state)
	fmt.Println("mode=", mode)
	fmt.Println("dev=", dev)

	model.UpdateDeviceById(dev)

	path := RequestUrl(c)
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   "list_edit",
		"Device": dev,
	})
}

//post list_cloud
func HandleProjectDeviceUpdateCloud(c echo.Context) error {
	//list_cloud
	var page_name string
	page_name = c.FormValue("modelName")
	ids := c.FormValue("deviceId")
	token := c.FormValue("ccToken")
	port := c.FormValue("ccPort")
	host := c.FormValue("ccHost")
	id, _ := strconv.Atoi(ids)

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
		fmt.Println("list lan")
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
		qos = model.QueryQos(id)
		qos.UpRate = up
		qos.DownRate = down
		qos.TcpLimit = tcp
		qos.UdpLimit = udp

		model.UpdateQos(qos)
		//WanQos table
		printFormParams(c)
		wanQos = getWanQosFormParams(c, qos.Id)
		model.UpdateWanQoss(wanQos)

		fmt.Println("qos:", qos)
		fmt.Println("wanQos:", wanQos)

	}

	///////////////////////////  Update device info   ////////////////////////
	model.UpdateDeviceById(dev)

	fmt.Println(dev)
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   page_name,
		"Device": dev,
		"Qos":    qos,
		"WanQos": wanQos,
	})
}

func HandleProjectDeviceUpdateSSID(c echo.Context) error {
	printFormParams(c)

	//page_name := c.QueryParam("modelName")
	ports := c.FormValue("port")
	port, _ := strconv.Atoi(ports)

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
	fmt.Println("set ssid:", s)
	model.UpdateSsid(s)

	//pass param to template
	ssid := model.GetSsidByDeviceId(devId)

	dev := model.GetDeviceById(devId)
	path := RequestUrl(c)

	fmt.Println("path:", path)
	//fmt.Println("name:", page_name)
	fmt.Println("dev:", dev)
	fmt.Println("ssid:", ssid)

	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   "list_ssid",
		"Device": dev,
		"SSID":   ssid,
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

func HandleProjectDeviceEditSSID(c echo.Context) error {
	ids := c.FormValue("deviceId")
	ports := c.FormValue("port")
	port, _ := strconv.Atoi(ports)
	id, _ := strconv.Atoi(ids)

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
	fmt.Println("ssid:", ssid)
	fmt.Println("ssidJson:", s)
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
	WanPppoeAccount  string `json:wanPppoeAccount`
	WanPppoePassword string `json:wanPppoePassword`
	WanPrimaryDns    string `json:"wanPrimaryDns"`
	WanSecondaryDns  string `json:"wanSecondaryDns"`
	Success          string `json:"success"`
}

func HandleProjectDeviceGetWan(c echo.Context) error {
	ids := c.FormValue("id")
	ports := c.FormValue("port")
	id, _ := strconv.Atoi(ids)
	port, _ := strconv.Atoi(ports)
	wan := model.Wan{
		Port:        port,
		DeviceRefer: id,
	}
	wancfg := model.QueryWan(wan)
	fmt.Println("wan:", wancfg)
	wanJson := WanJson{
		Port:             wancfg.Port,
		WanMode:          wancfg.Mode,
		WanIp:            wancfg.FixIp,
		WanMask:          wancfg.FixMask,
		WanPppoeAccount:  wancfg.PPPoEAccount,
		WanPppoePassword: wancfg.PPPoEPassword,
		WanPrimaryDns:    wancfg.PrimaryDns,
		WanSecondaryDns:  wancfg.SecondaryDns,
		Success:          "true",
	}

	fmt.Println("wanJson:", wanJson)

	return c.JSON(http.StatusOK, wanJson)
}
func HandleProjectDeviceUpdateWan(c echo.Context) error {
	name := c.FormValue("modelName")
	ids := c.FormValue("id")
	ports := c.FormValue("port")

	id, _ := strconv.Atoi(ids)
	port, _ := strconv.Atoi(ports)

	mode := c.FormValue("wanMode")
	ip := c.FormValue("wanIp")
	wm := c.FormValue("wanMask")
	pa := c.FormValue("wanPppoeAccount")
	pp := c.FormValue("wanPppoePassword")
	sec := c.FormValue("wanSecondaryDns")
	pri := c.FormValue("wanPrimaryDns")
	gw := c.FormValue("wanGateway")

	md, _ := strconv.Atoi(mode)

	fmt.Println("port:", port)
	fmt.Println("ip:", ip)
	fmt.Println("mask:", wm)
	fmt.Println("pa:", pa)
	fmt.Println("pp:", pp)
	fmt.Println("sec:", sec)
	fmt.Println("pri:", pri)
	fmt.Println("mode:", mode)
	fmt.Println("gw:", gw)

	wan := model.Wan{
		Port:          port,
		Mode:          md,
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

	fmt.Println("path=", path)
	fmt.Println("name=", name)
	fmt.Println("dev=", dev)

	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   name,
		"Device": dev,
	})
}

func HandleProjectAddDev(c echo.Context) error {
	path := RequestUrl(c)

	return c.Render(http.StatusOK, "device_add.html", echo.Map{
		"Path": path,
	})
}

func HandleProjectAddDevSave(c echo.Context) error {
	mac := c.FormValue("mac")
	modelType := c.FormValue("model")
	name := c.FormValue("name")

	fmt.Println("mac:", mac)
	fmt.Println("model:", modelType)
	fmt.Println("name:", name)

	mm, _ := strconv.Atoi(modelType)

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
	}

	err := model.AddDeviceMac(dev)

	path := RequestUrl(c)
	if err != nil {
		return c.Render(http.StatusOK, "device_add.html", echo.Map{
			"Path":   path,
			"JSCode": "alert(\"failed\")",
		})
	}

	return c.Render(http.StatusOK, "device_add.html", echo.Map{
		"Path":   path,
		"JSCode": "alert(\"success\")",
	})
}
