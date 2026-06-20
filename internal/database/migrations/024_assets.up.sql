CREATE TABLE assets (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    location_id BIGINT REFERENCES locations(id) ON DELETE SET NULL,
    asset_type_id BIGINT NOT NULL REFERENCES asset_types(id) ON DELETE SET NULL,
    asset_status_id BIGINT REFERENCES asset_statuses(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    serial_number TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    manufacturer TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    installed_at TIMESTAMPTZ,
    warranty_expires TIMESTAMPTZ,
    custom_fields JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_assets_customer_id ON assets(customer_id);
CREATE INDEX idx_assets_location_id ON assets(location_id);
CREATE INDEX idx_assets_asset_type_id ON assets(asset_type_id);
CREATE INDEX idx_assets_asset_status_id ON assets(asset_status_id);
CREATE INDEX idx_assets_company_id ON assets(company_id);

-- Add asset_id to jobs for service history linking
ALTER TABLE jobs ADD COLUMN asset_id BIGINT REFERENCES assets(id) ON DELETE SET NULL;
CREATE INDEX idx_jobs_asset_id ON jobs(asset_id);
