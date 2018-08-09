package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"../model"
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

func HandleAuthIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "auth.html", echo.Map{})
}

/*
	Register
*/
func HandleAuthRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", echo.Map{})
}

func HandleAuthRegisterPost(c echo.Context) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")

	model.AddAuthUser(user, pass)

	cookie := new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(user + pass + "honeybot"))
	cookie.Expires = time.Now().Add(24 * time.Hour)
	//Store cookie value to db=
	model.SetAuthCookie(user, cookie.Value)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/auth/user.html")
}
func HandleAuthRegisterCheck(c echo.Context) error {
	user := c.QueryParam("user")

	result := model.CheckUserExist(user)

	var s string
	if result {
		s = "false"
	} else {
		s = "true"
	}
	return c.String(http.StatusOK, s)
}

/*
	Login
*/
func HandleAuthLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", echo.Map{})
}

func HandleAuthLoginPost(c echo.Context) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")

	cookie := new(http.Cookie)
	cookie.Name = "sessionId"
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(user + pass + "honeybot"))
	cookie.Expires = time.Now().Add(24 * time.Hour)
	//Store cookie value to db=
	fmt.Println("save cookie to db")
	model.SetAuthCookie(user, cookie.Value)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/auth/user.html")
}

/*
	User
*/

func HandleUserPage(c echo.Context) error {
	return c.String(http.StatusOK, "login ok")
}

/*
	Manage
*/
func HandleManage(c echo.Context) error {
	aus := model.GetAllAuthUsers()
	return c.Render(http.StatusOK, "manage.html", echo.Map{
		"Users": aus,
	})
}
