package auth

import (
	"database/sql"
	"errors"
)

// Repository interface untuk abstraction (berguna untuk testing/mocking)
type Repository interface {
	FindByEmail(email string) (*User, error)
}

// repository mengimplementasikan Repository interface
type repository struct {
	db *sql.DB
}

// NewRepository membuat instance baru dari auth repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

// FindByEmail mencari user berdasarkan alamat email
func (r *repository) FindByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, role, avatar_url, is_active
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var user User
	var avatarURL sql.NullString

	err := r.db.QueryRow(query, email).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User tidak ditemukan
		}
		return nil, err // Error lain dari database
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}

	return &user, nil
}
