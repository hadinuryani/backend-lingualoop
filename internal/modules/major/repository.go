package major

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Major, error)
	FindByID(ctx context.Context, id string) (*Major, error)
	FindByCode(ctx context.Context, code string) (*Major, error)
	FindByName(ctx context.Context, name string) (*Major, error)
	Create(ctx context.Context, major *Major) error
	Update(ctx context.Context, major *Major) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

const selectMajor = `
SELECT 
	id, 
	code, 
	name, 
	description, 
	created_at, 
	updated_at 
FROM majors
`

type Scanner interface {
	Scan(dest ...any) error
}

func scanMajor(scanner Scanner) (*Major, error) {
	var m Major
	var desc sql.NullString
	if err := scanner.Scan(&m.ID, &m.Code, &m.Name, &desc, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return nil, err
	}
	if desc.Valid {
		m.Description = desc.String
	}
	return &m, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Major, error) {
	query := selectMajor + `WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var majors []*Major
	for rows.Next() {
		m, err := scanMajor(rows)
		if err != nil {
			return nil, err
		}
		majors = append(majors, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return majors, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Major, error) {
	query := selectMajor + `WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	m, err := scanMajor(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	return m, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Major, error) {
	query := `SELECT id FROM majors WHERE code = ? AND deleted_at IS NULL LIMIT 1`
	var m Major
	err := r.db.QueryRowContext(ctx, query, code).Scan(&m.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *repository) FindByName(ctx context.Context, name string) (*Major, error) {
	query := `SELECT id FROM majors WHERE name = ? AND deleted_at IS NULL LIMIT 1`
	var m Major
	err := r.db.QueryRowContext(ctx, query, name).Scan(&m.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMajorNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *repository) Create(ctx context.Context, major *Major) error {
	query := `INSERT INTO majors (id, code, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, major.ID, major.Code, major.Name, major.Description, major.CreatedAt, major.UpdatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, major *Major) error {
	query := `UPDATE majors SET code = ?, name = ?, description = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, major.Code, major.Name, major.Description, major.UpdatedAt, major.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMajorNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `UPDATE majors SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMajorNotFound
	}
	return nil
}
