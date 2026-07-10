-- ============================================================
-- Migration 019: Tabel Material Classes (Materi ↔ Kelas Target)
-- Pivot: guru bisa share 1 materi ke beberapa kelas sekaligus
-- ============================================================

CREATE TABLE IF NOT EXISTS material_classes (
    id              CHAR(36) PRIMARY KEY,
    material_id     CHAR(36) NOT NULL,
    class_id        CHAR(36) NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    UNIQUE KEY uk_material_classes (material_id, class_id),
    INDEX idx_mc_material (material_id),
    INDEX idx_mc_class (class_id),

    CONSTRAINT fk_mc_material FOREIGN KEY (material_id) 
        REFERENCES materials(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_mc_class FOREIGN KEY (class_id) 
        REFERENCES classes(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
