package entity

import "time"

type PriorityLevel string

const (
	PriorityLevelLow  PriorityLevel = "low"
	PriorityLevelMed  PriorityLevel = "med"
	PriorityLevelHigh PriorityLevel = "high"
)

type Todo struct {
	ID          int32         `json:"id"`
	UserID      int32         `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Completed   bool          `json:"completed"`
	Priority    PriorityLevel `json:"priority"`
	DueDate     time.Time     `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   time.Time     `json:"deleted_at"`
}
