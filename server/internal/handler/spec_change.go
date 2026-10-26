package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/multica-ai/multica/server/internal/featureflags"
	db "github.com/multica-ai/multica/server/pkg/db/generated"
)

// SpecChangeResponse is deliberately metadata-only. Artifact contents live in
// the user's repository and will be read through the local daemon in a later
// slice; the API must not expose arbitrary local filesystem paths.
type SpecChangeResponse struct {
	ID             string  `json:"id"`
	WorkspaceID    string  `json:"workspace_id"`
	ProjectID      string  `json:"project_id"`
	ResourceID     *string `json:"resource_id,omitempty"`
	ChangeKey      string  `json:"change_key"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	CurrentVersion int32   `json:"current_version"`
	CreatedBy      *string `json:"created_by,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func specChangeToResponse(change db.SpecChange) SpecChangeResponse {
	return SpecChangeResponse{
		ID:             uuidToString(change.ID),
		WorkspaceID:    uuidToString(change.WorkspaceID),
		ProjectID:      uuidToString(change.ProjectID),
		ResourceID:     uuidToPtr(change.ResourceID),
		ChangeKey:      change.ChangeKey,
		Title:          change.Title,
		Status:         change.Status,
		CurrentVersion: change.CurrentVersion,
		CreatedBy:      uuidToPtr(change.CreatedBy),
		CreatedAt:      timestampToString(change.CreatedAt),
		UpdatedAt:      timestampToString(change.UpdatedAt),
	}
}

func (h *Handler) requireOpenSpecWorkflow(w http.ResponseWriter, r *http.Request) bool {
	if featureflags.OpenSpecWorkflowEnabled(r.Context(), h.FeatureFlags) {
		return true
	}
	writeFeatureDisabled(w, "openspec_workflow_disabled", "OpenSpec Workflow is not enabled")
	return false
}

// ListProjectSpecChanges is the first read-only API slice for OpenSpec
// discovery. The daemon will populate these rows once local discovery lands;
// keeping the project authorization check here makes the future write paths
// safe by construction.
func (h *Handler) ListProjectSpecChanges(w http.ResponseWriter, r *http.Request) {
	if !h.requireOpenSpecWorkflow(w, r) {
		return
	}
	workspaceID := h.resolveWorkspaceID(r)
	workspaceUUID, ok := parseUUIDOrBadRequest(w, workspaceID, "workspace_id")
	if !ok {
		return
	}
	projectUUID, ok := parseUUIDOrBadRequest(w, chi.URLParam(r, "id"), "project id")
	if !ok {
		return
	}
	if _, err := h.Queries.GetProjectInWorkspaceForSpecChanges(r.Context(), db.GetProjectInWorkspaceForSpecChangesParams{
		ID: projectUUID, WorkspaceID: workspaceUUID,
	}); err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	changes, err := h.Queries.ListSpecChangesForProject(r.Context(), db.ListSpecChangesForProjectParams{
		WorkspaceID: workspaceUUID,
		ProjectID:   projectUUID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list OpenSpec changes")
		return
	}
	response := make([]SpecChangeResponse, 0, len(changes))
	for _, change := range changes {
		response = append(response, specChangeToResponse(change))
	}
	writeJSON(w, http.StatusOK, map[string]any{"changes": response})
}
