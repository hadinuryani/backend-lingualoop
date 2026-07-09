-- ============================================================
-- Migration 016: Tabel Student Answers (Jawaban Siswa per Soal)
-- Jawaban siswa untuk setiap soal dalam tugas interaktif (QUIZ)
-- Mendukung auto-grading untuk PG & True/False
-- ============================================================

CREATE TABLE IF NOT EXISTS student_answers (
    id                      CHAR(36) PRIMARY KEY,
    student_assignment_id   CHAR(36) NOT NULL,
    question_id             CHAR(36) NOT NULL,

    answer_text             TEXT,                    -- Jawaban tertulis (essay/short answer)
    selected_option         VARCHAR(10),             -- Key pilihan ganda ("A", "B", dll)
    answer_file_url         TEXT,                    -- File lampiran (opsional)

    is_correct              BOOLEAN,                 -- Hasil auto-grading (PG/True-False)
    score                   DECIMAL(5,2),            -- Nilai per soal
    teacher_feedback        TEXT,                    -- Feedback guru per soal

    answered_at             DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- 1 jawaban per soal per pengerjaan
    UNIQUE KEY uk_student_answers (student_assignment_id, question_id),
    INDEX idx_sa_answers_sa (student_assignment_id),
    INDEX idx_sa_answers_question (question_id),
    INDEX idx_sa_answers_is_correct (is_correct),

    CONSTRAINT fk_sa_answers_sa FOREIGN KEY (student_assignment_id) 
        REFERENCES student_assignments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_sa_answers_question FOREIGN KEY (question_id) 
        REFERENCES questions(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
