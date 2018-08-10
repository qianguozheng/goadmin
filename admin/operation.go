package admin

import (
	//	"fmt"
	"net/http"
	//	"strconv"

	//	"../model"
	"github.com/labstack/echo"
)

func HandleTrustIpsList(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "trustip_list.html", echo.Map{
		"Path": path,
	})
}
func HandleTrustIpsAdd(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "trustip_add.html", echo.Map{
		"Path": path,
	})
}
