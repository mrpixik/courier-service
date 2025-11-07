package domain

import "time"

// Courier основная сущность БД
type Courier struct {
	Id        int
	Name      string
	Phone     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
