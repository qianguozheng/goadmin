package admin

import (
	"encoding/base64"
	"fmt"
	"strings"

	"net/http"
	"time"

	"../model"
	"github.com/labstack/echo"
)

// adminGrp.GET("/home.html", homeCtx.HandleHome)
// adminGrp.GET("/v3/core/index", homeCtx.HandleProjectIndex)
// adminGrp.GET("/v3/project/nav", homeCtx.HandleProjectIndex)

// homeCtx := admin.NewHomeCtx()
// 	e.GET("/index", homeCtx.HandleLogin)
// 	e.GET("/", homeCtx.HandleLogin)
// 	e.GET("/login.html", homeCtx.HandleLogin)
// 	e.POST("/login", homeCtx.HandleLoginPost)
// 	e.GET("/reset.html", homeCtx.HandleReset)

type HomeController struct{}

func (self HomeController) RegisterRoute(g *echo.Group) {
	// TODO: implement /v3/core/logout, /v3/core/personal/v_edit?modelName=edit
	g.GET("/index", self.Login)
	g.GET("/", self.Login)
	g.GET("/login.html", self.Login)
	g.POST("/login", self.LoginPost)
	g.GET("/reset.html", self.Reset)

	g.GET("/home.html", self.Home, checkCookie)
	g.GET("/v3/core/index", self.ProjectIndex, checkCookie)
	g.GET("/v3/project/nav", self.ProjectIndex, checkCookie)
}

// func NewHomeCtx() *HomeController {
// 	home := HomeCtx{}
// 	return &home
// }

func (home HomeController) Login(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "login.html", echo.Map{
		"Path": path,
	})
}

func (home HomeController) LoginPost(c echo.Context) error {

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
	cookie.Value = base64.StdEncoding.EncodeToString([]byte(userName + password + "honeybot"))
	cookie.Expires = time.Now().Add(24 * time.Hour)
	//Store cookie value to db
	model.SetCookie(userName, cookie.Value)
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/home.html")
}

func (home HomeController) Reset(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "reset.html", echo.Map{
		"Path": path,
	})
}

func (home HomeController) Home(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "home.html", echo.Map{
		"Path": path,
	})
}

func (home HomeController) ProjectIndex(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "project_index.html", echo.Map{
		"Path": path,
	})
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionId")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.Redirect(http.StatusFound, "/login.html")
				//return c.String(http.StatusUnauthorized, "You don't have the right cookie")
			}
			fmt.Println("cookie not present")
			return err
		}

		if true == model.GetCookie(cookie.Value) {

			return next(c)
		}

		return c.Redirect(http.StatusFound, "/login.html")
	}
}
