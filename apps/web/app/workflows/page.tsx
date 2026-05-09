import { redirect } from "next/navigation"

import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Textarea } from "@/components/ui/textarea"
import { apiGet, getToken, type ApiListResponse, type WorkflowListItem } from "@/lib/api"

function versionLabel(item: WorkflowListItem) {
  return item.active_version?.version_number ? `v${item.active_version.version_number}` : "—"
}

export default async function WorkflowsPage() {
  const token = await getToken()
  if (!token) redirect("/login")

  const workflows = await apiGet<ApiListResponse<WorkflowListItem>>("/api/v1/workflows?page_size=100", token)

  return (
    <AppShell title="Workflows" description="JSON-first definitions with validation, Mermaid preview, versioning, and activation.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Badge>Type</Badge>
          <Badge>Status</Badge>
        </div>
        <div className="toolbar-actions">
          <Button asChild variant="outline">
            <Link href="/ai-workflow">AI Builder</Link>
          </Button>
          <Button type="button">Create workflow</Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Workflow list</CardTitle>
          <CardDescription>Versioned definitions and activation state.</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Active Version</TableHead>
                <TableHead>Updated At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {workflows.data.map((workflow) => (
                <TableRow key={workflow.id}>
                  <TableCell>{workflow.name}</TableCell>
                  <TableCell>{workflow.type}</TableCell>
                  <TableCell><Badge>{workflow.status}</Badge></TableCell>
                  <TableCell>{versionLabel(workflow)}</TableCell>
                  <TableCell>{new Date(workflow.updated_at).toLocaleDateString("id-ID")}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Workflow JSON editor</CardTitle>
          <CardDescription>Validate first, then save a new version or activate it.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="editor-grid">
            <Textarea className="code-area" defaultValue={`{\n  "name": "Purchase Request Approval",\n  "type": "purchase_request",\n  "version": 3,\n  "steps": []\n}`} />
            <div className="preview-pane">
              <p className="eyebrow">Validation</p>
              <p className="muted-copy">No steps defined yet. The backend validator will reject this version until the required approval steps are provided.</p>
              <div className="action-row">
                <Button type="button" variant="outline">Validate</Button>
                <Button type="button">Save new version</Button>
                <Button type="button" variant="outline">Activate version</Button>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </AppShell>
  )
}
