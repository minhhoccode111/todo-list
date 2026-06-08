package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/request"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/response"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/validatorx"
)

const cookieName = "todo-list-refresh-token"

// @Summary     Register
// @Description Register a user with name, email and password
// @ID          register
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body request.Register true "comment"
// @Success     201 {object} response.Auth
// @Failure     400 {object} response.Message
// @Failure     409 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /register [post]
func (r *V1) register(c *gin.Context) {
	var body request.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - register - c.ShouldBindJSON")

		messageResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - register - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		messageResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

	token, refresh, err := r.u.Register(
		c.Request.Context(),
		r.cfg,
		&entity.User{
			Email:        body.Email,
			Name:         body.Name,
			PasswordHash: body.Password,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateEntry):
			messageResponse(c, http.StatusConflict, "email already existed")
		default:
			r.l.Error(err, "restapi - v1 - register - r.u.Register")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    refresh,
		Path:     "/",
		Expires:  time.Now().Add(r.cfg.RefreshToken.Expiration),
		Secure:   r.cfg.RefreshToken.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, response.Auth{Token: token})
}

// @Summary     Login
// @Description Login a user with email and password
// @ID          login
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body request.Login true "comment"
// @Success     200 {object} response.Auth
// @Failure     400 {object} response.Message
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /login [post]
func (r *V1) login(c *gin.Context) {
	var body request.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - login - c.ShouldBindJSON")

		messageResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - login - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		messageResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

	token, refresh, err := r.u.Login(
		c.Request.Context(),
		r.cfg,
		&entity.User{
			Email:        body.Email,
			PasswordHash: body.Password,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRows):
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")
		case errors.Is(err, entity.ErrUnauthorized):
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")
		default:
			r.l.Error(err, "restapi - v1 - login - r.u.Login")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    refresh,
		Path:     "/",
		Expires:  time.Now().Add(r.cfg.RefreshToken.Expiration),
		Secure:   r.cfg.RefreshToken.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, response.Auth{Token: token})
}

// @Summary     Refresh
// @Description Refresh access token using refresh token from cookie
// @ID          refresh
// @Tags        Auth
// @Produce     json
// @Success     200 {object} response.Auth
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /refresh [post]
func (r *V1) refresh(c *gin.Context) {
	refresh, err := c.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	token, newRefresh, err := r.u.Refresh(
		c.Request.Context(),
		r.cfg,
		refresh,
	)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRows):
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")
		default:
			r.l.Error(err, "restapi - v1 - refresh - r.u.Refresh")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.SetCookieData(&http.Cookie{
		Name: cookieName, Value: newRefresh, Path: "/",
		Expires: time.Now().Add(r.cfg.RefreshToken.Expiration), Secure: r.cfg.RefreshToken.Secure,
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, response.Auth{Token: token})
}

// @Summary     Logout
// @Description Logout current session
// @ID          logout
// @Tags        Auth
// @Produce     json
// @Security    BearerAuth
// @Success     204 "No Content"
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /logout [post]
func (r *V1) logout(c *gin.Context) {
	userIDRaw, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	userID, ok := userIDRaw.(int32)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	refreshToken, err := c.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")

			return
		}

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	err = r.u.SelfLogout(c.Request.Context(), userID, refreshToken)
	if err != nil {
		r.l.Error(err, "restapi - v1 - logout - r.u.SelfLogout")
		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.SetCookieData(&http.Cookie{
		Name: cookieName, Value: "", Path: "/",
		Expires: time.Unix(0, 0), Secure: r.cfg.RefreshToken.Secure,
		HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})

	c.Status(http.StatusNoContent)
}

// @Summary     Logout All
// @Description Logout from all devices by deleting all refresh tokens for the user
// @ID          logoutAll
// @Tags        Auth
// @Produce     json
// @Security    BearerAuth
// @Success     204 "No Content"
// @Failure     500 {object} response.Message
// @Router      /logout/all [post]
func (r *V1) logoutAll(c *gin.Context) {
	userIDRaw, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	userID, ok := userIDRaw.(int32)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	err := r.u.LogoutAll(c.Request.Context(), userID)
	if err != nil {
		r.l.Error(err, "restapi - v1 - logoutAll - r.u.LogoutAll")

		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   r.cfg.RefreshToken.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.Status(http.StatusNoContent)
}

// @Summary     List Sessions
// @Description List all active sessions for the current user
// @ID          listSessions
// @Tags        Auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} response.Session
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /sessions [get]
func (r *V1) listSessions(c *gin.Context) {
	userIDRaw, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	userID, ok := userIDRaw.(int32)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	var refresh string

	refreshCookie, err := c.Cookie(cookieName)
	if err == nil {
		refresh = refreshCookie
	}

	sessions, err := r.u.ListSessions(c.Request.Context(), userID, refresh)
	if err != nil {
		r.l.Error(err, "restapi - v1 - listSessions - r.u.ListSessions")
		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	res := make([]response.Session, len(sessions))
	for i, s := range sessions {
		res[i] = response.Session{
			ID:        s.ID,
			Device:    s.DeviceInfo,
			CreatedAt: s.CreatedAt,
			ExpiresAt: s.ExpiresAt,
			IsCurrent: s.IsCurrent,
		}
	}

	c.JSON(http.StatusOK, res)
}

// @Summary     Delete Session
// @Description Logout a specific session by ID (remote logout)
// @ID          deleteSession
// @Tags        Auth
// @Produce     json
// @Param       id path int true "Session ID"
// @Security    BearerAuth
// @Success     204 "No Content"
// @Failure     401 {object} response.Message
// @Failure     404 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /sessions/{id} [delete]
func (r *V1) deleteSession(c *gin.Context) {
	userIDRaw, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	userID, ok := userIDRaw.(int32)
	if !ok {
		messageResponse(c, http.StatusUnauthorized, "Unauthorized")

		return
	}

	sessionIDStr := c.Param("id")

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 32)
	if err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid session id")

		return
	}

	err = r.u.DeleteSession(c.Request.Context(), userID, int32(sessionID))
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			messageResponse(c, http.StatusNotFound, "session not found")
		} else {
			r.l.Error(err, "restapi - v1 - deleteSession - r.u.DeleteSession")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.Status(http.StatusNoContent)
}
