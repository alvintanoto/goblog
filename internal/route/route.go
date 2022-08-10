package route

import (
	h "alvintanoto.id/blog/internal/route/handler"
	m "alvintanoto.id/blog/internal/route/middleware"
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo/v4"
)

func Init(port string) {
	logger := log.Get()

	middleware := &m.Middleware{
		Logger: logger,
	}

	handler := &h.Handler{
		Logger: logger,
	}

	e := echo.New()
	e.Use(middleware.LogRequest)

	e.GET("/healthz", handler.Healthz)

	logger.InfoLog.Fatal(e.Start(port))
}
