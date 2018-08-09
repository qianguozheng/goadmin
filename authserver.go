package main

import (
	"./model"
)

func goAuth() {
	e := echo.New()

	e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
		Root:      "html_auth",
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
		},
		DisableCache: true,
	})

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/assets", "static/assets")
	manage := e.Group("")
	manage.Use(checkAuthCookie)
}

func checkAuthCookie(next echo.HandlerFunc) echo.HandlerFunc {
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

		if true == model.GetAuthCookie(cookie.Value) {

			return next(c)
		}

		return c.Redirect(http.StatusFound, "/auth.html")
	}
}
