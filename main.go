package main

import (
	"flag"
	"fmt"

	"html/template"
	"io"

	"./admin"
	"./model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strings"
	"net/http"
)

func main() {

	portPtr := flag.String("port", ":80", "port to serve the service")

	flag.Parse()

	fmt.Println("port=", *portPtr)
	fmt.Println("goAdmin standalone web server")
	e := echo.New()

	model.DB = model.InitDB()

	render := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}
	e.Renderer = render

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Basic Auth no need login page
	//e.Use(middleware.BasicAuth(func(username,password string, c echo.Context)(bool){
	//	found := model.GetUser(username, password)
	//	if found {
	//		return true
	//	} else {
	//		return false
	//	}
	//}))

	e.Static("/assets", "static/assets")

	//Homepage
	homeCtx := admin.NewHomeCtx()
	e.GET("/index", homeCtx.HandleLogin)
	e.GET("/", homeCtx.HandleLogin)
	e.GET("/login.html", homeCtx.HandleLogin)


	e.GET("/reset.html", homeCtx.HandleReset)
	e.POST("/login", homeCtx.HandleLoginPost)

	//After Login
	adminGrp := e.Group("")
	//adminGrp.Static("/", "static/")
	adminGrp.Use(checkCookie)
	adminGrp.GET("/home.html", homeCtx.HandleHome)
	adminGrp.GET("/v3/core/index", homeCtx.HandleHome)
	adminGrp.GET("/v3/project/nav", homeCtx.HandleProjectIndex)
	adminGrp.GET("/v3/project/device/v_list", homeCtx.HandleProjectDeviceList)
	adminGrp.GET("/v3/project/device_offline/v_list_period", homeCtx.HandleProjectDeviceOffline)
	adminGrp.GET("/v3/project/firmware/v_list", homeCtx.HandleProjectUpgradeManage)

	//Test JWT
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", homeCtx.HandleRestricted)

	e.Logger.Fatal(e.Start(*portPtr))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//Check Cookie

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionId")

		if err != nil {
			if (strings.Contains(err.Error(),"named cookie not present")) {
				return c.Redirect(http.StatusFound, "/login.html")
				//return c.String(http.StatusUnauthorized, "You don't have the right cookie")
			}
			fmt.Println("cookie not present")
			return err
		}

		if true == model.GetCookie(cookie.Value) {

			return next(c)
		}

		//if cookie.Value == "some_string" {
		//
		//}
		//return c.String(http.StatusUnauthorized, "you a not login in")
		return c.Redirect(http.StatusFound, "/login.html")
	}
}