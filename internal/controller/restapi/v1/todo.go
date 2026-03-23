package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/middleware"
	"github.com/minhhoccode111/todo-list/internal/controller/restapi/v1/request"
	"github.com/minhhoccode111/todo-list/internal/entity"
	"github.com/minhhoccode111/todo-list/pkg/validatorx"
)

// @Summary     Create Todo
// @Description Create a Todo item with title and description
// @ID          create-todo
// @Tags        todo
// @Accept      json
// @Produce     json
// @Param       request body request.CreateTodo true "comment"
// @Success     200 {object} entity.Todo
// @Failure     400 {object} response.Message
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /todos [post]
func (r *V1) createTodo(c *gin.Context) {
	var body request.CreateTodo

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - createTodo - c.ShouldBindJSON")

		messageResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - createTodo - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		messageResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

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

	if body.Priority == nil {
		med := entity.PriorityLevelMed
		body.Priority = &med
	}

	t, err := r.to.CreateTodo(c, &entity.Todo{
		UserID:      userID,
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUnauthorized):
			messageResponse(c, http.StatusUnauthorized, "Unauthorized")
		default:
			r.l.Error(err, "restapi - v1 - createTodo - r.u.CreateTodo")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.JSON(http.StatusCreated, t)
}

// @Summary     Get Todos
// @Description Get paginated list of Todo items
// @ID          get-todos
// @Tags        todo
// @Produce     json
// @Param       page query int false "Page number"
// @Param       limit query int false "Items per page"
// @Success     200 {object} entity.Todos
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /todos [get]
func (r *V1) getTodos(c *gin.Context) {
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

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	result, err := r.to.GetTodos(c, userID, int32(page), int32(limit)) //nolint:gosec // intended
	if err != nil {
		r.l.Error(err, "restapi - v1 - getTodos - r.to.GetTodos")
		messageResponse(c, http.StatusInternalServerError, "internal server error")

		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     Update Todo
// @Description Update an existing Todo item
// @ID          update-todo
// @Tags        todo
// @Accept      json
// @Produce     json
// @Param       id path int true "Todo ID"
// @Param       request body request.UpdateTodo true "comment"
// @Success     200 {object} entity.Todo
// @Failure     400 {object} response.Message
// @Failure     401 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /todos/{id} [put]
func (r *V1) updateTodo(c *gin.Context) { //nolint:funlen // identical pattern to createTodo
	var body request.UpdateTodo
	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTodo - c.ShouldBindJSON")

		messageResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - updateTodo - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		messageResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))

		return
	}

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

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid todo id")

		return
	}

	if body.Priority == nil {
		med := entity.PriorityLevelMed
		body.Priority = &med
	}

	t, err := r.to.UpdateTodo(c, &entity.Todo{
		ID:          int32(todoID),
		UserID:      userID,
		Title:       body.Title,
		Description: body.Description,
		Completed:   body.Completed,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
	})
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrForbidden):
			messageResponse(c, http.StatusForbidden, "Forbidden")
		default:
			r.l.Error(err, "restapi - v1 - updateTodo - r.u.UpdateTodo")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.JSON(http.StatusOK, t)
}

// @Summary     Delete Todo
// @Description Delete a Todo item
// @ID          delete-todo
// @Tags        todo
// @Param       id path int true "Todo ID"
// @Success     204
// @Failure     401 {object} response.Message
// @Failure     403 {object} response.Message
// @Failure     500 {object} response.Message
// @Router      /todos/{id} [delete]
func (r *V1) deleteTodo(c *gin.Context) {
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

	todoID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		messageResponse(c, http.StatusBadRequest, "invalid todo id")

		return
	}

	err = r.to.DeleteTodo(c, int32(todoID), userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrForbidden):
			messageResponse(c, http.StatusForbidden, "Forbidden")
		default:
			r.l.Error(err, "restapi - v1 - deleteTodo - r.u.DeleteTodo")
			messageResponse(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	c.Status(http.StatusNoContent)
}
