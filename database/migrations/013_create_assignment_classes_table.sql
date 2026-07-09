-- ============================================================
-- Migration 013: Tabel Assignment Classes (Tugas ↔ Kelas Target)
-- Pivot: guru bisa assign 1 tugas ke 1 atau beberapa kelas sekaligus
-- ============================================================

CREATE TABLE IF NOT EXISTS assignment_classes (
    id              CHAR(36) PRIMARY KEY,
    assignment_id   CHAR(36) NOT NULL,
    class_id        CHAR(36) NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_assignment_classes (assignment_id, class_id),
    INDEX idx_ac_assignment (assignment_id),
    INDEX idx_ac_class (class_id),

    CONSTRAINT fk_ac_assignment FOREIGN KEY (assignment_id) 
        REFERENCES assignments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_ac_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
