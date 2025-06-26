package handler

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
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
	sess, _ := session.Get("session", c)
	if sess.Values["authenticated"] != true {
		return c.Redirect(http.StatusFound, "/auth/login")
	}

	key := c.Request().URL.Path
	if key == "/" {
		key = "index.html"
	}

	obj, err := h.proxyUsecase.GetObject(c.Request().Context(), key)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Stream(http.StatusOK, obj.ContentType, obj.Body)
}
