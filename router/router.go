package router

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/walnuts1018/s3-oauth2-proxy/config"
	"github.com/walnuts1018/s3-oauth2-proxy/router/handler"
	"github.com/walnuts1018/s3-oauth2-proxy/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func NewRouter(cfg *config.AppConfig, authHandler *handler.AuthHandler, proxyHandler *handler.ProxyHandler, healthHandler *handler.HealthHandler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: skipper,
	}))
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(tracer.ServiceName, otelecho.WithSkipper(
		skipper,
	)))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.SessionSecret))))

	e.GET("/auth/login", authHandler.Login)
	e.GET("/auth/callback", authHandler.Callback)

	e.GET("/livez", healthHandler.Liveness)
	e.GET("/readyz", healthHandler.Readiness)

	e.Any("/*", proxyHandler.GetObject)

	return e
}

func skipper(c echo.Context) bool {
	// Skip tracing for health check endpoints
	return c.Path() == "/livez" || c.Path() == "/readyz"
}
