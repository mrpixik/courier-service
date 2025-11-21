package domain

import "time"

// Courier основная сущность БД
type Courier struct {
	Id        int
	Name      string
	Phone     string
	Status    string // available || busy || paused
	CreatedAt time.Time
	UpdatedAt time.Time
}
