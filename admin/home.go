package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type HomeCtx struct{}

func NewHomeCtx() *HomeCtx {
	home := HomeCtx{}
	return &home
}

func (home HomeCtx) Handle(c echo.Context) error {
	//return c.String(http.StatusOK, "")

	//return c.File("html/index.html")
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"accessNumber": 1,
	})
}

func (home HomeCtx) HandleTheme(c echo.Context) error {
	//return c.File("static/theme/index.html")

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"accessNumber": 0,
	})
}

func (home HomeCtx) HandleLogin(c echo.Context) error {

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"accessNumber": 0,
	})
}

func (home HomeCtx) HandleLoginPost(c echo.Context) error {

	userName := c.FormValue("username")
	password := c.FormValue("password")

	fmt.Printf("username=%s, password=%s\n", userName, password)

	//TODO: store cookie into database
	cookie := new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = userName
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	//TODO: judge username and password

	//	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
	//		"accessNumber": 0,
	//	})
	return c.Redirect(http.StatusFound, "/home.html")
}

func (home HomeCtx) HandleReset(c echo.Context) error {

	return c.Render(http.StatusOK, "reset.html", map[string]interface{}{
		"accessNumber": 0,
	})
}

func (home HomeCtx) HandleHome(c echo.Context) error {

	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"accessNumber": 0,
	})
}

func (home HomeCtx) HandleProjectIndex(c echo.Context) error {

	return c.Render(http.StatusOK, "project_index.html", map[string]interface{}{
		"accessNumber": 0,
	})
}
