package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/request"
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
func (r *V1) login(c *gin.Context) {}
