CREATE UNIQUE INDEX IF NOT EXISTS spec_artifact_pkey_uidx ON spec_artifact (id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_spec_artifact_change_version_path
    ON spec_artifact (change_id, version, relative_path);
