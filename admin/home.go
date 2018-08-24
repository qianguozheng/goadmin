package admin

import (
	"strings"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	. "github.com/qianguozheng/goadmin/http"
	mw "github.com/qianguozheng/goadmin/http/middleware"
	"github.com/qianguozheng/goadmin/logic"
	"github.com/qianguozheng/goadmin/model"
)

type HomeController struct{}

func (self HomeController) RegisterRoute(g *echo.Group) {
	g.Any("/index", self.Login)
	g.Any("/", self.Login)
	g.Any("/login.html", self.Login)
	g.POST("/login", self.LoginPost)
	g.GET("/reset.html", self.Reset) //TODO:implement email send and reset
	g.GET("/v3/core/logout", self.Logout, mw.NeedLogin)
	g.GET("/v3/core/personal/v_edit", self.PersonalEdit, mw.NeedLogin)
	g.GET("/v3/core/personal/v_check_pwd", self.CheckPass, mw.NeedLogin)
	g.GET("/v3/core/personal/v_check_new_pwd", self.CheckNewPass, mw.NeedLogin)
	g.POST("/v3/core/personal/o_update", self.Update, mw.NeedLogin)

	g.GET("/home.html", self.Home, mw.NeedLogin)
	g.GET("/v3/core/index", self.ProjectIndex, mw.NeedLogin)
	g.GET("/v3/project/nav", self.ProjectIndex, mw.NeedLogin)
}

func (home HomeController) CheckPass(c echo.Context) error {
	old := c.QueryParam("origPwd")
	user := c.Get("user")
	if old == user.(model.User).Password {
		return c.String(http.StatusOK, "true") //true or false
	} else {
		return c.String(http.StatusOK, "false") //true or false
	}
}

func (home HomeController) CheckNewPass(c echo.Context) error {
	pass := c.QueryParam("password")
	user := c.Get("user")
	if pass == user.(model.User).Password {
		return c.String(http.StatusOK, "false") //true or false
	} else {
		return c.String(http.StatusOK, "true") //true or false
	}
}

// origPwd: ubuntu123
// password: weeds123
// email: richard.qian@magicwifi.com.cn
// realName: 钱国正
// mobile: 13538273761
// id: 267
// username: richard.qian
// rank: 1
// status: true
// imageUri: /mnt/resource/admin/core/20180824/jpg/09/7b/097b2399c5164fb68b1db2660ecf0eb8.jpg
// modelName: edit

func (home HomeController) Update(c echo.Context) error {
	pass := c.FormValue("password")
	email := c.FormValue("email")
	phone := c.FormValue("mobile")
	username := c.FormValue("username")
	nick := c.FormValue("realName")
	user := model.User{
		Name:     username,
		Password: pass,
		Phone:    phone,
		Email:    email,
		NickName: nick,
	}
	logic.DefaultUser.UpdateUser(user)
	return c.Redirect(http.StatusFound, "/v3/core/personal/v_edit?modelName=edit")
}

func (home HomeController) Login(c echo.Context) error {
	path := RequestUrl(c)
	user := c.Get("user")
	return c.Render(http.StatusOK, "login.html", echo.Map{
		"Path": path,
		"User": user,
	})
}

func (home HomeController) LoginPost(c echo.Context) error {

	// userName := c.FormValue("username")
	// password := c.FormValue("password")
	//
	// cookie := new(http.Cookie)
	// cookie.Name = "sessionId"
	// cookie.Value = base64.StdEncoding.EncodeToString([]byte(userName + password + "honeybot"))
	// cookie.Expires = time.Now().Add(24 * time.Hour)
	// //Store cookie value to db
	// model.SetCookie(userName, cookie.Value)
	// c.SetCookie(cookie)

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || c.Request().Method != "POST" {
		return c.Redirect(http.StatusFound, "/")
	}

	// fmt.Println("post, check user pass:", username, password)
	//Check username and password in db
	if !model.GetUser(username, password) {
		// fmt.Println("no user exist")
		return c.Redirect(http.StatusFound, "/")
	}

	// fmt.Println("set cookie")

	SetLoginCookie(c, username)
	return c.Redirect(http.StatusFound, "/home.html")
}

func (home HomeController) Reset(c echo.Context) error {
	path := RequestUrl(c)
	user := c.Get("user")
	return c.Render(http.StatusOK, "reset.html", echo.Map{
		"Path": path,
		"User": user,
	})
}

func (home HomeController) Home(c echo.Context) error {

	// fmt.Println("USER:", c.Get("user"))
	path := RequestUrl(c)
	user := c.Get("user")
	return c.Render(http.StatusOK, "home.html", echo.Map{
		"Path": path,
		"User": user,
	})
}

func (home HomeController) ProjectIndex(c echo.Context) error {
	path := RequestUrl(c)
	user := c.Get("user")
	return c.Render(http.StatusOK, "project_index.html", echo.Map{
		"Path": path,
		"User": user,
	})
}

func (host HomeController) Logout(c echo.Context) error {
	// cookie, err := c.Cookie("sessionId")
	// if err != nil {
	// 	return c.Redirect(http.StatusFound, "/login.html")
	// }
	// model.Logout(cookie.Value)
	// return c.Redirect(http.StatusFound, "/login.html")
	// 删除cookie信息
	session := GetCookieSession(c)
	session.Options = &sessions.Options{Path: "/", MaxAge: -1}
	session.Save(Request(c), ResponseWriter(c))
	// 重定向得到原页面
	return c.Redirect(http.StatusSeeOther, c.Request().Referer())
}

func (HomeController) PersonalEdit(c echo.Context) error {
	user := c.Get("user")
	if strings.Compare(c.QueryParam("modelName"), "edit") == 0 {
		return c.Render(http.StatusOK, "personal_edit.html", echo.Map{
			"User": user,
		})
	}
	return c.Render(http.StatusOK, "personal_editextra.html", echo.Map{
		"User": user,
	})
}
