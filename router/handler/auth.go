package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/s3-oauth2-proxy/usecase"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Login(c echo.Context) error {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	state := base64.StdEncoding.EncodeToString(b)

	_, err = rand.Read(b)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	nonce := base64.StdEncoding.EncodeToString(b)

	sess, _ := session.Get("session", c)
	sess.Values["state"] = state
	sess.Values["nonce"] = nonce
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusFound, h.authUsecase.GetAuthorizationURL(state, nonce))
}

func (h *AuthHandler) Callback(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if c.QueryParam("state") != sess.Values["state"] {
		return c.String(http.StatusBadRequest, "invalid state")
	}

	expectedNonce, ok := sess.Values["nonce"].(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "nonce not found in session")
	}

	_, err := h.authUsecase.Login(c.Request().Context(), c.QueryParam("code"), expectedNonce)
	if err != nil {
		return c.String(http.StatusForbidden, err.Error())
	}

	sess.Values["authenticated"] = true
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusFound, "/")
}
