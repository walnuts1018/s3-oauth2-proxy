package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/s3-oauth2-proxy/usecase"
	"github.com/walnuts1018/s3-oauth2-proxy/util/random"
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

	if err := useAuthSession(c, func(values map[any]any) error {
		values[sessionKeyState] = state
		values[sessionKeyNonce] = nonce
		return nil
	}); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to set state and nonce in session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusFound, h.authUsecase.GetAuthorizationURL(state, nonce))
}

func (h *AuthHandler) Callback(c echo.Context) error {
	var expectedState string
	var expectedNonce string
	var returnTo string
	if err := useAuthSession(c, func(values map[any]any) error {
		expectedState, _ = values[sessionKeyState].(string)
		expectedNonce, _ = values[sessionKeyNonce].(string)
		returnTo, _ = values[sessionKeyReturnTo].(string)
		delete(values, sessionKeyState)
		delete(values, sessionKeyNonce)
		delete(values, sessionKeyReturnTo)
		return nil
	}); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to get state and nonce from session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	if c.QueryParam("state") != expectedState {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	if returnTo == "" {
		returnTo = "/"
	}

	if !isAcceptableRedirectURL(returnTo) {
		returnTo = "/"
	}

	if _, err := h.authUsecase.Login(c.Request().Context(), c.QueryParam("code"), expectedNonce); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to login", "error", err)
		return c.String(http.StatusForbidden, "failed to login")
	}

	if err := useDefaultSession(c, func(values map[any]any) error {
		values[sessionKeyAuthenticated] = true
		return nil
	}); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to set authenticated in session", "error", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	return c.Redirect(http.StatusFound, returnTo)
}

func isAuthenticated(c echo.Context) bool {
	var authStatus bool
	if err := useDefaultSession(c, func(values map[any]any) error {
		authStatus, _ = values[sessionKeyAuthenticated].(bool)
		return nil
	}); err != nil {
		slog.ErrorContext(c.Request().Context(), "failed to initialize authenticated in session", "error", err)
		return false
	}

	return authStatus
}

func isAcceptableRedirectURL(url string) bool {
	return len(url) > 0 && url[0] == '/'
}
