-- ============================================================
-- Migration 004: Tabel Teachers (Guru)
-- Profil detail guru, terhubung 1:1 ke users
-- ============================================================

CREATE TABLE IF NOT EXISTS teachers (
    id              CHAR(36) PRIMARY KEY,
    user_id         CHAR(36) NOT NULL,
    nip             VARCHAR(30) NOT NULL,
    full_name       VARCHAR(255) NOT NULL,
    gender          ENUM('L', 'P') NOT NULL,
    birth_place     VARCHAR(100),
    birth_date      DATE,
    phone           VARCHAR(20),
    address_region  VARCHAR(255),
    address_detail  TEXT,
    photo           TEXT,
    status          ENUM('ACTIVE', 'INACTIVE') NOT NULL DEFAULT 'ACTIVE',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_teachers_user_id (user_id),
    UNIQUE KEY uk_teachers_nip (nip),
    INDEX idx_teachers_status (status),
    INDEX idx_teachers_gender (gender),

    CONSTRAINT fk_teachers_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
