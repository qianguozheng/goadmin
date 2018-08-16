package main

import (
	"flag"
	"fmt"

	"html/template"
	"io"

	"net/http"
	"strings"
	"time"

	"./admin"
	"./auth"
	"./model"
	"./rpc"
	"github.com/foolin/echo-template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	portPtr := flag.String("port", ":80", "port to serve the service")

	flag.Parse()

	//RPC Service for Auth
	s := rpc.RPCServerService()
	defer s.Stop()

	go auth.GoAuth()

	fmt.Println("port=", *portPtr)
	fmt.Println("goAdmin standalone web server")
	e := echo.New()

	model.DB = model.InitDB("goadmin.db")
	model.InitModels()
	model.InitUpgrade()
	model.InitAllDeviceConfig()

	//Terminal Free manage
	model.DB2 = model.InitDB2()
	model.AuthUserTest()

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
			"cmpGormID": func(a int, id uint) bool {
				return uint(a) == id
			},
			"timeStr": func(a int64) string {
				t1 := time.Unix(a, 0) //Parse("2016-12-04 15:39:06 +0800 CST")
				return t1.Format("2006-01-02 15:04:05")
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
	//adminGrp.GET("/v3/core/index", homeCtx.HandleHome)
	adminGrp.GET("/v3/core/index", homeCtx.HandleProjectIndex)
	adminGrp.GET("/v3/project/nav", homeCtx.HandleProjectIndex)

	adminGrp.GET("/v3/project/device/v_list", admin.HandleProjectDeviceList)
	adminGrp.POST("/v3/project/device/v_list", admin.HandleProjectDeviceListPost)

	adminGrp.GET("/v3/project/device/v_edit_cfg", admin.HandleProjectDeviceEdit)
	adminGrp.POST("/v3/project/device/o_update", admin.HandleProjectDeviceUpdateEdit)
	adminGrp.POST("/v3/project/device/o_update_config", admin.HandleProjectDeviceUpdateCloud)
	adminGrp.POST("/v3/project/device/o_save_ssid", admin.HandleProjectDeviceUpdateSSID)
	adminGrp.GET("/v3/project/device/v_ajax_edit_ssid", admin.HandleProjectDeviceEditSSID)
	adminGrp.GET("/v3/project/device/v_ajax_restart", admin.HandleProjectDeviceRestart)

	adminGrp.GET("/v3/project/device/v_edit_delete", admin.HandleProjectDeviceDelDev)

	adminGrp.POST("/v3/project/device/v_ajax_read_wan", admin.HandleProjectDeviceGetWan)
	adminGrp.POST("/v3/project/device/o_update_config_wan", admin.HandleProjectDeviceUpdateWan)
	adminGrp.GET("/v3/project/device/o_read_cfg", admin.HandlePrejectReadConfig)
	adminGrp.GET("/v3/project/device/v_ajax_update_cfg", admin.HandleProjectDeviceVPN)
	//v_ajax_update_mutiWan

	//Project Management
	adminGrp.GET("/v3/project/project/v_list", admin.HandleProjectList)
	adminGrp.GET("/v3/project/project/v_add", admin.HandleProjectAdd)
	adminGrp.POST("/v3/project/project/o_save", admin.HandleProjectSave)
	adminGrp.GET("/v3/project/project/o_delete", admin.HandleProjectDelete)
	adminGrp.POST("/v3/project/project/o_delete", admin.HandleProjectDelete)
	adminGrp.GET("/v3/project/project/v_edit", admin.HandleProjectEdit)

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

	//Terminal Free list
	adminGrp.GET("/v3/project/terminal_free/v_list", admin.TerminalList)
	adminGrp.POST("/v3/project/terminal_free/v_list", admin.TerminalListPost)
	adminGrp.GET("/v3/project/terminal_free/v_add", admin.TerminalAdd)
	adminGrp.POST("/v3/project/terminal_free/o_save", admin.TerminalAddPost)
	adminGrp.GET("/v3/project/terminal_free/o_delete", admin.TerminalDelete)
	adminGrp.GET("/v3/project/terminal_free/v_ajax_check_mac", admin.TerminalCheck)

	//Trust Ips

	adminGrp.GET("/v3/project/trust_ip/v_list", admin.HandleTrustIpsList)
	adminGrp.GET("/v3/project/trust_ip/v_add", admin.HandleTrustIpsAdd)
	adminGrp.GET("/v3/project/trust_ip/v_edit", admin.HandleTrustIpsEdit)
	adminGrp.POST("/v3/project/trust_ip/o_save", admin.HandleTrustIpsSave)
	adminGrp.GET("/v3/project/trust_ip/o_delete", admin.HandleTrustIpsDelete)
	adminGrp.POST("/v3/project/trust_ip/o_update", admin.HandleTrustIpsUpdate)

	//Trust Domain
	adminGrp.GET("/v3/project/trust_domain/v_list", admin.HandleTrustDomainsList)
	adminGrp.GET("/v3/project/trust_domain/v_add", admin.HandleTrustDomainsAdd)
	adminGrp.GET("/v3/project/trust_domain/v_edit", admin.HandleTrustDomainsEdit)
	adminGrp.POST("/v3/project/trust_domain/o_save_domain_strategy", admin.HandleTrustDomainsSave)
	adminGrp.GET("/v3/project/trust_domain/o_delete", admin.HandleTrustDomainsDelete)
	adminGrp.POST("/v3/project/trust_domain/o_update", admin.HandleTrustDomainsUpdate)

	//Dns Bogus
	adminGrp.GET("/v3/project/dns_bogus/v_list", admin.HandleDnsBogusList)
	adminGrp.GET("/v3/project/dns_bogus/v_add", admin.HandleDnsBogusAdd)
	adminGrp.GET("/v3/project/dns_bogus/v_edit", admin.HandleDnsBogusEdit)
	adminGrp.POST("/v3/project/dns_bogus/o_save", admin.HandleDnsBogusSave)
	adminGrp.GET("/v3/project/dns_bogus/o_delete", admin.HandleDnsBogusDelete)
	adminGrp.POST("/v3/project/dns_bogus/o_update", admin.HandleDnsBogusUpdate)
	adminGrp.POST("/v3/project/dns_bogus/o_send", admin.HandleDnsBogusSend)
	adminGrp.GET("/v3/project/dns_bogus/v_ajax_edit", admin.HandleDnsBogusAJAX)

	adminGrp.GET("/v3/project/v_search_pwd", admin.HandleSearchPassword)
	adminGrp.POST("/v3/project/v_search_pwd", admin.HandleSearchPasswordPost)

	//Authentication Server
	//	e.GET("/auth", auth.HandleAuth)
	e.GET("/notify", auth.HandleNotify)

	//Operation management

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
