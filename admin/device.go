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
	fmt.Println("ssid:", ssid[0])
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   name,
		"Device": dev,
		"Qos":    qos,
		"WanQos": wanQos,
		"SSID":   ssid[0],
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
	fmt.Println("ssid:", ssid[0])
	return c.Render(http.StatusOK, "device_edit.html", echo.Map{
		"Path":   path,
		"Name":   name,
		"Device": dev,
		"Qos":    qos,
		"WanQos": wanQos,
		"SSID":   ssid[0],
	})
}
