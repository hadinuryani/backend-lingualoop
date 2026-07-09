-- ============================================================
-- Migration 020: Tabel Grade Components (Komponen Penilaian)
-- Bobot penilaian per komponen per mapel per kelas per semester
-- Contoh: Tugas Harian 30%, UTS 30%, UAS 40%
-- ============================================================

CREATE TABLE IF NOT EXISTS grade_components (
    id                  CHAR(36) PRIMARY KEY,
    subject_id          CHAR(36) NOT NULL,
    class_id            CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    semester            ENUM('GANJIL', 'GENAP') NOT NULL,

    component_name      VARCHAR(100) NOT NULL,         -- 'Tugas Harian', 'UTS', 'UAS', 'Praktik'
    assignment_type     ENUM('MINGGUAN', 'TUGAS_HARIAN', 'UTS', 'UAS', 'QUIZ', 'PRAKTIK', 'LAINNYA') 
                        NOT NULL,
    weight_percentage   DECIMAL(5,2) NOT NULL,         -- bobot persentase (e.g., 30.00)

    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 1 tipe komponen per mapel per kelas per semester
    UNIQUE KEY uk_grade_components (subject_id, class_id, academic_year_id, semester, assignment_type),
    INDEX idx_gc_subject (subject_id),
    INDEX idx_gc_class (class_id),
    INDEX idx_gc_semester (academic_year_id, semester),

    CONSTRAINT fk_gc_subject FOREIGN KEY (subject_id) 
        REFERENCES subjects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_gc_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_gc_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
