-- OpenSpec-backed Change metadata. The repository remains the source of truth
-- for artifact contents; this table only indexes local-repository work.
CREATE TABLE spec_change (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL,
    project_id UUID NOT NULL REFERENCES project(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES project_resource(id) ON DELETE SET NULL,
    change_key TEXT NOT NULL CHECK (char_length(change_key) BETWEEN 1 AND 128),
    title TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'discovered'
        CHECK (status IN ('discovered', 'draft', 'active', 'blocked', 'completed', 'archived')),
    current_version INTEGER NOT NULL DEFAULT 1 CHECK (current_version > 0),
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
