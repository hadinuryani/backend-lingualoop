-- ============================================================
-- Migration 003: Tabel Levels (Tingkatan Kelas)
-- Data statis: Kelas X, XI, XII
-- ============================================================

CREATE TABLE IF NOT EXISTS levels (
    id      VARCHAR(10) PRIMARY KEY,
    name    VARCHAR(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Seed data tingkatan kelas
INSERT INTO levels (id, name) VALUES 
    ('lvl-10', 'X'),
    ('lvl-11', 'XI'),
    ('lvl-12', 'XII');
