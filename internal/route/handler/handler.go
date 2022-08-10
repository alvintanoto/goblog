package route

import (
	"net/http"
	"time"

	model "alvintanoto.id/blog/internal/model/response"
	"alvintanoto.id/blog/pkg/helper"
	"alvintanoto.id/blog/pkg/log"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Logger *log.Log
}

func (h *Handler) Healthz(c echo.Context) error {
	resp := &model.HealthzResponse{
		BaseResponse: model.BaseResponse{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Time: time.Now(),
	}

	b, err := helper.PrettyJson(resp)
	if err != nil {
		h.Logger.ErrorLog.Println(err)
	}

	h.Logger.InfoLog.Printf("Response JSON: \n%s", string(b))
	return c.JSON(http.StatusOK, resp)
}
