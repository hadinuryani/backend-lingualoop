-- ============================================================
-- Migration 022: Tabel Activity Logs (Log Aktivitas)
-- Audit trail untuk semua aksi pengguna di sistem
-- Digunakan untuk dashboard dan riwayat aktivitas
-- ============================================================

CREATE TABLE IF NOT EXISTS activity_logs (
    id              CHAR(36) PRIMARY KEY,
    user_id         CHAR(36),
    action          VARCHAR(100) NOT NULL,             -- 'CREATE_ASSIGNMENT', 'SUBMIT_ANSWER', dll
    target_type     VARCHAR(50),                       -- 'assignment', 'material', 'student', dll
    target_id       CHAR(36),
    details         JSON,                              -- Detail tambahan dalam JSON
    ip_address      VARCHAR(45),
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    INDEX idx_al_user (user_id),
    INDEX idx_al_action (action),
    INDEX idx_al_target (target_type, target_id),
    INDEX idx_al_created (created_at DESC),

    CONSTRAINT fk_al_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
