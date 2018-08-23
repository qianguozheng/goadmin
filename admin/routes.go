package admin

import "github.com/labstack/echo"

func RegisterRoutes(g *echo.Group) {
	new(ProjectController).RegisterRoute(g)
	new(UpgradeController).RegisterRoute(g)
	new(TerminalController).RegisterRoute(g)
	new(PasswordController).RegisterRoute(g)
	new(DnsBogusController).RegisterRoute(g)
	new(TrustIpsController).RegisterRoute(g)
	new(TrustDomainsController).RegisterRoute(g)
	new(DeviceController).RegisterRoute(g)
}

func RegisterRoutesHome(g *echo.Group) {
	new(HomeController).RegisterRoute(g)
}
