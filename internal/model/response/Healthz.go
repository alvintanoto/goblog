package model

import "time"

type HealthzResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
