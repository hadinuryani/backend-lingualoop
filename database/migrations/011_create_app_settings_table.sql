-- ============================================================
-- Migration 011: Tabel App Settings (Pengaturan Aplikasi)
-- Key-value store untuk konfigurasi sekolah dan aplikasi
-- ============================================================

CREATE TABLE IF NOT EXISTS app_settings (
    `key`           VARCHAR(100) PRIMARY KEY,
    value           TEXT NOT NULL,
    description     VARCHAR(255),
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed data pengaturan default
INSERT INTO app_settings (`key`, value, description) VALUES
    ('school_name', 'SMK Negeri 1 LinguaLoop', 'Nama sekolah'),
    ('school_npsn', '20100001', 'NPSN sekolah'),
    ('school_address', 'Jl. Pendidikan No. 1, Kota Bandung, Jawa Barat', 'Alamat sekolah'),
    ('school_phone', '(022) 1234567', 'Telepon sekolah'),
    ('school_email', 'info@smkn1lingualoop.sch.id', 'Email sekolah'),
    ('default_teacher_password', 'guru123', 'Password default akun guru baru'),
    ('default_student_password', 'siswa123', 'Password default akun siswa baru'),
    ('max_students_per_class', '36', 'Kapasitas maksimal siswa per kelas'),
    ('grading_system', 'numeric', 'Sistem penilaian: numeric atau letter'),
    ('passing_grade', '75', 'Nilai KKM default'),
    ('app_name', 'LinguaLoop', 'Nama aplikasi'),
    ('enable_student_login', 'true', 'Aktifkan login siswa'),
    ('enable_teacher_login', 'true', 'Aktifkan login guru'),
    ('maintenance_mode', 'false', 'Mode maintenance');
