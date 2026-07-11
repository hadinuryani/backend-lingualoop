package academic_year

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*AcademicYear, error)
	FindByID(ctx context.Context, id string) (*AcademicYear, error)
	GetIDByYear(ctx context.Context, year string) (string, error)
	CheckActiveYearExcept(ctx context.Context, id string) (bool, error)
	ActivateYear(ctx context.Context, id string) error
	Create(ctx context.Context, ay *AcademicYear) error
	Update(ctx context.Context, ay *AcademicYear) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectAcademicYear = `
SELECT 
	id, 
	year, 
	start_date, 
	end_date, 
	status,
	sem_ganjil_start_date,
	sem_ganjil_end_kbm,
	sem_ganjil_end_assessment,
	sem_ganjil_status,
	sem_genap_start_date,
	sem_genap_end_kbm,
	sem_genap_end_assessment,
	sem_genap_status,
	created_at, 
	updated_at
FROM academic_years
`

const insertAcademicYear = `INSERT INTO academic_years (id, year, start_date, end_date, status, sem_ganjil_start_date, sem_ganjil_end_kbm, sem_ganjil_end_assessment, sem_ganjil_status, sem_genap_start_date, sem_genap_end_kbm, sem_genap_end_assessment, sem_genap_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateAcademicYear = `UPDATE academic_years SET year = ?, start_date = ?, end_date = ?, status = ?, sem_ganjil_start_date = ?, sem_ganjil_end_kbm = ?, sem_ganjil_end_assessment = ?, sem_ganjil_status = ?, sem_genap_start_date = ?, sem_genap_end_kbm = ?, sem_genap_end_assessment = ?, sem_genap_status = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteAcademicYear = `UPDATE academic_years SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const selectAcademicYearIDByYear = `SELECT id FROM academic_years WHERE year = ? AND deleted_at IS NULL LIMIT 1`
const checkActiveExcept = `SELECT EXISTS(SELECT 1 FROM academic_years WHERE status = ? AND id != ? AND deleted_at IS NULL)`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAY(scanner rowScanner) (*AcademicYear, error) {
	var ay AcademicYear
	err := scanner.Scan(
		&ay.ID,
		&ay.Year,
		&ay.StartDate,
		&ay.EndDate,
		&ay.Status,
		&ay.SemGanjilStartDate,
		&ay.SemGanjilEndKBM,
		&ay.SemGanjilEndAssessment,
		&ay.SemGanjilStatus,
		&ay.SemGenapStartDate,
		&ay.SemGenapEndKBM,
		&ay.SemGenapEndAssessment,
		&ay.SemGenapStatus,
		&ay.CreatedAt,
		&ay.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ay, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*AcademicYear, error) {
	query := selectAcademicYear + ` WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []*AcademicYear
	for rows.Next() {
		ay, err := scanAY(rows)
		if err != nil {
			return nil, err
		}
		years = append(years, ay)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return years, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*AcademicYear, error) {
	query := selectAcademicYear + ` WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	ay, err := scanAY(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAcademicYearNotFound
		}
		return nil, err
	}
	return ay, nil
}

func (r *repository) GetIDByYear(ctx context.Context, year string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, selectAcademicYearIDByYear, year).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrAcademicYearNotFound
		}
		return "", err
	}
	return id, nil
}

func (r *repository) CheckActiveYearExcept(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, checkActiveExcept, StatusActive, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) ActivateYear(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	// Kunci baris yang aktif agar tidak ada dua proses yang mengecek secara bersamaan
	err = tx.QueryRowContext(ctx, checkActiveExcept+" FOR UPDATE", StatusActive, id).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrMultipleActiveYears
	}

	query := `UPDATE academic_years SET status = ?, sem_ganjil_status = ?, sem_genap_status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	result, err := tx.ExecContext(ctx, query, StatusActive, SemStatusActive, SemStatusNotActive, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrAcademicYearNotFound
	}

	return tx.Commit()
}

func (r *repository) Create(ctx context.Context, ay *AcademicYear) error {
	result, err := r.db.ExecContext(ctx, insertAcademicYear,
		ay.ID, ay.Year, ay.StartDate, ay.EndDate, ay.Status,
		ay.SemGanjilStartDate, ay.SemGanjilEndKBM, ay.SemGanjilEndAssessment, ay.SemGanjilStatus,
		ay.SemGenapStartDate, ay.SemGenapEndKBM, ay.SemGenapEndAssessment, ay.SemGenapStatus,
		ay.CreatedAt, ay.UpdatedAt,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert academic year")
	}

	return nil
}

func (r *repository) Update(ctx context.Context, ay *AcademicYear) error {
	result, err := r.db.ExecContext(ctx, updateAcademicYear,
		ay.Year, ay.StartDate, ay.EndDate, ay.Status,
		ay.SemGanjilStartDate, ay.SemGanjilEndKBM, ay.SemGanjilEndAssessment, ay.SemGanjilStatus,
		ay.SemGenapStartDate, ay.SemGenapEndKBM, ay.SemGenapEndAssessment, ay.SemGenapStatus,
		ay.UpdatedAt, ay.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrAcademicYearNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, softDeleteAcademicYear, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrAcademicYearNotFound
	}

	return nil
}
