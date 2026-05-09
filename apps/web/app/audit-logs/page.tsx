import { AppShell } from "@/components/app-shell"
import { Card, Chip, Table } from "@/components/ui"

const logs = [
  ["2026-05-08 10:11", "Siti", "request.submitted", "purchase_request", "PR-1024"],
  ["2026-05-08 10:25", "Mira", "approval.approved", "approval_step_instance", "asi-022"],
  ["2026-05-08 10:31", "Siti", "csv.exported", "export", "-"],
]

export default function AuditLogsPage() {
  return (
    <AppShell title="Audit Logs" description="Read-only event trail for meaningful mutations, approvals, exports, and AI generation.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Chip>Action</Chip>
          <Chip>Entity type</Chip>
          <Chip>Actor</Chip>
          <Chip>Date range</Chip>
        </div>
      </div>

      <Card title="Event trail">
        <Table headers={["Timestamp", "Actor", "Action", "Entity Type", "Entity ID"]} rows={logs.map((row) => row.map((cell) => cell))} />
      </Card>
    </AppShell>
  )
}
