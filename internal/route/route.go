package route

import (
	"net/http"

	model "alvintanoto.id/blog/internal/model/response"
	h "alvintanoto.id/blog/internal/route/handler"
	m "alvintanoto.id/blog/internal/route/middleware"
	"alvintanoto.id/blog/pkg/helper"
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

	echo.NotFoundHandler = func(c echo.Context) error {
		resp := model.BaseResponse{
			Code:    http.StatusNotFound,
			Message: "Not Found",
		}

		b, err := helper.PrettyJson(resp)
		if err != nil {
			logger.ErrorLog.Println(err)
		}

		logger.InfoLog.Printf("Response JSON: \n%s", string(b))

		return c.JSON(http.StatusNotFound, resp)
	}

	e := echo.New()
	e.Use(middleware.LogRequest)

	e.GET("/healthz", handler.Healthz)

	e.Static("/static", "./ui/static")

	logger.InfoLog.Fatal(e.Start(port))
}
