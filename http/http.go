package http

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/qianguozheng/goadmin/config"
)

var Store = sessions.NewCookieStore([]byte(config.ConfigFile.MustValue("global", "cookie_secret", "defaultcookie_secret@rich")))

func SetLoginCookie(ctx echo.Context, username string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.FormValue("remember_me") != "1" {
		// 浏览器关闭，cookie删除，否则保存30天(github.com/gorilla/sessions 包的默认值)
		session.Options = &sessions.Options{
			Path:     "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = username
	req := Request(ctx)
	resp := ResponseWriter(ctx)
	session.Save(req, resp)
}

func SetCookie(ctx echo.Context, key, value string) {
	Store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	session.Values[key] = value
	req := Request(ctx)
	resp := ResponseWriter(ctx)
	session.Save(req, resp)
}

func GetFromCookie(ctx echo.Context, key string) string {
	session := GetCookieSession(ctx)
	val, ok := session.Values[key]
	if ok {
		return val.(string)
	}
	return ""
}

// 必须是 http.Request
func GetCookieSession(ctx echo.Context) *sessions.Session {
	session, _ := Store.Get(Request(ctx), "user")
	return session
}

func Request(ctx echo.Context) *http.Request {
	return ctx.Request() //.(*http.Request).Request
}

func ResponseWriter(ctx echo.Context) http.ResponseWriter {
	return ctx.Response() //.(*http.Response).ResponseWriter
}
