package main

import (
	"flag"
	"fmt"

	"html/template"
	"io"

	"./admin"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	portPtr := flag.String("port", ":80", "port to serve the service")

	flag.Parse()

	fmt.Println("port=", *portPtr)
	fmt.Println("goAdmin standalone web server")
	e := echo.New()

	render := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("html/*.html")),
	}
	e.Renderer = render

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/", "static/")

	//Homepage
	homeCtx := admin.NewHomeCtx()
	e.GET("/index", homeCtx.HandleLogin)
	e.GET("/", homeCtx.HandleLogin)
	e.GET("/login.html", homeCtx.HandleLogin)
	e.GET("/reset.html", homeCtx.HandleReset)
	e.POST("/login", homeCtx.HandleLoginPost)
	//After Login
	e.GET("/home.html", homeCtx.HandleHome)
	e.GET("/v3/project/nav", homeCtx.HandleProjectIndex)

	//API
	adm := e.Group("/admin/")
	{
		adminCtx := admin.NewAdminCtx()
		adm.GET("/admin/", adminCtx.Handle)

		//Upload file
		//adm.GET("/admin/upload", adminCtx.HandleUpload)
		//adm.POST("/admin/upload", adminCtx.HandleUpload)
	}

	e.Logger.Fatal(e.Start(*portPtr))
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
