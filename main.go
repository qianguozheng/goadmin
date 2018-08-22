package main

import (
	"flag"
	"fmt"
	"time"

	"html/template"
	"io"

	"net/http"
	"strings"

	"./admin"
	"./auth"
	"./model"
	"./rpc"
	echotemplate "github.com/foolin/echo-template"
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
	Pprof("127.0.0.1:8096")

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
	e.POST("/login", homeCtx.HandleLoginPost)
	e.GET("/reset.html", homeCtx.HandleReset)

	//Routes
	newGrp := e.Group("/v3/project", checkCookie)
	admin.RegisterRoutes(newGrp)

	adminGrp := e.Group("")
	//adminGrp.Static("/", "static/")
	adminGrp.Use(checkCookie)
	adminGrp.GET("/home.html", homeCtx.HandleHome)
	adminGrp.GET("/v3/core/index", homeCtx.HandleProjectIndex)
	adminGrp.GET("/v3/project/nav", homeCtx.HandleProjectIndex)

	//Authentication Server
	e.GET("/notify", auth.HandleNotify)

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

var render = echotemplate.New(echotemplate.TemplateConfig{
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
