package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"html/template"
	"io"

	"net/http"
	"strings"

	"./admin"
	"./auth"
	"./rpc"
	echotemplate "github.com/foolin/echo-template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/qianguozheng/goadmin/config"
	mw "github.com/qianguozheng/goadmin/http/middleware"
	"github.com/qianguozheng/goadmin/model"
)

func main() {

	// portPtr := flag.String("port", ":80", "port to serve the service")
	//
	// flag.Parse()

	//RPC Service for Auth
	s := rpc.RPCServerService()
	defer s.Stop()

	go auth.GoAuth()
	// Pprof("127.0.0.1:8096")
	Pprof(ConfigFile.MustValue("global", "pprof", "127.0.0.1:8096"))

	// fmt.Println("port=", *portPtr)
	fmt.Println("goAdmin standalone web server")
	e := echo.New()

	dbType := ConfigFile.MustValue("global", "db", "sqlite3")
	if strings.Compare(dbType, "mysql") == 0 {
		fmt.Println("use mysql database")
		model.DB = model.InitDB(model.Mysql(), "mysql")
	} else {
		fmt.Println("use sqlite3 database")
		model.DB = model.InitDB("goadmin.db", "sqlite3")
	}

	model.InitModels()
	// model.InitUpgrade()
	// model.InitAllDeviceConfig()

	//Terminal Free manage， auth server
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
	homeGrp := e.Group("")
	admin.RegisterRoutesHome(homeGrp)

	//Routes
	manageGrp := e.Group("/v3/project", mw.NeedLogin)
	admin.RegisterRoutes(manageGrp)

	//Authentication Server
	e.GET("/notify", auth.HandleNotify)

	e.Logger.Fatal(e.Start(getAddr()))
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
	Partials: []string{"public/sidebar", "public/footer", "public/pagenav", "public/header",
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
		"timeStr": func(a int64) string {
			t1 := time.Unix(a, 0) //Parse("2016-12-04 15:39:06 +0800 CST")
			return t1.Format("2006-01-02 15:04:05")
		},
	},
	DisableCache: true,
})

func getAddr() string {
	// host := ConfigFile.MustValue("listen", "host", "")
	// if host == "" {
	// 	global.App.Host = "localhost"
	// } else {
	// 	global.App.Host = host
	// }
	// global.App.Port = ConfigFile.MustValue("listen", "port", "8081")
	// return host + ":" + global.App.Port
	host := ConfigFile.MustValue("listen", "host", "0.0.0.0")
	port := ConfigFile.MustValue("listen", "port", "8081")
	return host + ":" + port
}

func savePid() {
	pidFilename := ROOT + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
	pid := os.Getpid()

	ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0755)
}
