package admin

import (
	"github.com/labstack/echo"
)

func RequestUrl(c echo.Context) string {
	return c.Request().URL.Path
}
