package route

import (
	"alvintanoto.id/blog/pkg/log"
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
