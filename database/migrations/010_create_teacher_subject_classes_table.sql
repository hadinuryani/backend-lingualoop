-- ============================================================
-- Migration 010: Tabel Teacher Subject Classes (Guru ↔ Mapel ↔ Kelas)
-- Pivot: penugasan guru mengajar mata pelajaran di kelas tertentu
-- Constraint: 1 mapel = 1 guru per kelas per tahun ajaran
-- ============================================================

CREATE TABLE IF NOT EXISTS teacher_subject_classes (
    id                  CHAR(36) PRIMARY KEY,
    teacher_id          CHAR(36) NOT NULL,
    subject_id          CHAR(36) NOT NULL,
    class_id            CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    classroom           VARCHAR(100),                -- ruangan mengajar
    status              VARCHAR(20) NOT NULL DEFAULT 'Aktif',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    -- 1 mapel = 1 guru per kelas per tahun ajaran
    UNIQUE KEY uk_tsc_subject_class_year (subject_id, class_id, academic_year_id),
    INDEX idx_tsc_teacher (teacher_id),
    INDEX idx_tsc_subject (subject_id),
    INDEX idx_tsc_class (class_id),
    INDEX idx_tsc_academic_year (academic_year_id),

    CONSTRAINT fk_tsc_teacher FOREIGN KEY (teacher_id) 
        REFERENCES teachers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_tsc_subject FOREIGN KEY (subject_id) 
        REFERENCES subjects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_tsc_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_tsc_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
