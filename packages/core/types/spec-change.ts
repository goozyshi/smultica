export type SpecChangeStatus =
  | "discovered"
  | "draft"
  | "active"
  | "blocked"
  | "completed"
  | "archived";

export interface SpecChange {
  id: string;
  workspace_id: string;
  project_id: string;
  resource_id?: string;
  change_key: string;
  title: string;
  status: SpecChangeStatus;
  current_version: number;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface ListProjectSpecChangesResponse {
  changes: SpecChange[];
}
