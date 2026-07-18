package student

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Student, error)
	FindByID(ctx context.Context, id string) (*Student, error)
	FindByNIS(ctx context.Context, nis string) (*Student, error)
	GetIDByNIS(ctx context.Context, nis string) (string, error)
	Create(ctx context.Context, student *Student, user *User) error
	Update(ctx context.Context, student *Student, user *User) error
	UpdateStatus(ctx context.Context, studentID string, userID string, status string) error
	Delete(ctx context.Context, studentID string, userID string) error
	FindAllByAcademicYear(ctx context.Context, academicYearID string) ([]*Student, error)
	InsertStudentClassesBatchTx(ctx context.Context, tx *sql.Tx, studentClasses []*StudentClass) error
	UpdateLevelAndStatusBatchTx(ctx context.Context, tx *sql.Tx, students []*Student) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

const selectStudent = `
	SELECT 
		s.id, 
		s.user_id, 
		u.username,
		u.email,
		s.nis, 
		s.full_name, 
		s.gender, 
		s.birth_place, 
		s.birth_date, 
		s.phone, 
		s.address_region, 
		s.address_detail, 
		s.photo, 
		s.major_id,
		s.class_level,
		s.status, 
		s.created_at, 
		s.updated_at
	FROM students s
	JOIN users u ON s.user_id = u.id
	`

