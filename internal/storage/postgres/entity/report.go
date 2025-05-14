package Entity

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID          uint
	DeveloperID uuid.UUID
	CreatedAt   time.Time
	Developer   *Developer
	Tasks       []Task
}
