package class

import (
	"context"
	"database/sql"
	"errors"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*Class, error)
	FindAllByAcademicYear(ctx context.Context, academicYearID string) ([]*Class, error)
	FindByID(ctx context.Context, id string) (*Class, error)
	GetIDByNameAndYear(ctx context.Context, className, academicYearID string) (string, error)
	GetLevelByGrade(ctx context.Context, grade int) (id string, name string, err error)
	GetMajorCodeByID(ctx context.Context, majorID string) (string, error)
	FindNamesByLevelMajorYear(ctx context.Context, levelID, majorID, academicYearID string) ([]string, error)
	Create(ctx context.Context, class *Class) error
	CreateBatch(ctx context.Context, classes []*Class) error
	CreateBatchTx(ctx context.Context, tx *sql.Tx, classes []*Class) error
	Update(ctx context.Context, class *Class) error
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

const selectClass = `
SELECT 
	id, 
	academic_year_id, 
	major_id, 
	level_id, 
	grade_level,
	class_number,
	class_name, 
	classroom, 
	capacity, 
	homeroom_teacher_id, 
	is_active, 
	created_at, 
	updated_at
FROM classes
`

const insertClass = `INSERT INTO classes (id, academic_year_id, major_id, level_id, grade_level, class_number, class_name, classroom, capacity, homeroom_teacher_id, is_active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
const updateClass = `UPDATE classes SET capacity = ?, homeroom_teacher_id = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
const softDeleteClass = `UPDATE classes SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

const selectClassIDByNameYear = `SELECT id FROM classes WHERE class_name = ? AND academic_year_id = ? AND deleted_at IS NULL LIMIT 1`
const selectClassNamesByLevelMajorYear = `SELECT class_name FROM classes WHERE level_id = ? AND major_id = ? AND academic_year_id = ? AND deleted_at IS NULL`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanClass(scanner rowScanner) (*Class, error) {
	var c Class
	var majorID, classroom, homeroomID sql.NullString

	err := scanner.Scan(
		&c.ID,
		&c.AcademicYearID,
		&majorID,
		&c.LevelID,
		&c.GradeLevel,
		&c.ClassNumber,
		&c.ClassName,
		&classroom,
		&c.Capacity,
		&homeroomID,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if majorID.Valid {
		c.MajorID = &majorID.String
	}
	if classroom.Valid {
		c.Classroom = &classroom.String
	}
	if homeroomID.Valid {
		c.HomeroomTeacherID = &homeroomID.String
	}

	return &c, nil
}

func (r *repository) FindAll(ctx context.Context) ([]*Class, error) {
	query := selectClass + ` WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*Class
	for rows.Next() {
		c, err := scanClass(rows)
		if err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Class, error) {
	query := selectClass + ` WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	c, err := scanClass(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrClassNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *repository) GetIDByNameAndYear(ctx context.Context, className, academicYearID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, selectClassIDByNameYear, className, academicYearID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *repository) FindNamesByLevelMajorYear(ctx context.Context, levelID, majorID, academicYearID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, selectClassNamesByLevelMajorYear, levelID, majorID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func (r *repository) Create(ctx context.Context, class *Class) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, insertClass, class.ID, class.AcademicYearID, class.MajorID, class.LevelID, class.GradeLevel, class.ClassNumber, class.ClassName, class.Classroom, class.Capacity, class.HomeroomTeacherID, class.IsActive, class.CreatedAt, class.UpdatedAt)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("failed to insert class")
	}

	return tx.Commit()
}

func (r *repository) CreateBatch(ctx context.Context, classes []*Class) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, insertClass)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, class := range classes {
		result, err := stmt.ExecContext(ctx, class.ID, class.AcademicYearID, class.MajorID, class.LevelID, class.GradeLevel, class.ClassNumber, class.ClassName, class.Classroom, class.Capacity, class.HomeroomTeacherID, class.IsActive, class.CreatedAt, class.UpdatedAt)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("failed to insert class in batch")
		}
	}

	return tx.Commit()
}

func (r *repository) Update(ctx context.Context, class *Class) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, updateClass, class.Capacity, class.HomeroomTeacherID, class.UpdatedAt, class.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrClassNotFound
	}

	return tx.Commit()
}

func (r *repository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, softDeleteClass, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrClassNotFound
	}

	return tx.Commit()
}

func (r *repository) GetLevelByGrade(ctx context.Context, grade int) (string, string, error) {
	var id, name string
	err := r.db.QueryRowContext(ctx, "SELECT id, name FROM levels WHERE grade_level = ?", grade).Scan(&id, &name)
	return id, name, err
}

func (r *repository) GetMajorCodeByID(ctx context.Context, majorID string) (string, error) {
	var code string
	err := r.db.QueryRowContext(ctx, "SELECT major_code FROM majors WHERE id = ?", majorID).Scan(&code)
	return code, err
}

func (r *repository) FindAllByAcademicYear(ctx context.Context, academicYearID string) ([]*Class, error) {
	query := selectClass + ` WHERE academic_year_id = ? AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []*Class
	for rows.Next() {
		c, err := scanClass(rows)
		if err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, nil
}

func (r *repository) CreateBatchTx(ctx context.Context, tx *sql.Tx, classes []*Class) error {
	stmt, err := tx.PrepareContext(ctx, insertClass)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, class := range classes {
		result, err := stmt.ExecContext(ctx, class.ID, class.AcademicYearID, class.MajorID, class.LevelID, class.GradeLevel, class.ClassNumber, class.ClassName, class.Classroom, class.Capacity, class.HomeroomTeacherID, class.IsActive, class.CreatedAt, class.UpdatedAt)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("failed to insert class in batch tx")
		}
	}

	return nil
}
