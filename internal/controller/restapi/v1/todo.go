package v1

import "github.com/gin-gonic/gin"

// @Summary     Create Todo
// @Description Create a Todo item with title and description
// @ID          create-todo
// @Tags        todo
// @Accept      json
// @Produce     json
// @Param       request body request.CreateTodo true "comment"
// @Success     200 {object} entity.Todo
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /login [post]
func (r *V1) createTodo(c *gin.Context) {}
