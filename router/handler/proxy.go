package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/s3-oauth2-proxy/usecase"
)

type ProxyHandler struct {
	proxyUsecase usecase.ProxyUsecase
}

func NewProxyHandler(proxyUsecase usecase.ProxyUsecase) *ProxyHandler {
	return &ProxyHandler{proxyUsecase: proxyUsecase}
}

func (h *ProxyHandler) GetObject(c echo.Context) error {
	if !isAuthenticated(c) {
		return c.Redirect(http.StatusFound, c.Echo().Reverse("auth.login"))
	}

	key := c.Request().URL.Path
	if strings.HasSuffix(key, "/") {
		key = path.Join(key, "index.html")
	}

	obj, err := h.proxyUsecase.GetObject(c.Request().Context(), key)
	if err != nil {
		if errors.Is(err, usecase.ErrObjectNotFound) {
			slog.ErrorContext(c.Request().Context(), "object not found", "key", key)
			return c.String(http.StatusNotFound, "Object not found")
		}
		slog.ErrorContext(c.Request().Context(), "failed to get object", "key", key, "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	c.Response().Header().Set("Cache-Control", "private")
	c.Response().Header().Set("Pragma", "no-cache")

	return c.Stream(http.StatusOK, obj.ContentType, obj.Body)
}
