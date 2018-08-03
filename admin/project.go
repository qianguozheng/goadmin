package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"github.com/labstack/echo"
)

func HandleProjectList(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)
	prjs := model.GetProjects()

	return c.Render(http.StatusOK, "project_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func HandleProjectAdd(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)

	return c.Render(http.StatusOK, "project_add.html", echo.Map{
		"Path": path,
	})
}

func HandleProjectSave(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)

	name := c.FormValue("name")
	comment := c.FormValue("comment")

	err := model.InsertProject(name, comment)
	if err != nil {
		return c.String(http.StatusForbidden, "Already exist")
	}

	prjs := model.GetProjects()
	return c.Render(http.StatusOK, "project_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func HandleProjectEdit(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)

	ids := c.QueryParam("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		return c.String(http.StatusOK, "id not assign")
	}

	prj := model.GetProjectById(id)
	return c.Render(http.StatusOK, "project_edit.html", echo.Map{
		"Path":    path,
		"Project": prj,
	})
}

func HandleProjectDelete(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)

	if c.Request().Method == "POST" {
		ids := getIDSFromParams(c)
		for k, v := range ids {
			fmt.Println("k=", k, " v:", v)
			model.DeleteProjectById(v)
		}

	} else if c.Request().Method == "GET" {
		ids := c.QueryParam("ids")

		id, err := strconv.Atoi(ids)
		if err != nil {
			return c.String(http.StatusOK, "id not assign")
		}

		fmt.Println("id:", id)
		model.DeleteProjectById(id)
	}

	prjs := model.GetProjects()

	fmt.Println("prjs:", prjs)
	return c.Render(http.StatusOK, "project_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}
