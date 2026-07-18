-- ============================================================
-- Migration 027: Create Schedules Table
-- ============================================================

CREATE TABLE IF NOT EXISTS schedules (
    id CHAR(36) PRIMARY KEY,
    academic_year_id VARCHAR(50) NOT NULL,
    class_id CHAR(36) NOT NULL,
    subject_id CHAR(36) NOT NULL,
    teacher_id CHAR(36) NOT NULL,
    day VARCHAR(20) NOT NULL, -- e.g., 'Senin'
    period INT NOT NULL,      -- e.g., jam ke-1
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    
    -- Ensure a class cannot have two subjects at the exact same day and period for an academic year
    UNIQUE KEY idx_class_time (academic_year_id, class_id, day, period),
    
    -- Fast lookups
    KEY idx_schedules_teacher (teacher_id),
    KEY idx_schedules_class (class_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
