package main

import (
	"flag"
	"fmt"

	"html/template"
	"io"

	"net/http"
	"strings"

	"./admin"
	"./model"
	"github.com/foolin/echo-template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	portPtr := flag.String("port", ":80", "port to serve the service")

	flag.Parse()

	fmt.Println("port=", *portPtr)
	fmt.Println("goAdmin standalone web server")
	e := echo.New()

	model.DB = model.InitDB()
	model.InitModels()
	model.InitUpgrade()
	model.InitAllDeviceConfig()

	//Old format template
	//	render := &TemplateRenderer{
	//		templates: template.Must(template.ParseGlob("html/*.html")),
	//	}
	//e.Renderer = render

	e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
		Root:      "html",
		Extension: ".html",
		Master:    "",
		Partials: []string{"public/sidebar", "public/footer",
			"device/device_config", "device/list_cloud", "device/list_wan", "device/list_lan",
			"device/list_edit", "device/list_rf", "device/list_ssid", "device/list_vpn",
			"device/list_qos", "device/list_dhcp"},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
			"contains": func(a, b string) bool {
				return strings.Contains(a, b)
			},
			"var": newVar,
			"set": setVar,
			"cmp": func(x *interface{}, e int) bool {
				return (*x).(int) == e
			},
			"dec": func(a int) int {
				return (a - 1)
			},
		},
		DisableCache: true,
	})

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
	e.POST("/login", homeCtx.HandleLoginPost)

	e.GET("/reset.html", homeCtx.HandleReset)

	//After Login
	adminGrp := e.Group("")
	//adminGrp.Static("/", "static/")
	adminGrp.Use(checkCookie)
	adminGrp.GET("/home.html", homeCtx.HandleHome)
	adminGrp.GET("/v3/core/index", homeCtx.HandleHome)
	adminGrp.GET("/v3/project/nav", homeCtx.HandleProjectIndex)

	adminGrp.GET("/v3/project/device/v_list", admin.HandleProjectDeviceList)
	adminGrp.GET("/v3/project/device/v_edit_cfg", admin.HandleProjectDeviceEdit)
	adminGrp.POST("/v3/project/device/o_update", admin.HandleProjectDeviceUpdateEdit)
	adminGrp.POST("/v3/project/device/o_update_config", admin.HandleProjectDeviceUpdateCloud)
	adminGrp.POST("/v3/project/device/o_save_ssid", admin.HandleProjectDeviceUpdateSSID)
	adminGrp.GET("/v3/project/device/v_ajax_edit_ssid", admin.HandleProjectDeviceEditSSID)

	adminGrp.POST("/v3/project/device/v_ajax_read_wan", admin.HandleProjectDeviceGetWan)
	adminGrp.POST("/v3/project/device/o_update_config_wan", admin.HandleProjectDeviceUpdateWan)
	//v_ajax_update_mutiWan

	//
	adminGrp.GET("/v3/project/adddev/v_list", admin.HandleProjectAddDev)

	adminGrp.POST("/v3/project/adddev/o_save", admin.HandleProjectAddDevSave)

	adminGrp.GET("/v3/project/device_offline/v_list_period", admin.HandleProjectDeviceOffline)

	adminGrp.GET("/v3/project/firmware/v_list", admin.HandleProjectUpgradeManage)
	adminGrp.GET("/v3/project/firmware/v_add", admin.HandleProjectUpgradeAdd)
	adminGrp.POST("/v3/project/firmware/o_save", admin.HandleProjectUpgradeSave)
	adminGrp.GET("/v3/project/firmware/v_edit", admin.HandleProjectUpgradeEdit)
	adminGrp.POST("/v3/project/firmware/v_delete", admin.HandleProjectUpgradeDelete)
	adminGrp.GET("/v3/project/firmware/o_delete", admin.HandleProjectUpgradeDelete)
	adminGrp.POST("/v3/project/firmware/o_delete", admin.HandleProjectUpgradeDelete)

	adminGrp.POST("/v3/project/firmware/o_update", admin.HandleProjectUpgradeUpdate)

	//Test JWT
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", homeCtx.HandleRestricted)

	//Example of html template usage
	e.GET("/example", admin.HandleExample)
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

// Add global variable for template
func newVar(v interface{}) (*interface{}, error) {
	x := interface{}(v)
	return &x, nil
}
func setVar(x *interface{}, v interface{}) (string, error) {
	*x = v
	return "", nil
}
