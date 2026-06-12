CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status_id BIGINT REFERENCES statuses(id) ON DELETE SET NULL,
    location_id BIGINT REFERENCES locations(id) ON DELETE SET NULL,
    completion_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_customer_id ON projects(customer_id);

CREATE TABLE jobs (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    project_id BIGINT REFERENCES projects(id) ON DELETE SET NULL,
    location_id BIGINT REFERENCES locations(id) ON DELETE SET NULL,
    customer_contact_id BIGINT REFERENCES customer_contacts(id) ON DELETE SET NULL,
    job_type TEXT NOT NULL,
    subtitle TEXT NOT NULL DEFAULT '',
    status_id BIGINT REFERENCES statuses(id) ON DELETE SET NULL,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    due_date TIMESTAMPTZ,
    arrival_window_start TIMESTAMPTZ,
    arrival_window_end TIMESTAMPTZ,
    notes TEXT NOT NULL DEFAULT '',
    tech_notes TEXT NOT NULL DEFAULT '',
    billing_type TEXT NOT NULL DEFAULT 'flat_rate'
        CHECK (billing_type IN ('flat_rate', 'hourly', 't_and_m')),
    visits JSONB NOT NULL DEFAULT '[]',
    assignments JSONB NOT NULL DEFAULT '[]',
    custom_fields JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_customer_id ON jobs(customer_id);
CREATE INDEX idx_jobs_project_id ON jobs(project_id);
CREATE INDEX idx_jobs_status_id ON jobs(status_id);
CREATE INDEX idx_jobs_start_time ON jobs(start_time);
