package schedule

import (
	"encoding/json"
	"time"
)

type ScheduleConfig struct {
	ID                int
	PeriodsPerDay     int
	PeriodDuration    int
	StartTime         string
	BreakAfterPeriods json.RawMessage // stored as JSON array string
	BreakDurations    json.RawMessage // stored as JSON array string
	ActiveDays        json.RawMessage // stored as JSON array string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Schedule struct {
	ID             string
	AcademicYearID string
	ClassID        string
	SubjectID      string
	TeacherID      string
	Day            string
	Period         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
