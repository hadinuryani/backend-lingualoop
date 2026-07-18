-- ============================================================
-- Migration 025: Add Logo File ID to Majors
-- ============================================================

ALTER TABLE majors ADD COLUMN logo_file_id CHAR(36) NULL;

-- Asumsikan kita butuh indeks untuk optimasi
CREATE INDEX idx_majors_logo_file_id ON majors (logo_file_id);

-- Optional: Foreign key ke files.id jika database mendukung referential integrity ketat
-- ALTER TABLE majors ADD CONSTRAINT fk_majors_logo FOREIGN KEY (logo_file_id) REFERENCES files(id) ON DELETE SET NULL;
