package admin

import (
	"github.com/labstack/echo"
	"net/http"
)

func UpgradeInfo(c echo.Context) error {
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
	})
}

func DeleteUpgradeInfo(c echo.Context) error {
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
	})
}
func AddUpgradeInfo(c echo.Context) error {
	return c.Render(http.StatusOK, "firmware_manage.html", map[string]interface{}{
	})
}