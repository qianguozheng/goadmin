package admin

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type AdminCtx struct{}

func NewAdminCtx() *AdminCtx {
	admin := AdminCtx{}
	return &admin
}

func (admin AdminCtx) Handle(c echo.Context) error {
	fmt.Println("admin page")

	path := RequestUrl
	//for _, cookie := range c.Cookies() {
	//	fmt.Printf("Name:[%s][%s]\n", cookie.Name, cookie.Value)
	//}

	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return err
	}
	fmt.Printf("sessionId:%s\n", cookie.Value)
	//TODO: judge the cookie Name and Value exist in db

	//return c.File("html/admin.html")
	return c.Render(http.StatusOK, "admin.html", echo.Map{
		"name": "ooooooo",
		"pass": "fuck",
		"Path": path,
	})
}
