package domain

import "time"

// Delivery сущность из таблицы delivery
type Delivery struct {
	Id         int
	CourierId  int
	OrderId    string
	AssignedAt time.Time
	Deadline   time.Time
}
