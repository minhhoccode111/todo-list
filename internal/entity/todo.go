package entity

import "time"

type PriorityLevel string

const (
	PriorityLevelLow  PriorityLevel = "low"
	PriorityLevelMed  PriorityLevel = "med"
	PriorityLevelHigh PriorityLevel = "high"
)

type Todo struct {
	ID          int32          `json:"id"          validate:"required"`
	UserID      int32          `json:"-"`
	Title       string         `json:"title"       validate:"required"`
	Description string         `json:"description" validate:"required"`
	Completed   bool           `json:"completed"   validate:"required"`
	Priority    *PriorityLevel `json:"priority"    validate:"required"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"  validate:"required"`
	UpdatedAt   time.Time      `json:"updated_at"  validate:"required"`
	DeletedAt   *time.Time     `json:"-"`
}

type Todos struct {
	Data  []Todo `json:"data"  validate:"required"`
	Page  int32  `json:"page"  validate:"required"`
	Limit int32  `json:"limit" validate:"required"`
	Total int32  `json:"total" validate:"required"`
}