const insertUser = `INSERT INTO users (id, email, username, password_hash, full_name, role, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateUserFullName = `UPDATE users SET full_name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const updateUserStatus = `UPDATE users SET is_active = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteUser = `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const insertStudent = `INSERT INTO students (id, user_id, nis, full_name, gender, birth_place, birth_date, phone, address_region, address_detail, photo, major_id, class_level, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateStudent = `UPDATE students SET nis = ?, full_name = ?, gender = ?, birth_place = ?, birth_date = ?, phone = ?, address_region = ?, address_detail = ?, photo = ?, major_id = ?, class_level = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const updateStudentStatus = `UPDATE students SET status = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteStudent = `UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const selectStudentIDByNIS = `SELECT id FROM students WHERE nis = ? AND deleted_at IS NULL LIMIT 1`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanStudent(scanner rowScanner) (*Student, error) {
	var s Student
	var birthPlace, phone, addrRegion, addrDetail, photo, majorID, classLevel sql.NullString
	var birthDate sql.NullTime

	err := scanner.Scan(
		&s.ID,
		&s.UserID,
		&s.Username,
		&s.Email,
		&s.NIS,
		&s.FullName,
		&s.Gender,
		&birthPlace,
		&birthDate,
		&phone,
		&addrRegion,
		&addrDetail,
		&photo,
		&majorID,
		&classLevel,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if birthPlace.Valid {
		s.BirthPlace = &birthPlace.String
	}
	if birthDate.Valid {
		s.BirthDate = &birthDate.Time
	}
	if phone.Valid {
		s.Phone = &phone.String
	}
	if addrRegion.Valid {
		s.AddressRegion = &addrRegion.String
	}
	if addrDetail.Valid {
		s.AddressDetail = &addrDetail.String
	}
	if photo.Valid {
		s.Photo = &photo.String
	}
	if majorID.Valid {
		s.MajorID = &majorID.String
	}
	if classLevel.Valid {
		s.ClassLevel = &classLevel.String
	}

	return &s, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Student, error) {
	query := selectStudent + ` WHERE s.deleted_at IS NULL AND u.deleted_at IS NULL ORDER BY s.created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Student, error) {
	query := selectStudent + ` WHERE s.id = ? AND s.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1`
	s, err := scanStudent(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *repository) FindByNIS(ctx context.Context, nis string) (*Student, error) {
	query := selectStudent + ` WHERE s.nis = ? AND s.deleted_at IS NULL AND u.deleted_at IS NULL LIMIT 1`
	s, err := scanStudent(r.db.QueryRowContext(ctx, query, nis))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *repository) GetIDByNIS(ctx context.Context, nis string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, selectStudentIDByNIS, nis).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *repository) Create(ctx context.Context, student *Student, user *User) error {
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

	result, err = tx.ExecContext(ctx, insertStudent, student.ID, student.UserID, student.NIS, student.FullName, student.Gender, student.BirthPlace, student.BirthDate, student.Phone, student.AddressRegion, student.AddressDetail, student.Photo, student.MajorID, student.ClassLevel, student.Status, student.CreatedAt, student.UpdatedAt)
	if err != nil {
		return err
	}
	affected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert student")
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, student *Student, user *User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateStudent, student.NIS, student.FullName, student.Gender, student.BirthPlace, student.BirthDate, student.Phone, student.AddressRegion, student.AddressDetail, student.Photo, student.MajorID, student.ClassLevel, student.UpdatedAt, student.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrStudentNotFound
	}

	result, err = tx.ExecContext(ctx, updateUserFullName, user.FullName, user.UpdatedAt, user.ID)
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

func (r *repository) UpdateStatus(ctx context.Context, studentID string, userID string, status string) error {
	var isActive bool
	switch status {
	case StudentActive:
		isActive = true
	case StudentGraduated, StudentTransfer, StudentInactive:
		isActive = false
	default:
		return ErrInvalidStatus
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateStudentStatus, status, studentID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrStudentNotFound
	}

	_, err = tx.ExecContext(ctx, updateUserStatus, isActive, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Delete(ctx context.Context, studentID string, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, softDeleteStudent, studentID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrStudentNotFound
	}

	_, err = tx.ExecContext(ctx, softDeleteUser, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) FindAllByAcademicYear(ctx context.Context, academicYearID string) ([]*Student, error) {
	query := `
		SELECT 
			s.id, 
			s.user_id, 
			u.username,
			u.email,
			s.nis, 
			s.full_name, 
			s.gender, 
			s.birth_place, 
			s.birth_date, 
			s.phone, 
			s.address_region, 
			s.address_detail, 
			s.photo, 
			s.major_id,
			s.class_level,
			s.status, 
			s.created_at, 
			s.updated_at,
			sc.class_id
		FROM students s
		JOIN users u ON s.user_id = u.id
		JOIN student_classes sc ON s.id = sc.student_id
		WHERE sc.academic_year_id = ? AND sc.is_active = TRUE AND s.deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*Student
	for rows.Next() {
		var s Student
		var birthPlace, phone, addressRegion, addressDetail, photo, majorID, classLevel, currentClassID sql.NullString
		var birthDate sql.NullTime

		err := rows.Scan(
			&s.ID, &s.UserID, &s.Username, &s.Email, &s.NIS, &s.FullName, &s.Gender,
			&birthPlace, &birthDate, &phone, &addressRegion, &addressDetail, &photo,
			&majorID, &classLevel, &s.Status, &s.CreatedAt, &s.UpdatedAt, &currentClassID,
		)
		if err != nil {
			return nil, err
		}

		if birthPlace.Valid {
			s.BirthPlace = &birthPlace.String
		}
		if birthDate.Valid {
			s.BirthDate = &birthDate.Time
		}
		if phone.Valid {
			s.Phone = &phone.String
		}
		if addressRegion.Valid {
			s.AddressRegion = &addressRegion.String
		}
		if addressDetail.Valid {
			s.AddressDetail = &addressDetail.String
		}
		if photo.Valid {
			s.Photo = &photo.String
		}
		if majorID.Valid {
			s.MajorID = &majorID.String
		}
		if classLevel.Valid {
			s.ClassLevel = &classLevel.String
		}
		if currentClassID.Valid {
			s.CurrentClassID = &currentClassID.String
		}

		students = append(students, &s)
	}
	return students, nil
}

func (r *repository) InsertStudentClassesBatchTx(ctx context.Context, tx *sql.Tx, studentClasses []*StudentClass) error {
	if len(studentClasses) == 0 {
		return nil
	}
	query := `INSERT INTO student_classes (id, student_id, class_id, academic_year_id, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, sc := range studentClasses {
		_, err := stmt.ExecContext(ctx, sc.ID, sc.StudentID, sc.ClassID, sc.AcademicYearID, sc.IsActive, sc.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) UpdateLevelAndStatusBatchTx(ctx context.Context, tx *sql.Tx, students []*Student) error {
	if len(students) == 0 {
		return nil
	}
	query := `UPDATE students SET class_level = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range students {
		_, err := stmt.ExecContext(ctx, s.ClassLevel, s.Status, s.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
