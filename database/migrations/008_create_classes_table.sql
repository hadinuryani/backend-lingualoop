-- ============================================================
-- Migration 008: Tabel Classes (Kelas / Rombongan Belajar)
-- Kelas per tahun akademik, terhubung ke jurusan dan tingkat
-- ============================================================

CREATE TABLE IF NOT EXISTS classes (
    id                      CHAR(36) PRIMARY KEY,
    academic_year_id        CHAR(36) NOT NULL,
    major_id                CHAR(36),
    level_id                VARCHAR(10) NOT NULL,
    class_name              VARCHAR(100) NOT NULL,
    classroom               VARCHAR(100),             -- nama ruangan fisik
    capacity                INT NOT NULL DEFAULT 36,
    homeroom_teacher_id     CHAR(36),                 -- wali kelas
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_classes_name_year (class_name, academic_year_id),
    INDEX idx_classes_academic_year (academic_year_id),
    INDEX idx_classes_major (major_id),
    INDEX idx_classes_level (level_id),
    INDEX idx_classes_homeroom (homeroom_teacher_id),
    INDEX idx_classes_is_active (is_active),

    CONSTRAINT fk_classes_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_classes_major FOREIGN KEY (major_id) 
        REFERENCES majors(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_classes_level FOREIGN KEY (level_id) 
        REFERENCES levels(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_classes_homeroom FOREIGN KEY (homeroom_teacher_id) 
        REFERENCES teachers(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
