package subject

import "time"

type Subject struct {
	ID          string
	Code        string
	Name        string
	Description *string
	MajorID     *string
	LevelID     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
