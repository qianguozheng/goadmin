package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"../model"
	"github.com/labstack/echo"
)

type ProjectController struct{}

//Project Management
// adminGrp.GET("/v3/project/project/v_list", admin.HandleProjectList)
// adminGrp.GET("/v3/project/project/v_add", admin.HandleProjectAdd)
// adminGrp.POST("/v3/project/project/o_save", admin.HandleProjectSave)
// adminGrp.GET("/v3/project/project/o_delete", admin.HandleProjectDelete)
// adminGrp.POST("/v3/project/project/o_delete", admin.HandleProjectDelete)
// adminGrp.GET("/v3/project/project/v_edit", admin.HandleProjectEdit)
// 注册路由
func (self ProjectController) RegisterRoute(g *echo.Group) {
	g.GET("/project/v_list", self.List)
	g.GET("/project/v_add", self.Add)
	g.POST("/project/o_save", self.Save)
	g.Match([]string{"GET", "POST"}, "/v3/project/project/o_delete", self.Delete)
	g.GET("/project/v_edit", self.Edit)
}

func (self ProjectController) List(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)
	prjs := model.GetProjects()

	return c.Render(http.StatusOK, "project_list.html", echo.Map{
		"Path":     path,
		"Projects": prjs,
	})
}

func (self ProjectController) Add(c echo.Context) error {
	path := RequestUrl(c)
	fmt.Println("path=", path)

	return c.Render(http.StatusOK, "project_add.html", echo.Map{
		"Path": path,
	})
}

func (self ProjectController) Save(c echo.Context) error {
	// path := RequestUrl(c)
	// fmt.Println("path=", path)

	name := c.FormValue("name")
	comment := c.FormValue("comment")

	err := model.InsertProject(name, comment)
	if err != nil {
		return c.String(http.StatusForbidden, "Already exist")
	}

	// prjs := model.GetProjects()
	// return c.Render(http.StatusOK, "project_list.html", echo.Map{
	// 	"Path":     path,
	// 	"Projects": prjs,
	// })
	return c.Redirect(http.StatusFound, "/v3/project/project/v_list")
}

func (self ProjectController) Edit(c echo.Context) error {
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

func (self ProjectController) Delete(c echo.Context) error {
	// path := RequestUrl(c)
	// fmt.Println("path=", path)

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

	// prjs := model.GetProjects()
	//
	// fmt.Println("prjs:", prjs)
	// return c.Render(http.StatusOK, "project_list.html", echo.Map{
	// 	"Path":     path,
	// 	"Projects": prjs,
	// })
	return c.Redirect(http.StatusFound, "/v3/project/project/v_list")
}
