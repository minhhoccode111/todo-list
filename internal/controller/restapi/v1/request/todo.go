package request

import (
	"time"

	"github.com/minhhoccode111/todo-list/internal/entity"
)

type CreateTodo struct {
	Title       string                `json:"title"       validate:"required,max=255"`
	Description string                `json:"description" validate:"required,max=10000"`
	Priority    *entity.PriorityLevel `json:"priority"    validate:"omitempty,oneof=low med high"`
	DueDate     *time.Time            `json:"due_date"    validate:"omitempty,future"`
}

// func (ct *CreateTodo) Trim() {
// 	ct.Title = strings.TrimSpace(ct.Title)
// 	ct.Description = strings.TrimSpace(ct.Description)
// }
