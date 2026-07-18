-- ============================================================
-- Migration 029: Tambah grade_level dan class_number
-- ============================================================

-- Tambah grade_level di tabel levels
ALTER TABLE levels ADD COLUMN grade_level INT NOT NULL DEFAULT 0;

-- Update data seed untuk levels (berdasarkan id lvl-10, lvl-11, lvl-12)
UPDATE levels SET grade_level = 10 WHERE id = 'lvl-10';
UPDATE levels SET grade_level = 11 WHERE id = 'lvl-11';
UPDATE levels SET grade_level = 12 WHERE id = 'lvl-12';

-- Tambah grade_level dan class_number di tabel classes
ALTER TABLE classes ADD COLUMN grade_level INT NOT NULL DEFAULT 0;
ALTER TABLE classes ADD COLUMN class_number INT NOT NULL DEFAULT 0;

-- (Opsional) Update grade_level di classes yang sudah ada dengan JOIN ke levels
UPDATE classes c
JOIN levels l ON c.level_id = l.id
SET c.grade_level = l.grade_level;
