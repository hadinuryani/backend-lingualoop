-- ============================================================
-- Migration 001: Tabel Users
-- Tabel autentikasi sentral untuk semua role (admin, teacher, student)
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id              CHAR(36) PRIMARY KEY,
    email           VARCHAR(255) NOT NULL,
    username        VARCHAR(100) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    role            ENUM('admin', 'teacher', 'student') NOT NULL,
    avatar_url      TEXT,
    is_active       BOOLEAN DEFAULT TRUE,
    last_login_at   DATETIME NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_users_email (email),
    UNIQUE KEY uk_users_username (username),
    INDEX idx_users_role (role),
    INDEX idx_users_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
