package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/minhhoccode111/todo-list/config"
	"github.com/minhhoccode111/todo-list/internal/usecase"
	"github.com/minhhoccode111/todo-list/pkg/logger"
)

// V1 -.
type V1 struct {
	l   logger.Interface
	v   *validator.Validate
	cfg *config.Config
	tr  usecase.Translation
	u   usecase.User
	to  usecase.Todo
}
