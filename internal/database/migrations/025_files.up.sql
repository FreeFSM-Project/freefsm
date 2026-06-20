CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT,
    object_type TEXT NOT NULL,
    object_id BIGINT NOT NULL,
    original_name TEXT NOT NULL,
    stored_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    uploaded_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_files_object ON files(object_type, object_id);
CREATE INDEX idx_files_created ON files(created_at);
