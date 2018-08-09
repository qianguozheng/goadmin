package auth

import (
	"fmt"
	"strings"

	"html/template"
	"net/http"

	"../model"
	"github.com/foolin/echo-template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func GoAuth() {
	e := echo.New()

	e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
		Root:      "html_auth",
		Extension: ".html",
		Master:    "",
		Partials:  []string{},
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
		},
		DisableCache: true,
	})

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/assets", "static/assets")

	e.GET("/", HandleAuthIndex)
	e.GET("/auth/register", HandleAuthRegister)
	e.POST("/auth/register", HandleAuthRegisterPost)
	e.GET("/auth/v_ajax_check_username", HandleAuthRegisterCheck)
	e.GET("/auth/login", HandleAuthLogin)
	e.POST("/auth/login", HandleAuthLoginPost)

	manage := e.Group("")
	manage.Use(checkAuthCookie)
	manage.GET("/auth/manage.html", HandleManage)

	user := e.Group("")
	user.Use(checkAuthCookie)
	user.GET("/auth/user.html", HandleUserPage)

	e.Logger.Fatal(e.Start(":9000"))
}

func checkAuthCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionId")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.Redirect(http.StatusFound, "/auth/login")
				//return c.String(http.StatusUnauthorized, "You don't have the right cookie")
			}
			fmt.Println("cookie not present")
			return err
		}

		if true == model.GetAuthCookie(cookie.Value) {

			return next(c)
		}

		return c.Redirect(http.StatusFound, "/auth/login")
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
