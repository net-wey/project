package Entity

import (
	"time"

	"github.com/google/uuid"
)

type Developer struct {
	ID         uuid.UUID
	Firstname  string
	LastName   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time
}
