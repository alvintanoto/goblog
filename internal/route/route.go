package route

import (
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo/v4"
)

func Init(port string) {
	logger := log.Get()
	e := echo.New()

	logger.InfoLog.Fatal(e.Start(port))
}
