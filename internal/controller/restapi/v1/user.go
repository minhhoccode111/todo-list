package v1

import (
	"github.com/gin-gonic/gin"
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
func (r *V1) register(c *gin.Context) {}

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
