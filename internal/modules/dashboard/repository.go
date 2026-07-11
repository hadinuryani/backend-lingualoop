package dashboard

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetTotalTeachers(ctx context.Context) (int, error)
	GetTotalStudents(ctx context.Context) (int, error)
	GetTotalClasses(ctx context.Context) (int, error)
	GetTotalMajors(ctx context.Context) (int, error)
	GetGenderDemographics(ctx context.Context) ([]GenderStat, error)
	GetClassLevelDistribution(ctx context.Context) ([]LevelStat, error)
	GetRecentRegistrations(ctx context.Context) ([]RecentRegistration, error)
}

const (
	countTeachers = "SELECT COUNT(*) FROM teachers WHERE deleted_at IS NULL"
	countStudents = "SELECT COUNT(*) FROM students WHERE deleted_at IS NULL"
	countClasses  = "SELECT COUNT(*) FROM classes WHERE deleted_at IS NULL"
	countMajors   = "SELECT COUNT(*) FROM majors WHERE deleted_at IS NULL"

	genderDemographics = `
		SELECT gender, COUNT(*) 
		FROM students 
		WHERE deleted_at IS NULL 
		GROUP BY gender
	`

	classLevelDistribution = `
		SELECT level_id, COUNT(*) 
		FROM classes 
		WHERE deleted_at IS NULL 
		GROUP BY level_id
	`

	recentRegistrations = `
		(SELECT id, full_name, 'Teacher' as role, created_at FROM teachers WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT 5)
		UNION ALL
		(SELECT id, full_name, 'Student' as role, created_at FROM students WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT 5)
		ORDER BY created_at DESC
		LIMIT 5
	`
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) count(ctx context.Context, query string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *repository) GetTotalTeachers(ctx context.Context) (int, error) {
	return r.count(ctx, countTeachers)
}

func (r *repository) GetTotalStudents(ctx context.Context) (int, error) {
	return r.count(ctx, countStudents)
}

func (r *repository) GetTotalClasses(ctx context.Context) (int, error) {
	return r.count(ctx, countClasses)
}

func (r *repository) GetTotalMajors(ctx context.Context) (int, error) {
	return r.count(ctx, countMajors)
}

func (r *repository) GetGenderDemographics(ctx context.Context) ([]GenderStat, error) {
	rows, err := r.db.QueryContext(ctx, genderDemographics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []GenderStat
	for rows.Next() {
		var stat GenderStat
		var gender sql.NullString
		if err := rows.Scan(&gender, &stat.Count); err != nil {
			return nil, err
		}
		if gender.Valid && gender.String != "" {
			stat.Gender = gender.String
		} else {
			stat.Gender = "Tidak Diketahui"
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *repository) GetClassLevelDistribution(ctx context.Context) ([]LevelStat, error) {
	rows, err := r.db.QueryContext(ctx, classLevelDistribution)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []LevelStat
	for rows.Next() {
		var stat LevelStat
		var level sql.NullString
		if err := rows.Scan(&level, &stat.Count); err != nil {
			return nil, err
		}
		if level.Valid && level.String != "" {
			stat.Level = level.String
		} else {
			stat.Level = "Tidak Diketahui"
		}
		stats = append(stats, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *repository) GetRecentRegistrations(ctx context.Context) ([]RecentRegistration, error) {
	rows, err := r.db.QueryContext(ctx, recentRegistrations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regs []RecentRegistration
	for rows.Next() {
		var reg RecentRegistration
		var createdAt sql.NullTime
		if err := rows.Scan(&reg.ID, &reg.FullName, &reg.Role, &createdAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			reg.CreatedAt = createdAt.Time
		}
		regs = append(regs, reg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return regs, nil
}
