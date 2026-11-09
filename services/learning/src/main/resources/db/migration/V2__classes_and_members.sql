-- Migration V2: Tạo bảng classes và class_members cho module M3 (L01, L12, AC-L1.5)

CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    semester VARCHAR(32) NOT NULL,
    instructor_quota_account_ref VARCHAR(64),
    session_limit_per_student INT NOT NULL DEFAULT 1,
    status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_classes_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS class_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    user_id UUID,
    role VARCHAR(32) NOT NULL,
    student_id_number VARCHAR(32),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_class_member_email UNIQUE (class_id, email),
    CONSTRAINT uk_class_member_student_id UNIQUE (class_id, student_id_number)
);

CREATE INDEX idx_class_members_user_id ON class_members(user_id);
CREATE INDEX idx_class_members_email ON class_members(email);
CREATE INDEX idx_classes_status ON classes(status);
