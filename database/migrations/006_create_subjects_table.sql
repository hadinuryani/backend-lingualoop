-- ============================================================
-- Migration 006: Tabel Subjects (Mata Pelajaran)
-- Kurikulum mata pelajaran per jurusan dan tingkat kelas
-- ============================================================

CREATE TABLE IF NOT EXISTS subjects (
    id              CHAR(36) PRIMARY KEY,
    code            VARCHAR(20) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    major_id        CHAR(36),              -- NULL = mata pelajaran umum (semua jurusan)
    level_id        VARCHAR(10),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_subjects_code (code),
    INDEX idx_subjects_major (major_id),
    INDEX idx_subjects_level (level_id),

    CONSTRAINT fk_subjects_major FOREIGN KEY (major_id) 
        REFERENCES majors(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_subjects_level FOREIGN KEY (level_id) 
        REFERENCES levels(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
