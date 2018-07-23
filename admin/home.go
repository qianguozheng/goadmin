package admin

import (
	//"fmt"
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"encoding/base64"
	"../model"
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

func (home HomeCtx) HandleRestricted(c echo.Context) error{
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
func (home HomeCtx) HandleLoginPost(c echo.Context) error {

	userName := c.FormValue("username")
	password := c.FormValue("password")

	//if userName == "user" && password == "pass" {
	//	//Create token
	//	token:=jwt.New(jwt.SigningMethodHS256)
	//
	//	//Set claims
	//	claims := token.Claims.(jwt.MapClaims)
	//	claims["name"] = "Jon Snow"
	//	claims["admin"] = true
	//	claims["exp"] = time.Now().Add(time.Hour*72).Unix()
	//
	//	//Generate encoded token and send it as response
	//	t, err := token.SignedString([]byte("secret"))
	//	if err != nil {
	//		return err
	//	}
	//
	//	c.Response().Header().Set("Authorization", "Bearer " + t)
	//	//c.Response().WriteHeader(201)
	//	//return c.JSON(http.StatusOK, map[string]string{
	//	//	"token":t,
	//	//})
	//	return c.Redirect(http.StatusFound, "/restricted")
	//}
	//return echo.ErrUnauthorized

	//fmt.Printf("username=%s, password=%s\n", userName, password)
	//
	////TODO: store cookie into database

	cookie := new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(userName+password+"honeybot"))
	cookie.Expires = time.Now().Add(24 * time.Hour)
	//Store cookie value to db
	model.SetCookie(userName, cookie.Value)
	c.SetCookie(cookie)

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
	})
}

func (home HomeCtx) HandleProjectDeviceList(c echo.Context) error {
	return c.Render(http.StatusOK, "device_list.html", map[string]interface{}{
	})
}

func (home HomeCtx) HandleProjectDeviceOffline (c echo.Context) error {
	return c.Render(http.StatusOK, "device_offline.html", map[string]interface{}{
	})
}

func (home HomeCtx) HandleProjectUpgradeManage (c echo.Context) error {
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
	})
}
