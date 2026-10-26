CREATE TABLE spec_artifact (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    change_id UUID NOT NULL REFERENCES spec_change(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1 CHECK (version > 0),
    kind TEXT NOT NULL CHECK (kind IN ('proposal', 'spec', 'design', 'tasks', 'other')),
    relative_path TEXT NOT NULL CHECK (relative_path <> ''),
    digest TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'unknown'
        CHECK (status IN ('unknown', 'present', 'missing', 'invalid', 'conflict')),
    scanned_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
