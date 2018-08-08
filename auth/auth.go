package auth

import (
	"fmt"
	"net/http"

	"../rpc"
	"github.com/labstack/echo"
)

//func HandleAuth(c echo.Context) error {

//	c.JSON(http.StatusOK)
//}

func HandleNotify(c echo.Context) error {

	res := rpc.Notification("00782fe82e35", "myiPhonemac", 1441)
	fmt.Println("notify:", res)
	return c.JSON(http.StatusOK, res)
}
