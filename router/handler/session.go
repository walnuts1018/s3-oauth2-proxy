package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	sessionNameDefault = "session"
	sessionKeyState    = "state"
	sessionKeyNonce    = "nonce"
	sessionKeyReturnTo = "return_to"
)

func useAuthSession(c echo.Context, fn func(values map[any]any) error) error {
	sess, err := session.Get(sessionNameAuth, c)
	if err != nil {
		return err
	}

	if err := fn(sess.Values); err != nil {
		return err
	}
	sess.Options.SameSite = http.SameSiteLaxMode
	sess.Options = &sessions.Options{
		Path:     "/auth/",
		MaxAge:   int((1 * time.Hour).Seconds()), // 1 hour
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to save session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	return nil
}

const (
	sessionNameAuth         = "auth"
	sessionKeyAuthenticated = "authenticated"
)

func useDefaultSession(c echo.Context, fn func(values map[any]any) error) error {
	sess, err := session.Get(sessionNameDefault, c)
	if err != nil {
		return err
	}

	if err := fn(sess.Values); err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((3 * 24 * time.Hour).Seconds()), // 3 days
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return sess.Save(c.Request(), c.Response())
}
