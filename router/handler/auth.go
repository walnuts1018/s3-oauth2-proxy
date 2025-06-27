package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/s3-oauth2-proxy/usecase"
	"github.com/walnuts1018/s3-oauth2-proxy/util/random"
)

const (
	sessionKeyState         = "state"
	sessionKeyNonce         = "nonce"
	sessionKeyAuthenticated = "authenticated"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	random      random.Random
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, random random.Random) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, random: random}
}

func (h *AuthHandler) Login(c echo.Context) error {
	state, err := h.random.SecureString(32, random.Alphanumeric)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to generate state", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	nonce, err := h.random.SecureString(32, random.Alphanumeric)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to generate nonce", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	sess, _ := session.Get("session", c)
	sess.Values[sessionKeyState] = state
	sess.Values[sessionKeyNonce] = nonce
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to save session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusFound, h.authUsecase.GetAuthorizationURL(state, nonce))
}

func (h *AuthHandler) Callback(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if c.QueryParam("state") != sess.Values[sessionKeyState] {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	expectedNonce, ok := sess.Values[sessionKeyNonce].(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "nonce not found in session")
	}

	_, err := h.authUsecase.Login(c.Request().Context(), c.QueryParam("code"), expectedNonce)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to login", "error", err)
		return c.String(http.StatusForbidden, "failed to login")
	}

	sess.Values[sessionKeyAuthenticated] = true
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to save session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusFound, "/")
}

func isAuthenticated(c echo.Context) bool {
	sess, err := session.Get("session", c)
	if err != nil {
		return false
	}
	authStatus, ok := sess.Values[sessionKeyAuthenticated].(bool)
	return ok && authStatus
}
