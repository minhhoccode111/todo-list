package entity

import "time"

type PriorityLevel string

const (
	PriorityLevelLow  PriorityLevel = "low"
	PriorityLevelMed  PriorityLevel = "med"
	PriorityLevelHigh PriorityLevel = "high"
)

type Todo struct {
	ID          int32
	UserID      int32
	Title       string
	Description string
	Completed   bool
	Priority    PriorityLevel
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
