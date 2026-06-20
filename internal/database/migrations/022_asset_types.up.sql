CREATE TABLE asset_types (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT,
    name TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_asset_types_company_id ON asset_types(company_id);

-- Default asset types
INSERT INTO asset_types (name, sort_order) VALUES
    ('HVAC', 1),
    ('Generator', 2),
    ('Water Heater', 3),
    ('Refrigeration', 4),
    ('Other', 5);
