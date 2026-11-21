-- Migration V2: Tạo bảng classes và class_members cho module M3 (L01, L12, AC-L1.5, BR-M3.05, BR-M3.26)

CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    semester VARCHAR(32) NOT NULL,
    lead_instructor_id UUID NOT NULL,
    quota_account_ref VARCHAR(80) NOT NULL,
    weekly_minutes INT NOT NULL DEFAULT 1200,
    max_concurrent_sessions INT NOT NULL DEFAULT 30,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_classes_code UNIQUE (code),
    CONSTRAINT uk_classes_quota_account_ref UNIQUE (quota_account_ref)
);

CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    user_id UUID,
    role VARCHAR(32) NOT NULL,
    student_id_number VARCHAR(32),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    locale VARCHAR(8) NOT NULL DEFAULT 'vi',
    joined_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_class_member_student_id UNIQUE (class_id, student_id_number)
);

CREATE UNIQUE INDEX IF NOT EXISTS uk_class_member_email_lower ON class_members (class_id, LOWER(email));
CREATE INDEX IF NOT EXISTS idx_class_members_user ON class_members (user_id);
CREATE INDEX IF NOT EXISTS idx_class_members_unlinked_email ON class_members (LOWER(email)) WHERE user_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_classes_status ON classes (status);
