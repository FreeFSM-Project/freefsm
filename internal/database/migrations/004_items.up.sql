CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'service'
        CHECK (type IN ('service', 'product')),
    sku TEXT NOT NULL DEFAULT '',
    unit_price NUMERIC(12,2) NOT NULL DEFAULT 0,
    unit_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
    taxable BOOLEAN NOT NULL DEFAULT true,
    tax_rate TEXT NOT NULL DEFAULT '',
    track_inventory BOOLEAN NOT NULL DEFAULT false,
    description TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_items_name ON items(name);
CREATE UNIQUE INDEX idx_items_sku ON items(sku) WHERE sku != '';
