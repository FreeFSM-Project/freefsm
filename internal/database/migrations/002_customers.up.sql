CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    display_name TEXT NOT NULL,
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    company_name TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'lead'
        CHECK (status IN ('lead','opportunity','customer','lost','inactive')),
    account_type TEXT NOT NULL DEFAULT 'individual'
        CHECK (account_type IN ('individual','company')),
    assigned_to BIGINT REFERENCES users(id) ON DELETE SET NULL,
    pipeline_status_id BIGINT,
    lead_source_id BIGINT,

    billing_address_1 TEXT NOT NULL DEFAULT '',
    billing_address_2 TEXT NOT NULL DEFAULT '',
    billing_city TEXT NOT NULL DEFAULT '',
    billing_state TEXT NOT NULL DEFAULT '',
    billing_zip_code TEXT NOT NULL DEFAULT '',

    service_address_1 TEXT NOT NULL DEFAULT '',
    service_address_2 TEXT NOT NULL DEFAULT '',
    service_city TEXT NOT NULL DEFAULT '',
    service_state TEXT NOT NULL DEFAULT '',
    service_zip_code TEXT NOT NULL DEFAULT '',

    custom_fields JSONB NOT NULL DEFAULT '[]',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_display_name ON customers(display_name);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_assigned_to ON customers(assigned_to);

CREATE TABLE customer_contacts (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customer_contacts_customer_id ON customer_contacts(customer_id);
