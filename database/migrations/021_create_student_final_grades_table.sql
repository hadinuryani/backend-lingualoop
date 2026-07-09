-- ============================================================
-- Migration 021: Tabel Student Final Grades (Nilai Akhir / Rapor)
-- Nilai akhir siswa per mapel per semester, dihitung otomatis
-- dari weighted average komponen penilaian
-- ============================================================

CREATE TABLE IF NOT EXISTS student_final_grades (
    id                  CHAR(36) PRIMARY KEY,
    student_id          CHAR(36) NOT NULL,
    subject_id          CHAR(36) NOT NULL,
    class_id            CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    semester            ENUM('GANJIL', 'GENAP') NOT NULL,

    -- Nilai per komponen (disimpan sebagai JSON)
    -- Format: {"TUGAS_HARIAN": 85.5, "UTS": 78.0, "UAS": 90.0}
    component_scores    JSON,

    final_score         DECIMAL(5,2),                  -- Nilai akhir (weighted average)
    grade_letter        VARCHAR(5),                    -- 'A', 'B+', 'B', 'C+', 'C', 'D', 'E'
    is_passing          BOOLEAN,                       -- Lulus KKM atau tidak

    teacher_notes       TEXT,                          -- Catatan guru di rapor
    finalized_by        CHAR(36),
    finalized_at        DATETIME NULL,

    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 1 nilai akhir per siswa per mapel per semester
    UNIQUE KEY uk_student_final_grades (student_id, subject_id, academic_year_id, semester),
    INDEX idx_sfg_student (student_id),
    INDEX idx_sfg_subject (subject_id),
    INDEX idx_sfg_class (class_id),
    INDEX idx_sfg_semester (academic_year_id, semester),
    INDEX idx_sfg_finalized_by (finalized_by),
    INDEX idx_sfg_is_passing (is_passing),

    CONSTRAINT fk_sfg_student FOREIGN KEY (student_id) 
        REFERENCES students(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sfg_subject FOREIGN KEY (subject_id) 
        REFERENCES subjects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sfg_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sfg_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sfg_finalized_by FOREIGN KEY (finalized_by) 
        REFERENCES teachers(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
