CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    object_type TEXT NOT NULL,
    object_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    address_1 TEXT NOT NULL DEFAULT '',
    address_2 TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL DEFAULT '',
    zip_code TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    is_primary BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_locations_object ON locations(object_type, object_id);
