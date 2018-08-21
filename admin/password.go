package admin

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// adminGrp.GET("/v3/project/v_search_pwd", admin.HandleSearchPassword)
// adminGrp.POST("/v3/project/v_search_pwd", admin.HandleSearchPasswordPost)
type PasswordController struct{}

func (self PasswordController) RegisterRoute(g *echo.Group) {
	g.Match([]string{"GET", "POST"}, "/v_search_pwd", self.SearchPassword)
}
func (self PasswordController) SearchPassword(c echo.Context) error {
	path := RequestUrl(c)
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "password_search.html", echo.Map{
			"Path":     path,
			"Password": "",
			"Mac":      "",
		})
	} else if c.Request().Method == "POST" {
		mac := c.FormValue("mac")
		password := genPassword(mac)
		path := RequestUrl(c)
		return c.Render(http.StatusOK, "password_search.html", echo.Map{
			"Path":     path,
			"Mac":      mac,
			"Password": password,
		})
	}
	return c.Redirect(http.StatusFound, "/v3/project/v_search_pwd")
}

const passwordMap = string("1234567890abcdefghijklmopqrstuvwxyz!@#$%&()[]{}")

func genPassword(macaddr string) string {
	mac := strings.ToLower(macaddr)
	//	fmt.Println("Lower mac:", mac)
	length := len(passwordMap)
	h := md5.New()
	io.WriteString(h, mac)
	sum := h.Sum(nil)
	md5sum := hex.EncodeToString(sum)
	//	fmt.Println("md5:", md5sum)
	var password []byte
	//j := 0
	for i := 1; i <= 10; i++ {
		password = append(password, passwordMap[(int(md5sum[i])*i)%length])
		//j++
	}
	//fmt.Println("Password:", string(password))
	return string(password)
}
