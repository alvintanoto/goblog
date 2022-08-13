package route

import (
	"fmt"
	"net/http"

	"alvintanoto.id/blog/internal/database"
	"alvintanoto.id/blog/internal/database/connection"
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	Logger *log.Log
}

func (m *Middleware) LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		m.Logger.InfoLog.Printf("%s - %s %s %s", c.Request().RemoteAddr, c.Request().Proto, c.Request().Method, c.Request().RequestURI)
		return next(c)
	}
}

func (m *Middleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		exist := sess.Values["userID"] != nil
		fmt.Println("Middleware: start", sess.Values["userID"])

		if !exist {
			return next(c)
		}

		user, err := new(database.UserDB).Get(sess.Values["userID"].(int))
		if err != nil {
			if err == connection.ErrRecordNotFound {
				delete(sess.Values, "userID")
				sess.Save(c.Request(), c.Response().Writer)
				return next(c)
			}

			delete(sess.Values, "userID")
			sess.Save(c.Request(), c.Response().Writer)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "server error"})
		}

		fmt.Println("middleware", user)
		c.Set("user", user)
		return next(c)
	}
}
