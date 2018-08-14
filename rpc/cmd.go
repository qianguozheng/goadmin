package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"../model"

	"github.com/qianguozheng/gotcpserver/proto"
)

//1. Config Read
func ReadConfig(mac string) int {
	//compose config read request

	req := proto.ReqParam{
		Cmd:   "config_read_req",
		SeqId: "unique id",
		Mac:   mac,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("form reqJson failed", err.Error())
	}

	//send to rpc server
	res, err := RPCClientRequest(string(reqJson))
	if err != nil {
		return -1
	}

	if res != nil {
		fmt.Println("rpc resp:", res.(string))

		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(res.(string)), &msg)
		if err != nil {
			return -1
		}
		if msg["cmd"] != nil && msg["cmd"].(string) != "config_read_resp" {
			return -1
		}
	}

	//Parse response
	return 0
}

//2. Reboot device
func Restart(mac string) bool {

	req := proto.ReqParam{
		Cmd:   "reboot_req",
		SeqId: "unique id",
		Mac:   mac,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("form reqJson failed", err.Error())
	}

	//send to rpc server
	res, err := RPCClientRequest(string(reqJson))
	if err != nil {
		return false
	}

	if res != nil {
		fmt.Println("rpc resp:", res.(string))
		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(res.(string)), &msg)
		if err != nil {
			return false
		}
		if msg["cmd"] != nil && msg["cmd"].(string) != "reboot_resp" {
			return false
		}
	}

	//Parse response
	return true
}

//3. Upgrade
func Upgrade(mac, url, ver, md5 string) bool {

	vi, err := strconv.Atoi(ver)
	if err != nil {
		return false
	}
	req := proto.Upgrade{
		Cmd:   "upgrade_req",
		SeqId: "unique id",
		Mac:   mac,
		Ver:   vi,
		Md5:   md5,
		Url:   url,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("form reqJson failed", err.Error())
	}

	//send to rpc server
	res, err := RPCClientRequest(string(reqJson))
	if err != nil {
		return false
	}

	if res != nil {
		fmt.Println("rpc resp:", res.(string))
		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(res.(string)), &msg)
		if err != nil {
			return false
		}
		if msg["cmd"] != nil && msg["cmd"].(string) != "upgrade_resp" {
			return false
		}
	}

	//Parse response
	return true
}

//4. Notification
func Notification(mac, tmac string, valid int) bool {

	req := proto.Notification{
		Cmd:         "notification_req",
		SeqId:       "unique id",
		Mac:         mac,
		TerminalMac: tmac,
		Valid:       valid,
		AuthType:    0,
		AuthId:      "13538273761",
		UpRate:      256,
		DownRate:    4096,
		TcpLimit:    200,
		UdpLimit:    200,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("form reqJson failed", err.Error())
	}

	//send to rpc server
	res, err := RPCClientRequest(string(reqJson))
	if err != nil {
		return false
	}

	if res != nil {
		fmt.Println("rpc resp:", res.(string))
		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(res.(string)), &msg)
		if err != nil {
			return false
		}
		if msg["cmd"] != nil && msg["cmd"].(string) != "notification_resp" {
			return false
		}
	}

	//Parse response
	return true
}

//5. Verification
func Verification(request string) (string, error) {
	var req proto.VerificationReq
	err := json.Unmarshal([]byte(request), &req)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Println("Verification:", req)

	//free terminal
	free := model.CheckTermFreeByMac(req.TerminalMac)
	var valid int
	if free {
		valid = 1440
	} else {
		valid = 0
	}

	resp := proto.Verification{
		Cmd:   "verification_resp",
		SeqId: "unique id",
		Code:  "000",
		Data: proto.VerificationData{
			TerminalMac: req.TerminalMac,
			Valid:       valid,
			AuthType:    0,
			AuthId:      "13538273761",
			UpRate:      256,
			DownRate:    4096,
			TcpLimit:    200,
			UdpLimit:    200,
		},
	}

	msg, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	//TODO: Check user's mac usner device'mac and timeout

	return string(msg), nil
}

//6. VPN
func VPNNgrok(mac string) bool {

	req := proto.ReqParam{
		Cmd:   "rc_write_req",
		SeqId: "unique id",
		Mac:   mac,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("form reqJson failed", err.Error())
	}

	//send to rpc server
	res, err := RPCClientRequest(string(reqJson))
	if err != nil {
		return false
	}

	if res != nil {
		fmt.Println("rpc resp:", res.(string))
		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(res.(string)), &msg)
		if err != nil {
			return false
		}
		if msg["cmd"] != nil && msg["cmd"].(string) != "rc_write_resp" {
			return false
		}
	}

	//Parse response
	return true
}

//7. Dns Bogus: request
const (
	ProjectKind = iota
	DeviceKind
)

type DnsBogusParam struct {
	Mac string
	Id  int
}

func DnsBogusRequest(kind int, param interface{}) bool {
	resp := make(chan bool)
	inner := func(mac string, dns model.DnsBogus, resp chan bool) {
		dnsBogus := proto.DnsBogusWrite{
			Cmd:   "dns_bogus_write_req",
			SeqId: "unique id",
			Mac:   mac,
			Bogus: []proto.DnsBogus{
				{
					Domain: dns.Domain,
					Host:   dns.Ip,
				},
			},
		}
		reqJson, err := json.Marshal(dnsBogus)
		//fmt.Println("reqJson:", string(reqJson))
		if err != nil {
			fmt.Println("form reqJson failed", err.Error())
		}

		//send to rpc server
		res, err := RPCClientRequest(string(reqJson))
		if err != nil {
			resp <- false

		}

		if res != nil {
			fmt.Println("rpc resp:", res.(string))
			msg := make(map[string]interface{})
			err = json.Unmarshal([]byte(res.(string)), &msg)
			if err != nil {
				resp <- false

			}
			if msg["cmd"] != nil && msg["cmd"].(string) != "dns_bogus_write_resp" {
				resp <- false

			}
		}
		resp <- true

	}

	//	fmt.Println("Param:", param)
	//	fmt.Println("Kind:", kind)
	//	fmt.Println("ProjectKind:", ProjectKind)
	//	fmt.Println("DeviceKind:", DeviceKind)
	if kind == ProjectKind {
		dnsBogusId := param.(DnsBogusParam).Id
		//Send to each mac
		var device []model.Device

		dnsBogus := model.GetDnsBogusById(dnsBogusId)
		device = model.GetDeviceByProjectId(dnsBogus.ProjectRefer)

		for _, v := range device {
			mac := v.Mac
			go inner(mac, dnsBogus, resp)
		}

	} else if kind == DeviceKind {
		mac := param.(DnsBogusParam).Mac
		dnsBogusId := param.(DnsBogusParam).Id
		dnsBogus := model.GetDnsBogusById(dnsBogusId)

		go inner(mac, dnsBogus, resp)

	} else {
		fmt.Println("Not specified type")
		return false
	}

	select {
	case result := <-resp:
		fmt.Println("get response", result)
	case <-time.After(5 * time.Second):
		fmt.Println("timeout for get response")
	}
	return true
}
