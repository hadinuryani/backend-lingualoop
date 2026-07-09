-- ============================================================
-- Migration 009: Tabel Student Classes (Siswa ↔ Kelas)
-- Pivot: penempatan siswa ke kelas per tahun akademik
-- Constraint: 1 siswa hanya bisa di 1 kelas per tahun ajaran
-- ============================================================

CREATE TABLE IF NOT EXISTS student_classes (
    id                  CHAR(36) PRIMARY KEY,
    student_id          CHAR(36) NOT NULL,
    class_id            CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- 1 siswa = 1 kelas per tahun ajaran
    UNIQUE KEY uk_student_classes_student_year (student_id, academic_year_id),
    INDEX idx_student_classes_class (class_id),
    INDEX idx_student_classes_academic_year (academic_year_id),
    INDEX idx_student_classes_is_active (is_active),

    CONSTRAINT fk_sc_student FOREIGN KEY (student_id) 
        REFERENCES students(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sc_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sc_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
