package v1

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
			r.l.Error(err, "restapi - v1 - login - r.u.Login")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.SetCookieData(&http.Cookie{
		Name:     cookieName,
		Value:    newRefresh,
		Path:     "/",
		Expires:  time.Now().Add(r.cfg.RefreshToken.Expiration),
		Secure:   r.cfg.RefreshToken.Secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, response.Auth{Token: token})
}
