package teacher

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Teacher, error)
	FindByID(ctx context.Context, id string) (*Teacher, error)
	FindByNIP(ctx context.Context, nip string) (*Teacher, error)
	GetIDByNIP(ctx context.Context, nip string) (string, error)
	Create(ctx context.Context, teacher *Teacher, user *User) error
	Update(ctx context.Context, teacher *Teacher, user *User) error
	UpdateStatus(ctx context.Context, teacherID string, userID string, status string) error
	Delete(ctx context.Context, teacherID string, userID string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectTeacher = `
SELECT 
	t.id, 
	t.user_id, 
	u.username,
	u.email,
	t.nip, 
	t.full_name, 
	t.gender, 
	t.birth_place, 
	t.birth_date, 
	t.phone, 
	t.address_region, 
	t.address_detail, 
	t.photo, 
	t.status, 
	t.created_at, 
	t.updated_at
FROM teachers t
JOIN users u ON t.user_id = u.id
`

const insertUser = `INSERT INTO users (id, email, username, password_hash, full_name, role, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateUser = `UPDATE users SET full_name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const updateUserStatus = `UPDATE users SET is_active = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteUser = `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const insertTeacher = `INSERT INTO teachers (id, user_id, nip, full_name, gender, birth_place, birth_date, phone, address_region, address_detail, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateTeacher = `UPDATE teachers SET nip = ?, full_name = ?, gender = ?, birth_place = ?, birth_date = ?, phone = ?, address_region = ?, address_detail = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const updateTeacherStatus = `UPDATE teachers SET status = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteTeacher = `UPDATE teachers SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const selectTeacherIDByNIP = `SELECT id FROM teachers WHERE nip = ? AND deleted_at IS NULL LIMIT 1`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTeacher(scanner rowScanner) (*Teacher, error) {
	var t Teacher
	var birthPlace, phone, addrRegion, addrDetail, photo sql.NullString
	var birthDate sql.NullTime

	err := scanner.Scan(
		&t.ID,
		&t.UserID,
		&t.Username,
		&t.Email,
		&t.NIP,
		&t.FullName,
		&t.Gender,
		&birthPlace,
		&birthDate,
		&phone,
		&addrRegion,
		&addrDetail,
		&photo,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if birthPlace.Valid {
		t.BirthPlace = &birthPlace.String
	}
	if birthDate.Valid {
		t.BirthDate = &birthDate.Time
	}
	if phone.Valid {
		t.Phone = &phone.String
	}
	if addrRegion.Valid {
		t.AddressRegion = &addrRegion.String
	}
	if addrDetail.Valid {
		t.AddressDetail = &addrDetail.String
	}
	if photo.Valid {
		t.Photo = &photo.String
	}

	return &t, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Teacher, error) {
	query := selectTeacher + ` WHERE t.deleted_at IS NULL AND u.deleted_at IS NULL ORDER BY t.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []*Teacher
	for rows.Next() {
		t, err := scanTeacher(rows)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teachers, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Teacher, error) {
	query := selectTeacher + ` WHERE t.id = ? AND t.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1`
	t, err := scanTeacher(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeacherNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *repository) FindByNIP(ctx context.Context, nip string) (*Teacher, error) {
	query := selectTeacher + ` WHERE t.nip = ? AND t.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1`
	t, err := scanTeacher(r.db.QueryRowContext(ctx, query, nip))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeacherNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *repository) GetIDByNIP(ctx context.Context, nip string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, selectTeacherIDByNIP, nip).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *repository) Create(ctx context.Context, teacher *Teacher, user *User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, insertUser, user.ID, user.Email, user.Username, user.PasswordHash, user.FullName, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert user")
	}

	result, err = tx.ExecContext(ctx, insertTeacher, teacher.ID, teacher.UserID, teacher.NIP, teacher.FullName, teacher.Gender, teacher.BirthPlace, teacher.BirthDate, teacher.Phone, teacher.AddressRegion, teacher.AddressDetail, teacher.Status, teacher.CreatedAt, teacher.UpdatedAt)
	if err != nil {
		return err
	}
	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert teacher")
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, teacher *Teacher, user *User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateTeacher, teacher.NIP, teacher.FullName, teacher.Gender, teacher.BirthPlace, teacher.BirthDate, teacher.Phone, teacher.AddressRegion, teacher.AddressDetail, teacher.UpdatedAt, teacher.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTeacherNotFound
	}

	// Update user as well (currently updates full_name, but flexible for more fields)
	result, err = tx.ExecContext(ctx, updateUser, user.FullName, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to update user")
	}

	return tx.Commit()
}

func (r *repository) UpdateStatus(ctx context.Context, teacherID string, userID string, status string) error {
	var isActive bool
	switch status {
	case TeacherActive:
		isActive = true
	case TeacherInactive:
		isActive = false
	default:
		return ErrInvalidStatus
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateTeacherStatus, status, teacherID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTeacherNotFound
	}

	_, err = tx.ExecContext(ctx, updateUserStatus, isActive, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Delete(ctx context.Context, teacherID string, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, softDeleteTeacher, teacherID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrTeacherNotFound
	}

	_, err = tx.ExecContext(ctx, softDeleteUser, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
