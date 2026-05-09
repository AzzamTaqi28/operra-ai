import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Card, Chip, Table } from "@/components/ui"

const workflows = [
  ["Purchase Request Approval", "purchase_request", "active", "v3", "2026-05-08"],
  ["Travel Request Approval", "purchase_request", "draft", "v1", "2026-05-06"],
]

export default function WorkflowsPage() {
  return (
    <AppShell title="Workflows" description="JSON-first definitions with validation, Mermaid preview, versioning, and activation.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Chip>Type</Chip>
          <Chip>Status</Chip>
        </div>
        <div className="toolbar-actions">
          <Link href="/ai-workflow" className="button button-outline">AI Builder</Link>
          <button className="button button-solid" type="button">Create workflow</button>
        </div>
      </div>

      <Card title="Workflow list">
        <Table headers={["Name", "Type", "Status", "Active Version", "Updated At"]} rows={workflows.map((row) => row.map((cell) => cell))} />
      </Card>

      <Card title="Workflow JSON editor" description="Validate first, then save a new version or activate it.">
        <div className="editor-grid">
          <textarea className="textarea code-area" defaultValue={`{\n  "name": "Purchase Request Approval",\n  "type": "purchase_request",\n  "version": 3,\n  "steps": []\n}`} />
          <div className="preview-pane">
            <p className="eyebrow">Validation</p>
            <p className="muted-copy">No steps defined yet. The backend validator will reject this version until the required approval steps are provided.</p>
            <div className="action-row">
              <button className="button button-outline" type="button">Validate</button>
              <button className="button button-solid" type="button">Save new version</button>
              <button className="button button-outline" type="button">Activate version</button>
            </div>
          </div>
        </div>
      </Card>
    </AppShell>
  )
}
