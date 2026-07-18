package schedule

import "time"

type ScheduleConfigRequest struct {
	PeriodsPerDay     int      `json:"periods_per_day" binding:"required"`
	PeriodDuration    int      `json:"period_duration" binding:"required"`
	StartTime         string   `json:"start_time" binding:"required"`
	BreakAfterPeriods []int    `json:"break_after_periods"`
	BreakDurations    []int    `json:"break_durations"`
	ActiveDays        []string `json:"active_days"`
}

type ScheduleConfigResponse struct {
	PeriodsPerDay     int      `json:"periods_per_day"`
	PeriodDuration    int      `json:"period_duration"`
	StartTime         string   `json:"start_time"`
	BreakAfterPeriods []int    `json:"break_after_periods"`
	BreakDurations    []int    `json:"break_durations"`
	ActiveDays        []string `json:"active_days"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ScheduleRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	ClassID        string `json:"class_id" binding:"required"`
	SubjectID      string `json:"subject_id" binding:"required"`
	TeacherID      string `json:"teacher_id" binding:"required"`
	Day            string `json:"day" binding:"required"`
	Period         int    `json:"period" binding:"required"`
}

type ScheduleResponse struct {
	ID             string    `json:"id"`
	AcademicYearID string    `json:"academic_year_id"`
	ClassID        string    `json:"class_id"`
	SubjectID      string    `json:"subject_id"`
	TeacherID      string    `json:"teacher_id"`
	Day            string    `json:"day"`
	Period         int       `json:"period"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
