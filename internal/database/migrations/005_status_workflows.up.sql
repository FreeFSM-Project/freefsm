CREATE TABLE status_workflows (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    object_type TEXT NOT NULL
        CHECK (object_type IN ('job', 'invoice', 'estimate')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE statuses (
    id BIGSERIAL PRIMARY KEY,
    workflow_id BIGINT NOT NULL REFERENCES status_workflows(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '#6B7280',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_statuses_workflow_id ON statuses(workflow_id);

-- Default workflows
INSERT INTO status_workflows (name, object_type) VALUES ('Default Job Workflow', 'job');
INSERT INTO statuses (workflow_id, name, color, sort_order) VALUES
    (1, 'New', '#3B82F6', 1),
    (1, 'Scheduled', '#8B5CF6', 2),
    (1, 'In Progress', '#F59E0B', 3),
    (1, 'Completed', '#10B981', 4),
    (1, 'Canceled', '#EF4444', 5);

INSERT INTO status_workflows (name, object_type) VALUES ('Default Invoice Workflow', 'invoice');
INSERT INTO statuses (workflow_id, name, color, sort_order) VALUES
    (2, 'Draft', '#6B7280', 1),
    (2, 'Sent', '#3B82F6', 2),
    (2, 'Paid', '#10B981', 3),
    (2, 'Partially Paid', '#F59E0B', 4),
    (2, 'Void', '#EF4444', 5);

INSERT INTO status_workflows (name, object_type) VALUES ('Default Estimate Workflow', 'estimate');
INSERT INTO statuses (workflow_id, name, color, sort_order) VALUES
    (3, 'Draft', '#6B7280', 1),
    (3, 'Sent', '#3B82F6', 2),
    (3, 'Accepted', '#10B981', 3),
    (3, 'Declined', '#EF4444', 4);
