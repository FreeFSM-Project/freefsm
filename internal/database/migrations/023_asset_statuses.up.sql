CREATE TABLE asset_statuses (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT,
    name TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '#6B7280',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_asset_statuses_company_id ON asset_statuses(company_id);

-- Default asset statuses
INSERT INTO asset_statuses (name, color, sort_order) VALUES
    ('Active', '#10B981', 1),
    ('Inactive', '#6B7280', 2),
    ('Maintenance', '#F59E0B', 3),
    ('Retired', '#EF4444', 4);
