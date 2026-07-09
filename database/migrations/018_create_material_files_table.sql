-- ============================================================
-- Migration 018: Tabel Material Files (File Lampiran Materi)
-- 1 materi bisa punya banyak file (PDF, PPT, Video, dll)
-- ============================================================

CREATE TABLE IF NOT EXISTS material_files (
    id              CHAR(36) PRIMARY KEY,
    material_id     CHAR(36) NOT NULL,
    file_url        TEXT NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_type       VARCHAR(50),                   -- 'pdf', 'ppt', 'pptx', 'doc', 'mp4', 'jpg', dll
    file_size       BIGINT,                        -- ukuran dalam bytes
    sort_order      INT NOT NULL DEFAULT 0,
    uploaded_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_mf_material (material_id),
    INDEX idx_mf_sort (sort_order),

    CONSTRAINT fk_mf_material FOREIGN KEY (material_id) 
        REFERENCES materials(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
