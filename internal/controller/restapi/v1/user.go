package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/request"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/response"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/validatorx"
)

// @Summary     Register
// @Description Register a user with name, email and password
// @ID          register
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body request.Register true "comment"
// @Success     201 {object} response.Auth
// @Failure     400 {object} response.Error
// @Failure     409 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /register [post]
func (r *V1) register(c *gin.Context) {
	var body request.Register

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - register - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - register - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

	token, err := r.u.Register(
		c.Request.Context(),
		&entity.User{
			Email:        body.Email,
			Name:         body.Name,
			PasswordHash: body.Password,
		},
		&r.cfg.JWT,
	)
	if err != nil {
		r.l.Error(err, "restapi - v1 - register - r.u.Register")

		switch {
		case errors.Is(err, entity.ErrDuplicateEntry):
			errorResponse(c, http.StatusConflict, "email already existed")
		default:
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

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
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /login [post]
func (r *V1) login(c *gin.Context) {
	var body request.Login

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - login - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - login - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

	token, err := r.u.Login(
		c.Request.Context(),
		&entity.User{
			Email:        body.Email,
			PasswordHash: body.Password,
		},
		&r.cfg.JWT,
	)
	if err != nil {
		r.l.Error(err, "restapi - v1 - login - r.u.Login")

		switch {
		case errors.Is(err, entity.ErrNoRows):
			errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		case errors.Is(err, entity.ErrUnauthorized):
			errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		default:
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.JSON(http.StatusOK, response.Auth{Token: token})
}
