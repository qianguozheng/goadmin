package admin

import (
	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

const HOST string = "121.176.55.180"

var PageSize int = 20

type Page struct {
	PageSize     int // One page element number
	PageNum      int // total page number
	PageNo       int // which page are you in
	TotalElement int // total element
}

func GeneratePage(c echo.Context, elemNum int) Page {
	page := Page{
		PageSize:     20,
		PageNum:      1,
		PageNo:       1,
		TotalElement: elemNum,
	}

	cookie, err := c.Cookie("_cookie_page_size")
	if err != nil {
		return page
	}

	if cookie != nil {
		page.PageSize = goutils.MustInt(cookie.Value, 20)
	}

	page.PageNo = goutils.MustInt(c.FormValue("pageNo"), 1)
	page.PageNum = elemNum / page.PageSize
	if elemNum%page.PageSize > 0 {
		page.PageNum = page.PageNum + 1
	}
	return page
}
