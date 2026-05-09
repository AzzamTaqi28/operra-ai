import { redirect } from "next/navigation"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { apiGet, getToken, type ApiListResponse, type WorkflowListItem, type WorkflowVersionSummary } from "@/lib/api"
import { WorkflowEditorForm } from "@/components/workflow-editor-form"

export default async function WorkflowDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const token = await getToken()
  if (!token) redirect("/login")

  const { id } = await params

  const [workflowRes, versionsRes] = await Promise.all([
    apiGet<{ data: WorkflowListItem }>(`/api/v1/workflows/${id}`, token),
    apiGet<ApiListResponse<WorkflowVersionSummary>>(`/api/v1/workflows/${id}/versions?page_size=100`, token).catch(() => ({ data: [], pagination: { page: 1, page_size: 100, total: 0 } })),
  ])

  const workflow = workflowRes.data

  return (
    <AppShell title={`Workflow ${workflow.name}`} description="Edit workflow versions and activation state.">
      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Workflow summary</CardTitle>
            <CardDescription>Current workflow metadata and active version.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-2">
            <Badge>{workflow.type}</Badge>
            <Badge>{workflow.status}</Badge>
            <p className="text-sm text-slate-600">Active version: {workflow.active_version?.version_number ? `v${workflow.active_version.version_number}` : "none"}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Version editor</CardTitle>
            <CardDescription>Save a new JSON version or activate an existing one.</CardDescription>
          </CardHeader>
          <CardContent>
            <WorkflowEditorForm workflow={workflow} versions={versionsRes.data} />
          </CardContent>
        </Card>
      </div>
    </AppShell>
  )
}
