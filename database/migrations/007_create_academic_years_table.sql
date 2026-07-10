-- ============================================================
-- Migration 007: Tabel Academic Years (Tahun Akademik)
-- Tahun ajaran dengan konfigurasi semester ganjil dan genap
-- ============================================================

CREATE TABLE IF NOT EXISTS academic_years (
    id                          CHAR(36) PRIMARY KEY,
    year                        VARCHAR(20) NOT NULL,          -- '2025/2026'
    start_date                  DATE NOT NULL,
    end_date                    DATE NOT NULL,
    status                      ENUM('Draft', 'Aktif', 'Menunggu Kenaikan', 'Selesai') 
                                NOT NULL DEFAULT 'Draft',

    -- Semester Ganjil
    sem_ganjil_start_date       DATE,
    sem_ganjil_end_kbm          DATE,
    sem_ganjil_end_assessment   DATE,
    sem_ganjil_status           ENUM('Belum Aktif', 'Aktif', 'Masa Penilaian', 'Siap Ditutup', 'Terkunci') 
                                NOT NULL DEFAULT 'Belum Aktif',

    -- Semester Genap
    sem_genap_start_date        DATE,
    sem_genap_end_kbm           DATE,
    sem_genap_end_assessment    DATE,
    sem_genap_status            ENUM('Belum Aktif', 'Aktif', 'Masa Penilaian', 'Siap Ditutup', 'Terkunci') 
                                NOT NULL DEFAULT 'Belum Aktif',

    created_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_academic_years_year (year),
    INDEX idx_academic_years_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
