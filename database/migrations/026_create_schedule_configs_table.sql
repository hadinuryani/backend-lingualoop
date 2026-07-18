-- ============================================================
-- Migration 026: Create Schedule Configs Table
-- ============================================================

CREATE TABLE IF NOT EXISTS schedule_configs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    periods_per_day INT NOT NULL DEFAULT 10,
    period_duration INT NOT NULL DEFAULT 45,
    start_time VARCHAR(10) NOT NULL DEFAULT '07:00',
    break_after_periods JSON NOT NULL, -- e.g. '[3,6]'
    break_durations JSON NOT NULL,     -- e.g. '[15,15]'
    active_days JSON NOT NULL,         -- e.g. '["Senin", "Selasa", "Rabu", "Kamis", "Jumat"]'
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default row immediately since it's a singleton config
INSERT IGNORE INTO schedule_configs (id, periods_per_day, period_duration, start_time, break_after_periods, break_durations, active_days) 
VALUES (1, 10, 45, '07:00', '[3,6]', '[15,15]', '["Senin","Selasa","Rabu","Kamis","Jumat"]');
