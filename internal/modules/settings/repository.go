package settings

import (
	"context"
	"database/sql"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Repository interface {
	GetConfig(ctx context.Context) (*SchoolSettings, error)
	UpdateConfig(ctx context.Context, settings *SchoolSettings) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func scanSchoolSettings(scanner Scanner) (*SchoolSettings, error) {
	var s SchoolSettings
	err := scanner.Scan(
		&s.ID, &s.SchoolName, &s.SchoolNPSN, &s.SchoolAddress, &s.SchoolPhone, &s.SchoolEmail,
		&s.SchoolLogoFileID, &s.SchoolLogoPath,
		&s.EducationLogoFileID, &s.EducationLogoPath,
		&s.PrincipalName, &s.PrincipalNIP,
		&s.PrincipalSignatureFileID, &s.PrincipalSignaturePath,
		&s.MaxStudentsPerClass, &s.GradingSystem, &s.PassingGrade,
		&s.AppName, &s.EnableStudentLogin, &s.EnableTeacherLogin, &s.MaintenanceMode,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) GetConfig(ctx context.Context) (*SchoolSettings, error) {
	query := `
		SELECT 
			s.id, s.school_name, s.school_npsn, s.school_address, s.school_phone, s.school_email,
			s.school_logo_file_id, f1.storage_path as school_logo_path,
			s.education_logo_file_id, f2.storage_path as education_logo_path,
			s.principal_name, s.principal_nip,
			s.principal_signature_file_id, f3.storage_path as principal_signature_path,
			s.max_students_per_class, s.grading_system, s.passing_grade,
			s.app_name, s.enable_student_login, s.enable_teacher_login, s.maintenance_mode,
			s.created_at, s.updated_at
		FROM school_settings s
		LEFT JOIN files f1 ON s.school_logo_file_id = f1.id
		LEFT JOIN files f2 ON s.education_logo_file_id = f2.id
		LEFT JOIN files f3 ON s.principal_signature_file_id = f3.id
		WHERE s.id = 1
	`
	
	s, err := scanSchoolSettings(r.db.QueryRowContext(ctx, query))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSettingsNotFound
		}
		return nil, err
	}
	
	return s, nil
}

func (r *repository) UpdateConfig(ctx context.Context, s *SchoolSettings) error {
	query := `
		UPDATE school_settings
		SET 
			school_name = ?, school_npsn = ?, school_address = ?, school_phone = ?, school_email = ?,
			school_logo_file_id = ?, education_logo_file_id = ?, 
			principal_name = ?, principal_nip = ?, principal_signature_file_id = ?,
			max_students_per_class = ?, grading_system = ?, passing_grade = ?,
			app_name = ?, enable_student_login = ?, enable_teacher_login = ?, maintenance_mode = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		s.SchoolName, s.SchoolNPSN, s.SchoolAddress, s.SchoolPhone, s.SchoolEmail,
		s.SchoolLogoFileID, s.EducationLogoFileID,
		s.PrincipalName, s.PrincipalNIP, s.PrincipalSignatureFileID,
		s.MaxStudentsPerClass, s.GradingSystem, s.PassingGrade,
		s.AppName, s.EnableStudentLogin, s.EnableTeacherLogin, s.MaintenanceMode,
	)
	
	if err != nil {
		return err
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrSettingsNotFound
	}
	
	return nil
}
