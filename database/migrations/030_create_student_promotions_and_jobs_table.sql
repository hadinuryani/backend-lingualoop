-- ============================================================
-- Migration 030: Tabel Audit Kenaikan Kelas & Jobs
-- ============================================================

-- Tambah promotion_completed_at ke academic_years
ALTER TABLE academic_years ADD COLUMN promotion_completed_at DATETIME NULL DEFAULT NULL;

-- Tabel promotion_jobs (Background job tracking)
CREATE TABLE IF NOT EXISTS promotion_jobs (
    id CHAR(36) PRIMARY KEY,
    academic_year_id CHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL, -- PENDING, RUNNING, DONE, FAILED
    total_students INT DEFAULT 0,
    success_students INT DEFAULT 0,
    failed_students INT DEFAULT 0,
    executed_by CHAR(36),
    started_at DATETIME,
    finished_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_promotion_jobs_year (academic_year_id),
    INDEX idx_promotion_jobs_status (status),
    
    CONSTRAINT fk_promotion_jobs_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabel student_promotions (Riwayat mutasi siswa)
CREATE TABLE IF NOT EXISTS student_promotions (
    id CHAR(36) PRIMARY KEY,
    student_id CHAR(36) NOT NULL,
    from_class_id CHAR(36),
    to_class_id CHAR(36) NULL,
    from_academic_year_id CHAR(36) NOT NULL,
    to_academic_year_id CHAR(36) NULL,
    status VARCHAR(20) NOT NULL, -- PROMOTED, RETAINED, GRADUATED, FAILED
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_student_promotions_student (student_id),
    INDEX idx_student_promotions_from_year (from_academic_year_id),
    
    CONSTRAINT fk_sp_student FOREIGN KEY (student_id) 
        REFERENCES students(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sp_from_class FOREIGN KEY (from_class_id) 
        REFERENCES classes(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_sp_to_class FOREIGN KEY (to_class_id) 
        REFERENCES classes(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_sp_from_year FOREIGN KEY (from_academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sp_to_year FOREIGN KEY (to_academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
