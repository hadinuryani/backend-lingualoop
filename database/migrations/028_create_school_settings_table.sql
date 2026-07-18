CREATE TABLE school_settings (
    id INT PRIMARY KEY DEFAULT 1,
    school_name VARCHAR(255),
    school_npsn VARCHAR(50),
    school_address TEXT,
    school_phone VARCHAR(50),
    school_email VARCHAR(100),
    school_logo_file_id VARCHAR(36) NULL,
    education_logo_file_id VARCHAR(36) NULL,
    principal_name VARCHAR(100),
    principal_nip VARCHAR(50),
    principal_signature_file_id VARCHAR(36) NULL,
    max_students_per_class INT DEFAULT 36,
    grading_system VARCHAR(20) DEFAULT 'numeric',
    passing_grade INT DEFAULT 75,
    app_name VARCHAR(100) DEFAULT 'LinguaLoop',
    enable_student_login BOOLEAN DEFAULT TRUE,
    enable_teacher_login BOOLEAN DEFAULT TRUE,
    maintenance_mode BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (school_logo_file_id) REFERENCES files(id) ON DELETE SET NULL,
    FOREIGN KEY (education_logo_file_id) REFERENCES files(id) ON DELETE SET NULL,
    FOREIGN KEY (principal_signature_file_id) REFERENCES files(id) ON DELETE SET NULL
);

-- Insert the default config row
INSERT INTO school_settings (
    id, school_name, school_npsn, school_address, school_phone, school_email, 
    principal_name, principal_nip, max_students_per_class, grading_system, 
    passing_grade, app_name, enable_student_login, enable_teacher_login, maintenance_mode
) VALUES (
    1, 'SMK Negeri 1 LinguaLoop', '20100001', 'Jl. Pendidikan No. 1, Kota Bandung, Jawa Barat', 
    '(022) 1234567', 'info@smkn1lingualoop.sch.id', 'Drs. H. Mulyana, M.Pd.', '19740512 200003 1 002', 
    36, 'numeric', 75, 'LinguaLoop', TRUE, TRUE, FALSE
);
