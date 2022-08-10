package model

import "time"

type HealthzResponse struct {
	BaseResponse
	Time time.Time `json:"time"`
}
