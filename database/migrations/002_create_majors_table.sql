-- ============================================================
-- Migration 002: Tabel Majors (Jurusan)
-- Menyimpan data jurusan/program studi (Bahasa Inggris, Jepang, dll)
-- ============================================================

CREATE TABLE IF NOT EXISTS majors (
    id              CHAR(36) PRIMARY KEY,
    code            VARCHAR(20) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_majors_code (code),
    UNIQUE KEY uk_majors_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
