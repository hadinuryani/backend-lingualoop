package file

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, f *File) error
	FindByID(ctx context.Context, id string) (*File, error)
	FindByStoragePath(ctx context.Context, storagePath string) (*File, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, f *File) error {
	query := `
		INSERT INTO files (id, resource_type, storage_path, original_name, mime_type, size_bytes, uploaded_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		f.ID, f.ResourceType, f.StoragePath, f.OriginalName, f.MimeType, f.SizeBytes, f.UploadedBy, f.CreatedAt, f.UpdatedAt)
	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*File, error) {
	query := `
		SELECT id, resource_type, storage_path, original_name, mime_type, size_bytes, uploaded_by, created_at, updated_at
		FROM files
		WHERE id = ? AND deleted_at IS NULL
		LIMIT 1
	`
	var f File
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID, &f.ResourceType, &f.StoragePath, &f.OriginalName, &f.MimeType, &f.SizeBytes, &f.UploadedBy, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &f, nil
}

func (r *repository) FindByStoragePath(ctx context.Context, storagePath string) (*File, error) {
	query := `
		SELECT id, resource_type, storage_path, original_name, mime_type, size_bytes, uploaded_by, created_at, updated_at
		FROM files
		WHERE storage_path = ? AND deleted_at IS NULL
		LIMIT 1
	`
	var f File
	err := r.db.QueryRowContext(ctx, query, storagePath).Scan(
		&f.ID, &f.ResourceType, &f.StoragePath, &f.OriginalName, &f.MimeType, &f.SizeBytes, &f.UploadedBy, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &f, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `UPDATE files SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return ErrFileNotFound
	}
	
	return nil
}
