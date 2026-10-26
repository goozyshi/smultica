CREATE UNIQUE INDEX IF NOT EXISTS spec_change_pkey_uidx ON spec_change (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_spec_change_project_key
    ON spec_change (project_id, change_key);
CREATE INDEX IF NOT EXISTS idx_spec_change_workspace_project
    ON spec_change (workspace_id, project_id, updated_at DESC);
