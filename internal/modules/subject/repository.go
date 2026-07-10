package subject

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Subject, error)
	FindByID(ctx context.Context, id string) (*Subject, error)
	GetIDByCode(ctx context.Context, code string) (string, error)
	Create(ctx context.Context, subject *Subject) error
	Update(ctx context.Context, subject *Subject) error
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

const selectSubject = `
SELECT 
	id, 
	code, 
	name, 
	description, 
	major_id, 
	level_id, 
	created_at, 
	updated_at
FROM subjects
`

const selectSubjectIDByCode = `SELECT id FROM subjects WHERE code = ? AND deleted_at IS NULL LIMIT 1`

const insertSubject = `INSERT INTO subjects (id, code, name, description, major_id, level_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
const updateSubject = `UPDATE subjects SET code = ?, name = ?, description = ?, major_id = ?, level_id = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteSubject = `UPDATE subjects SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSubject(scanner rowScanner) (*Subject, error) {
	var s Subject
	var desc, majorID, levelID sql.NullString

	err := scanner.Scan(
		&s.ID,
		&s.Code,
		&s.Name,
		&desc,
		&majorID,
		&levelID,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if desc.Valid {
		s.Description = &desc.String
	}
	if majorID.Valid {
		s.MajorID = &majorID.String
	}
	if levelID.Valid {
		s.LevelID = &levelID.String
	}

	return &s, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Subject, error) {
	query := selectSubject + ` WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []*Subject
	for rows.Next() {
		s, err := scanSubject(rows)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Subject, error) {
	query := selectSubject + ` WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	s, err := scanSubject(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubjectNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *repository) GetIDByCode(ctx context.Context, code string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, selectSubjectIDByCode, code).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *repository) Create(ctx context.Context, subject *Subject) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, insertSubject, subject.ID, subject.Code, subject.Name, subject.Description, subject.MajorID, subject.LevelID, subject.CreatedAt, subject.UpdatedAt)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert subject")
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, subject *Subject) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateSubject, subject.Code, subject.Name, subject.Description, subject.MajorID, subject.LevelID, subject.UpdatedAt, subject.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrSubjectNotFound
	}

	return tx.Commit()
}

func (r *repository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, softDeleteSubject, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrSubjectNotFound
	}

	return tx.Commit()
}
