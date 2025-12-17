package model

import "time"

const (
	StatusAssigned   = "assigned"
	StatusUnassigned = "unassigned"
	StatusCompleted  = "completed"
)

// Delivery сущность из таблицы delivery
type Delivery struct {
	Id         int
	CourierId  int
	OrderId    string
	AssignedAt time.Time
	Deadline   time.Time
}
