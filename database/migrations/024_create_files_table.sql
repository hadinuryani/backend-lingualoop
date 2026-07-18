-- ============================================================
-- Migration 024: Create Files Table
-- ============================================================

CREATE TABLE IF NOT EXISTS files (
    id CHAR(36) PRIMARY KEY,
    resource_type VARCHAR(50) NOT NULL, -- e.g. 'majors', 'teachers', 'materials'
    storage_path VARCHAR(500) NOT NULL, -- e.g. 'public/majors/uuid.png'
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    uploaded_by CHAR(36) NULL, -- Can be linked to User ID
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    
    KEY idx_files_resource_type (resource_type),
    KEY idx_files_uploaded_by (uploaded_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
