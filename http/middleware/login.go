package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	. "github.com/qianguozheng/goadmin/http"
	"github.com/qianguozheng/goadmin/logic"
	"github.com/qianguozheng/goadmin/model"
)

func NeedLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// cookie, err := c.Cookie("sessionId")
		//
		// if err != nil {
		// 	if strings.Contains(err.Error(), "named cookie not present") {
		// 		return c.Redirect(http.StatusFound, "/login.html")
		// 		//return c.String(http.StatusUnauthorized, "You don't have the right cookie")
		// 	}
		// 	fmt.Println("cookie not present")
		// 	return err
		// }

		session := GetCookieSession(c)
		username, ok := session.Values["username"]
		if ok {
			//getCurrentUser(username)
			user := logic.DefaultUser.FindCurrentUser(username.(string))
			if user.(model.User).Name != username {
				return c.Redirect(http.StatusFound, "/login.html")
			}
			//Later show in page
			c.Set("user", user)
		} else {
			return c.Redirect(http.StatusFound, "/login.html")
		}

		if err := next(c); err != nil {
			return err
		}

		// if true == model.GetCookie(cookie.Value) {
		// 	return next(c)
		// }
		//
		return c.Redirect(http.StatusFound, "/login.html")
	}
}
