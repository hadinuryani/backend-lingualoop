-- ============================================================
-- Migration 031: Add composite indexes for Teacher Portal
-- Optimizing queries that filter by teacher_id and academic_year_id
-- ============================================================

-- Drop the old single indexes if we want to replace them, 
-- but adding composite indexes alongside is also fine.
-- Using composite index for teacher portal queries:

CREATE INDEX idx_tsc_teacher_academic ON teacher_subject_classes (teacher_id, academic_year_id);

CREATE INDEX idx_schedules_teacher_academic_day_period ON schedules (teacher_id, academic_year_id, day, period);
