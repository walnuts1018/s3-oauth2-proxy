package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Liveness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (h *HealthHandler) Readiness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
