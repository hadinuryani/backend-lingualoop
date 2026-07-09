-- ============================================================
-- Migration 012: Tabel Assignments (Tugas)
-- Tugas yang dibuat guru (Mingguan, UTS, UAS, Quiz, dll)
-- Mendukung mode soal STATIC (sama) dan RANDOM (acak per siswa)
-- Mendukung tipe pengumpulan FILE, QUIZ, atau BOTH
-- ============================================================

CREATE TABLE IF NOT EXISTS assignments (
    id                  CHAR(36) PRIMARY KEY,
    teacher_id          CHAR(36) NOT NULL,
    subject_id          CHAR(36) NOT NULL,
    academic_year_id    CHAR(36) NOT NULL,
    semester            ENUM('GANJIL', 'GENAP') NOT NULL,

    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    type                ENUM('MINGGUAN', 'TUGAS_HARIAN', 'UTS', 'UAS', 'QUIZ', 'PRAKTIK', 'LAINNYA') 
                        NOT NULL,

    -- Mode soal: STATIC = sama untuk semua, RANDOM = acak per siswa dari bank soal
    question_mode       ENUM('STATIC', 'RANDOM') NOT NULL DEFAULT 'STATIC',

    -- Cara pengumpulan: FILE = upload file, QUIZ = soal interaktif, BOTH = gabungan
    submission_type     ENUM('FILE', 'QUIZ', 'BOTH') NOT NULL DEFAULT 'FILE',

    -- File lampiran soal (untuk mode FILE/BOTH)
    attachment_url      TEXT,
    attachment_name     VARCHAR(255),

    -- Konfigurasi deadline: NONE, SOFT (bisa telat), HARD (tidak bisa telat)
    deadline_mode       ENUM('NONE', 'SOFT', 'HARD') NOT NULL DEFAULT 'SOFT',
    deadline_at         DATETIME NULL,

    -- Skor & bobot
    max_score           INT NOT NULL DEFAULT 100,
    passing_score       INT NOT NULL DEFAULT 75,

    -- Jumlah soal random per siswa (hanya untuk mode RANDOM)
    random_question_count INT NULL,

    status              ENUM('DRAFT', 'PUBLISHED', 'CLOSED', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT',
    published_at        DATETIME NULL,

    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_assignments_teacher (teacher_id),
    INDEX idx_assignments_subject (subject_id),
    INDEX idx_assignments_semester (academic_year_id, semester),
    INDEX idx_assignments_type (type),
    INDEX idx_assignments_status (status),
    INDEX idx_assignments_deadline (deadline_at),

    CONSTRAINT fk_assignments_teacher FOREIGN KEY (teacher_id) 
        REFERENCES teachers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_assignments_subject FOREIGN KEY (subject_id) 
        REFERENCES subjects(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_assignments_academic_year FOREIGN KEY (academic_year_id) 
        REFERENCES academic_years(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
