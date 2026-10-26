-- name: ListSpecChangesForProject :many
SELECT * FROM spec_change
WHERE workspace_id = $1 AND project_id = $2
ORDER BY updated_at DESC, id DESC;

-- name: ListSpecArtifactsForChange :many
SELECT * FROM spec_artifact
WHERE change_id = $1
ORDER BY version DESC, relative_path ASC;

-- name: GetSpecChangeInWorkspace :one
SELECT * FROM spec_change
WHERE id = $1 AND workspace_id = $2;

-- name: GetProjectInWorkspaceForSpecChanges :one
SELECT id FROM project
WHERE id = $1 AND workspace_id = $2;
