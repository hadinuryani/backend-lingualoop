-- ============================================================
-- Migration 014: Tabel Questions (Bank Soal)
-- Soal-soal untuk tugas interaktif (PG, Essay, True/False, Short Answer)
-- Mendukung auto-grading untuk PG & True/False
-- ============================================================

CREATE TABLE IF NOT EXISTS questions (
    id              CHAR(36) PRIMARY KEY,
    assignment_id   CHAR(36) NOT NULL,

    question_type   ENUM('MULTIPLE_CHOICE', 'ESSAY', 'TRUE_FALSE', 'SHORT_ANSWER') NOT NULL,
    question_text   TEXT NOT NULL,
    question_image  TEXT,                          -- URL gambar soal (opsional)

    -- Untuk pilihan ganda (disimpan sebagai JSON)
    -- Format: [{"key": "A", "text": "Jawaban A"}, {"key": "B", "text": "Jawaban B"}, ...]
    options         JSON,

    -- Kunci jawaban untuk auto-grading (PG: "A"/"B"/dll, True-False: "TRUE"/"FALSE")
    correct_answer  VARCHAR(50),

    -- Panduan jawaban / rubrik penilaian (untuk guru menilai essay)
    answer_guide    TEXT,

    points          INT NOT NULL DEFAULT 10,       -- Bobot poin soal
    sort_order      INT NOT NULL DEFAULT 0,        -- Urutan tampil soal

    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME NULL DEFAULT NULL,

    INDEX idx_questions_assignment (assignment_id),
    INDEX idx_questions_type (question_type),
    INDEX idx_questions_sort (sort_order),

    CONSTRAINT fk_questions_assignment FOREIGN KEY (assignment_id) 
        REFERENCES assignments(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
