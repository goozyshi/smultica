"use client";

import { FileText, GitBranch, ChevronRight } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { cn } from "@multica/ui/lib/utils";
import { projectSpecChangesOptions } from "@multica/core/projects/queries";
import type { SpecChange } from "@multica/core/types";

const statusStyles: Record<string, string> = {
  discovered: "bg-slate-400",
  draft: "bg-amber-400",
  active: "bg-blue-500",
  blocked: "bg-red-500",
  completed: "bg-emerald-500",
  archived: "bg-muted-foreground",
};

function ChangeRow({ change }: { change: SpecChange }) {
  return (
    <button
      type="button"
      className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-left transition-colors hover:bg-accent/60"
      title={`${change.change_key}: ${change.title}`}
    >
      <GitBranch className="size-3.5 shrink-0 text-muted-foreground" />
      <span className="min-w-0 flex-1 truncate text-caption">
        <span className="font-medium">{change.change_key}</span>
        {change.title && <span className="ml-1.5 text-muted-foreground">{change.title}</span>}
      </span>
      <span className={cn("size-2 shrink-0 rounded-full", statusStyles[change.status] ?? "bg-muted-foreground")} />
      <ChevronRight className="size-3 shrink-0 text-muted-foreground" />
    </button>
  );
}

export function OpenSpecWorkflowSection({ workspaceId, projectId }: { workspaceId: string; projectId: string }) {
  const { data: changes = [], isLoading } = useQuery(projectSpecChangesOptions(workspaceId, projectId));

  return (
    <div className="border-t pt-4">
      <div className="mb-2 flex items-center gap-1.5 px-2 text-caption font-medium">
        <FileText className="size-3.5 text-muted-foreground" />
        <span>OpenSpec Workflow</span>
        <span className="ml-auto text-muted-foreground">{changes.length}</span>
      </div>
      {isLoading ? (
        <div className="px-2 py-2 text-caption text-muted-foreground">Scanning changes…</div>
      ) : changes.length === 0 ? (
        <div className="px-2 py-2 text-caption text-muted-foreground">
          No OpenSpec changes discovered in this repository.
        </div>
      ) : (
        <div className="space-y-0.5">
          {changes.map((change) => <ChangeRow key={change.id} change={change} />)}
        </div>
      )}
    </div>
  );
}
