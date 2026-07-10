package auth

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectUser = `
SELECT 
	id, 
	email, 
	username, 
	password_hash, 
	full_name, 
	role, 
	avatar_url, 
	is_active
FROM users
`

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner scanner) (*User, error) {
	var user User
	var avatarURL sql.NullString

	err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&avatarURL,
		&user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return &user, nil
}

// FindByEmail mencari user berdasarkan alamat email
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := selectUser + `WHERE email = ? AND deleted_at IS NULL LIMIT 1`

	user, err := scanUser(r.db.QueryRowContext(ctx, query, email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User tidak ditemukan
		}
		return nil, err // Error lain dari database
	}

	return user, nil
}
