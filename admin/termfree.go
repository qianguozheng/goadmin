package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"github.com/labstack/echo"
)

func TerminalListPost(c echo.Context) error {
	cookie, err := c.Cookie("_cookie_page_size")
	if err != nil {
		//return c.String(http.StatusNotFound, "No cookie_page_size found")
		fmt.Println("no cookie_page_size")
	}
	var pageSize, pageNo int
	pageSize, err = strconv.Atoi(cookie.Value)
	if err != nil {
		pageSize = 20
	}
	pageNo, err = strconv.Atoi(c.FormValue("pageNo"))
	if err != nil {
		pageNo = 1
	}

	//fmt.Println("cookie_page_size:", cookie.Value)
	fmt.Println("pageNo:", pageNo)

	devNum := model.GetTotalTermFreeNum()
	pageNum := devNum / pageSize
	if devNum%pageSize > 0 {
		pageNum = pageNum + 1
	}
	fmt.Println("pageNum:", pageNum)

	//term := model.GetAllTermFree()
	terms := model.ListPageNoDeviceTermFree(pageNo, pageSize)

	fmt.Println("TermFree:", terms)

	path := RequestUrl(c)
	return c.Render(http.StatusOK, "terminal_list.html", echo.Map{
		"Path":        path,
		"TermFree":    terms,
		"PageNo":      pageNo,
		"PageNum":     pageNum,
		"TotalDevice": devNum, //Total free terminal device
		"PageSize":    pageSize,
	})
}
func TerminalList(c echo.Context) error {

	cookie, err := c.Cookie("_cookie_page_size")
	if err != nil {
		//return c.String(http.StatusNotFound, "No cookie_page_size found")
		fmt.Println("no cookie_page size")
	}
	var pageSize int
	if cookie != nil {
		pageSize, err = strconv.Atoi(cookie.Value)
	}
	if err != nil {
		pageSize = 20
	}

	devNum := model.GetTotalTermFreeNum()
	pageNum := devNum / pageSize
	if devNum%pageSize > 0 {
		pageNum = pageNum + 1
	}
	fmt.Println("pageNum:", pageNum)

	//term := model.GetAllTermFree()
	terms := model.ListPageNoDeviceTermFree(1, pageSize)

	fmt.Println("TermFree:", terms)

	path := RequestUrl(c)
	return c.Render(http.StatusOK, "terminal_list.html", echo.Map{
		"Path":        path,
		"TermFree":    terms,
		"PageNo":      1,
		"PageNum":     pageNum,
		"TotalDevice": devNum, //Total free terminal device
		"PageSize":    pageSize,
	})
}

func TerminalAdd(c echo.Context) error {

	path := RequestUrl(c)
	return c.Render(http.StatusOK, "terminal_add.html", echo.Map{
		"Path": path,
	})
}
func TerminalAddPost(c echo.Context) error {
	mac := c.FormValue("terminalMac")
	comment := c.FormValue("remark")

	fmt.Println("mac:", mac)
	fmt.Println("comment:", comment)

	model.AddTermFree(mac, comment)

	return c.Redirect(http.StatusFound, "/v3/project/terminal_free/v_list")
}

func TerminalDelete(c echo.Context) error {
	mac := c.QueryParam("mac")
	model.DeleteTermFreeByMac(mac)
	return c.Redirect(http.StatusFound, "/v3/project/terminal_free/v_list")
}

func TerminalCheck(c echo.Context) error {
	mac := c.QueryParam("terminalMac")

	result := model.CheckTermFreeByMac(mac)

	var s string
	if result {
		s = "false"
	} else {
		s = "true"
	}
	return c.String(http.StatusOK, s)
}
