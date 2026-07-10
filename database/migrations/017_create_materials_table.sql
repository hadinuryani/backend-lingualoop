-- ============================================================
-- Migration 017: Tabel Materials (Materi Pembelajaran)
-- Materi yang diupload/ditulis guru (file + rich-text)
-- ============================================================

CREATE TABLE IF NOT EXISTS materials (
    id                  CHAR(36) PRIMARY KEY,
    teacher_id          CHAR(36) NOT NULL,
    subject_id          CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    semester            ENUM('GANJIL', 'GENAP') NOT NULL,

    title               VARCHAR(255) NOT NULL,
    description         TEXT,

    -- Konten materi: FILE saja, TEXT saja, atau BOTH
    content_type        ENUM('FILE', 'TEXT', 'BOTH') NOT NULL DEFAULT 'FILE',
    content_text        MEDIUMTEXT,                            -- Rich-text / HTML content

    sort_order          INT NOT NULL DEFAULT 0,
    status              ENUM('DRAFT', 'PUBLISHED', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT',
    published_at        DATETIME NULL,

    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    INDEX idx_materials_teacher (teacher_id),
    INDEX idx_materials_subject (subject_id),
    INDEX idx_materials_semester (academic_year_id, semester),
    INDEX idx_materials_status (status),
    INDEX idx_materials_sort (sort_order),

    CONSTRAINT fk_materials_teacher FOREIGN KEY (teacher_id) 
        REFERENCES teachers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_materials_subject FOREIGN KEY (subject_id) 
        REFERENCES subjects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_materials_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
