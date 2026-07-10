-- ============================================================
-- Migration 005: Tabel Students (Siswa)
-- Profil detail siswa, terhubung 1:1 ke users
-- ============================================================

CREATE TABLE IF NOT EXISTS students (
    id              CHAR(36) PRIMARY KEY,
    user_id         CHAR(36) NOT NULL,
    nis             VARCHAR(20) NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    gender          ENUM('L', 'P') NOT NULL,
    birth_place     VARCHAR(100),
    birth_date      DATE,
    phone           VARCHAR(20),
    address_region  VARCHAR(255),
    address_detail  TEXT,
    photo           TEXT,
    major_id        CHAR(36),
    class_level     VARCHAR(10),
    status          ENUM('ACTIVE', 'GRADUATED', 'TRANSFER', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_students_user_id (user_id),
    UNIQUE KEY uk_students_nis (nis),
    INDEX idx_students_status (status),
    INDEX idx_students_major (major_id),
    INDEX idx_students_class_level (class_level),
    INDEX idx_students_gender (gender),

    CONSTRAINT fk_students_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_students_major FOREIGN KEY (major_id) 
        REFERENCES majors(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_students_level FOREIGN KEY (class_level) 
        REFERENCES levels(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
