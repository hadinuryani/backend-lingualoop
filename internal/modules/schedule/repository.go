package schedule

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	GetConfig(ctx context.Context) (*ScheduleConfig, error)
	SaveConfig(ctx context.Context, config *ScheduleConfig) error

	FindAll(ctx context.Context) ([]*Schedule, error)
	FindByClass(ctx context.Context, classID, academicYearID string) ([]*Schedule, error)
	FindByID(ctx context.Context, id string) (*Schedule, error)
	FindClassClash(ctx context.Context, classID, academicYearID, day string, period int, exceptID string) (bool, error)
	FindTeacherClash(ctx context.Context, teacherID, academicYearID, day string, period int, exceptID string) (bool, error)
	
	Create(ctx context.Context, s *Schedule) error
	Update(ctx context.Context, s *Schedule) error
	Delete(ctx context.Context, id string) error
	DeleteByClass(ctx context.Context, classID, academicYearID string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetConfig(ctx context.Context) (*ScheduleConfig, error) {
	query := `
		SELECT id, periods_per_day, period_duration, start_time, break_after_periods, break_durations, active_days, created_at, updated_at
		FROM schedule_configs
		WHERE id = 1
	`
	var c ScheduleConfig
	err := r.db.QueryRowContext(ctx, query).Scan(
		&c.ID, &c.PeriodsPerDay, &c.PeriodDuration, &c.StartTime,
		&c.BreakAfterPeriods, &c.BreakDurations, &c.ActiveDays,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &c, nil
}

func (r *repository) SaveConfig(ctx context.Context, c *ScheduleConfig) error {
	query := `
		INSERT INTO schedule_configs (id, periods_per_day, period_duration, start_time, break_after_periods, break_durations, active_days)
		VALUES (1, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			periods_per_day = VALUES(periods_per_day),
			period_duration = VALUES(period_duration),
			start_time = VALUES(start_time),
			break_after_periods = VALUES(break_after_periods),
			break_durations = VALUES(break_durations),
			active_days = VALUES(active_days)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.PeriodsPerDay, c.PeriodDuration, c.StartTime,
		c.BreakAfterPeriods, c.BreakDurations, c.ActiveDays,
	)
	return err
}

type Scanner interface {
	Scan(dest ...any) error
}

func scanSchedule(scanner Scanner) (*Schedule, error) {
	var s Schedule
	err := scanner.Scan(
		&s.ID, &s.AcademicYearID, &s.ClassID, &s.SubjectID, &s.TeacherID,
		&s.Day, &s.Period, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Schedule, error) {
	query := `
		SELECT id, academic_year_id, class_id, subject_id, teacher_id, day, period, created_at, updated_at
		FROM schedules
		WHERE deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Schedule
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) FindByClass(ctx context.Context, classID, academicYearID string) ([]*Schedule, error) {
	query := `
		SELECT id, academic_year_id, class_id, subject_id, teacher_id, day, period, created_at, updated_at
		FROM schedules
		WHERE class_id = ? AND academic_year_id = ? AND deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, classID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*Schedule
	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Schedule, error) {
	query := `
		SELECT id, academic_year_id, class_id, subject_id, teacher_id, day, period, created_at, updated_at
		FROM schedules
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`
	s, err := scanSchedule(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *repository) FindClassClash(ctx context.Context, classID, academicYearID, day string, period int, exceptID string) (bool, error) {
	query := `
		SELECT 1
		FROM schedules
		WHERE class_id = ? AND academic_year_id = ? AND day = ? AND period = ? AND deleted_at IS NULL
	`
	args := []any{classID, academicYearID, day, period}
	if exceptID != "" {
		query += " AND id != ?"
		args = append(args, exceptID)
	}
	query += " LIMIT 1"

	var exists int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) FindTeacherClash(ctx context.Context, teacherID, academicYearID, day string, period int, exceptID string) (bool, error) {
	query := `
		SELECT 1
		FROM schedules
		WHERE teacher_id = ? AND academic_year_id = ? AND day = ? AND period = ? AND deleted_at IS NULL
	`
	args := []any{teacherID, academicYearID, day, period}
	if exceptID != "" {
		query += " AND id != ?"
		args = append(args, exceptID)
	}
	query += " LIMIT 1"

	var exists int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) Create(ctx context.Context, s *Schedule) error {
	query := `
		INSERT INTO schedules (id, academic_year_id, class_id, subject_id, teacher_id, day, period)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.AcademicYearID, s.ClassID, s.SubjectID, s.TeacherID, s.Day, s.Period,
	)
	return err
}

func (r *repository) Update(ctx context.Context, s *Schedule) error {
	query := `
		UPDATE schedules
		SET subject_id = ?, teacher_id = ?, day = ?, period = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		s.SubjectID, s.TeacherID, s.Day, s.Period, s.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrScheduleNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `UPDATE schedules SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrScheduleNotFound
	}
	return nil
}

func (r *repository) DeleteByClass(ctx context.Context, classID, academicYearID string) error {
	query := `UPDATE schedules SET deleted_at = NOW() WHERE class_id = ? AND academic_year_id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, classID, academicYearID)
	return err
}
