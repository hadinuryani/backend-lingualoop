package teacher_portal

import (
	"context"
	"database/sql"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Repository interface {
	GetTeacherIDByUserID(ctx context.Context, userID string) (string, error)
	GetClassesByTeacherID(ctx context.Context, teacherID string) ([]TeacherClass, error)
	GetSchedulesByTeacherID(ctx context.Context, teacherID string) ([]TeacherSchedule, error)
	GetScheduleConfig(ctx context.Context) (*ScheduleConfig, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func scanTeacherClass(scanner Scanner) (TeacherClass, error) {
	var cls TeacherClass
	err := scanner.Scan(
		&cls.ID,
		&cls.Name,
		&cls.Major,
		&cls.Capacity,
		&cls.Subject,
		&cls.Room,
		&cls.Enrolled,
	)
	return cls, err
}

func scanTeacherSchedule(scanner Scanner) (TeacherSchedule, error) {
	var s TeacherSchedule
	var room sql.NullString
	
	err := scanner.Scan(
		&s.ID,
		&s.Day,
		&s.Period,
		&s.ClassName,
		&room,
		&s.Subject,
		&s.Major,
	)
	if err == nil && room.Valid {
		s.Room = room.String
	}
	return s, err
}

func (r *repository) GetTeacherIDByUserID(ctx context.Context, userID string) (string, error) {
	var teacherID string
	query := `SELECT id FROM teachers WHERE user_id = ? AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&teacherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrTeacherNotFound
		}
		return "", err
	}
	return teacherID, nil
}

func (r *repository) GetScheduleConfig(ctx context.Context) (*ScheduleConfig, error) {
	query := `SELECT periods_per_day, period_duration, start_time FROM schedule_configs WHERE id = 1`
	var c ScheduleConfig
	err := r.db.QueryRowContext(ctx, query).Scan(&c.PeriodsPerDay, &c.PeriodDuration, &c.StartTime)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetClassesByTeacherID(ctx context.Context, teacherID string) ([]TeacherClass, error) {
	query := `
		SELECT 
			c.id, 
			c.class_name, 
			COALESCE(m.name, 'Umum') as major_name,
			c.capacity, 
			s.name as subject_name, 
			tsc.classroom,
			COALESCE(sc_counts.total_enrolled, 0) as enrolled_count
		FROM teacher_subject_classes tsc
		JOIN classes c ON tsc.class_id = c.id
		JOIN subjects s ON tsc.subject_id = s.id
		JOIN academic_years ay ON tsc.academic_year_id = ay.id
		LEFT JOIN majors m ON c.major_id = m.id
		LEFT JOIN (
			SELECT class_id, COUNT(*) as total_enrolled 
			FROM student_classes 
			WHERE is_active = 1 AND deleted_at IS NULL 
			GROUP BY class_id
		) sc_counts ON sc_counts.class_id = c.id
		WHERE tsc.teacher_id = ? 
		AND ay.status = 'Aktif'
		AND ay.deleted_at IS NULL
		AND c.is_active = 1 
		AND tsc.deleted_at IS NULL
		AND c.deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []TeacherClass
	for rows.Next() {
		cls, err := scanTeacherClass(rows)
		if err != nil {
			return nil, err
		}
		classes = append(classes, cls)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (r *repository) GetSchedulesByTeacherID(ctx context.Context, teacherID string) ([]TeacherSchedule, error) {
	query := `
		SELECT 
			sch.id,
			sch.day,
			sch.period,
			c.class_name,
			tsc.classroom,
			s.name as subject_name,
			COALESCE(m.name, 'Umum') as major_name
		FROM schedules sch
		JOIN classes c ON sch.class_id = c.id
		JOIN subjects s ON sch.subject_id = s.id
		JOIN academic_years ay ON sch.academic_year_id = ay.id
		LEFT JOIN majors m ON c.major_id = m.id
		LEFT JOIN teacher_subject_classes tsc ON tsc.class_id = c.id 
			AND tsc.subject_id = s.id 
			AND tsc.teacher_id = sch.teacher_id 
			AND tsc.deleted_at IS NULL
		WHERE sch.teacher_id = ?
		AND ay.status = 'Aktif'
		AND ay.deleted_at IS NULL
		AND sch.deleted_at IS NULL
		ORDER BY FIELD(sch.day, 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'), sch.period ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var schedules []TeacherSchedule
	for rows.Next() {
		s, err := scanTeacherSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return schedules, nil
}
