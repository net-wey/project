package Entity

import (
	"time"
)

type Developer struct {
	ID         uint
	Firstname  string
	LastName   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time
}