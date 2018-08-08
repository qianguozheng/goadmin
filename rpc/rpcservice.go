package rpc

import (
	"fmt"

	"github.com/valyala/gorpc"
)

func DispatcherAddFunc() *gorpc.Dispatcher {
	d := gorpc.NewDispatcher()

	d.AddFunc("SendCommand", func(request string) (string, error) {
		//Parse Json
		return "", nil
	})
	d.AddFunc("Verification", func(request string) (string, error) {
		return Verification(request)
	})

	//TODO: Add RPC FUNC here,

	return d

}

func RPCClientRequest(msg string) (interface{}, error) {
	d := DispatcherAddFunc()
	//Client
	c := gorpc.NewTCPClient("127.0.0.1:12445")
	c.Start()
	defer c.Stop()

	dc := d.NewFuncClient(c)

	res, err := dc.Call("SendCommand", msg)
	fmt.Printf("res:%v, err:%v\n", res, err)

	//TODO: process response

	return res, err
}

//Purpose: used for AAA authentication
//Command: verification
func RPCServerService() *gorpc.Server {
	d := DispatcherAddFunc()

	//Start TCP Server
	s := gorpc.NewTCPServer("127.0.0.1:12448", d.NewHandlerFunc())
	if err := s.Start(); err != nil {
		fmt.Println("Cannot start rpc server")
	}
	//defer s.Stop()
	return s
}

//func main() {
//	s := RPCServerService()
//	defer s.Stop()

//	RPCClientRequest()

//}
