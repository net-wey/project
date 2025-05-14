package Entity

import "time"

type Report struct {
	ID          uint
	DeveloperID uint
	CreatedAt   time.Time
	Developer   *Developer
	Tasks       []Task
