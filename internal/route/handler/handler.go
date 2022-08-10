package route

import (
	"net/http"
	"time"

	model "alvintanoto.id/blog/internal/model/response"
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger *log.Log
}

func (h *Handler) Healthz(c echo.Context) error {
	resp := &model.HealthzResponse{
		Code:    http.StatusOK,
		Message: "Success",
		Time:    time.Time{},
	}
	return c.JSON(http.StatusOK, resp)
}
