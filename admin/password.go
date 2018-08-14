package admin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func HandleSearchPassword(c echo.Context) error {
	path := RequestUrl(c)
	return c.Render(http.StatusOK, "password_search.html", echo.Map{
		"Path":     path,
		"Password": "",
		"Mac":      "",
	})
}

const passwordMap = string("1234567890abcdefghijklmopqrstuvwxyz!@#$%&()[]{}")

func HandleSearchPasswordPost(c echo.Context) error {
	mac := c.FormValue("mac")
	mac = strings.ToLower(mac)
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
	fmt.Println("Password:", string(password))

	path := RequestUrl(c)
	return c.Render(http.StatusOK, "password_search.html", echo.Map{
		"Path":     path,
		"Mac":      mac,
		"Password": string(password),
	})
}
