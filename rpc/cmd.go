package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"

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
