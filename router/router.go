package router

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/walnuts1018/s3-oauth2-proxy/config"
	"github.com/walnuts1018/s3-oauth2-proxy/router/handler"
)

func NewRouter(cfg *config.AppConfig, authHandler *handler.AuthHandler, proxyHandler *handler.ProxyHandler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.SessionSecret))))

	e.GET("/auth/login", authHandler.Login)
	e.GET("/auth/callback", authHandler.Callback)

	e.Any("/*", proxyHandler.GetObject)

	return e
}
