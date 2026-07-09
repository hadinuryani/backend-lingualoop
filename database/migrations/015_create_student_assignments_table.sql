-- ============================================================
-- Migration 015: Tabel Student Assignments (Pengerjaan Tugas Siswa)
-- Catatan penugasan per siswa, termasuk soal random yang didapat,
-- file submission, status pengerjaan, dan nilai akhir
-- ============================================================

CREATE TABLE IF NOT EXISTS student_assignments (
    id                      CHAR(36) PRIMARY KEY,
    assignment_id           CHAR(36) NOT NULL,
    student_id              CHAR(36) NOT NULL,

    -- Soal yang didapat siswa ini (untuk RANDOM mode)
    -- Disimpan sebagai JSON array of question IDs, contoh: ["uuid1", "uuid2", ...]
    -- Jika STATIC mode, ini NULL (semua siswa dapat soal yang sama)
    assigned_question_ids   JSON,

    status                  ENUM('PENDING', 'IN_PROGRESS', 'SUBMITTED', 'GRADED', 'LATE') 
                            NOT NULL DEFAULT 'PENDING',

    -- File submission (untuk mode FILE/BOTH)
    submitted_file_url      TEXT,
    submitted_file_name     VARCHAR(255),

    submitted_at            DATETIME NULL,
    is_late                 BOOLEAN NOT NULL DEFAULT FALSE,

    -- Nilai akhir dari guru / auto-grading
    total_score             DECIMAL(5,2),
    graded_by               CHAR(36),
    graded_at               DATETIME NULL,
    teacher_notes           TEXT,                    -- Catatan/feedback guru

    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- 1 siswa = 1 pengerjaan per tugas
    UNIQUE KEY uk_student_assignments (assignment_id, student_id),
    INDEX idx_sa_student (student_id),
    INDEX idx_sa_assignment (assignment_id),
    INDEX idx_sa_status (status),
    INDEX idx_sa_is_late (is_late),
    INDEX idx_sa_graded_by (graded_by),

    CONSTRAINT fk_sa_assignment FOREIGN KEY (assignment_id) 
        REFERENCES assignments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sa_student FOREIGN KEY (student_id) 
        REFERENCES students(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sa_graded_by FOREIGN KEY (graded_by) 
        REFERENCES teachers(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
