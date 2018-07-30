package admin

import (
	"fmt"
	"strconv"
	"strings"

	"../model"

	"github.com/labstack/echo"
)

func RequestUrl(c echo.Context) string {
	return c.Request().URL.Path
}

func printFormParams(c echo.Context) {
	x, _ := c.FormParams()
	for xx, xv := range x {
		fmt.Println(xx, xv)
	}
}

func getWanQosFormParams(c echo.Context, qosId int) []model.WanQos {
	wanQos := make([]model.WanQos, 5)

	x, _ := c.FormParams()
	for k, v := range x {
		//fmt.Println(xx, xv)
		if strings.Compare(k, "ports") == 0 {
			//ports := strings.Split(v, " ")
			for pk, pv := range v {
				p, _ := strconv.Atoi(pv)
				//				if len(wanQos) < 5 {
				//					var wq model.WanQos
				//					wq.Port = p
				//					wanQos = append(wanQos, wq)
				//				} else {
				wanQos[pk].Port = p
				//				}

			}
		}

		if strings.Compare(k, "txRates") == 0 {
			//txRates := strings.Split(v, " ")
			for pk, pv := range v {
				p, _ := strconv.Atoi(pv)
				//				if len(wanQos) < 5 {
				//					var wq model.WanQos
				//					wq.Up = p
				//					wanQos = append(wanQos, wq)
				//				} else {
				wanQos[pk].Up = p
				//				}

			}
		} else if strings.Compare(k, "rxRates") == 0 {
			//rxRates := strings.Split(v, " ")
			for pk, pv := range v {
				p, _ := strconv.Atoi(pv)
				//				if len(wanQos) < 5 {
				//					var wq model.WanQos
				//					wq.Down = p
				//					wanQos = append(wanQos, wq)
				//				} else {
				wanQos[pk].Down = p
				//				}
			}
		}
	}

	for k, _ := range wanQos {
		wanQos[k].QosRefer = qosId
	}

	fmt.Println("wanQos len:", len(wanQos), wanQos)

	return wanQos
}