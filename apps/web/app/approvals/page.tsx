import { AppShell } from "@/components/app-shell"
import { Card, Table } from "@/components/ui"

const approvals = [
  ["PR-1024", "Laptop refresh", "Siti", "Finance", "IDR 45,000,000", "Director Approval", "1h ago"],
  ["PR-1019", "Security training", "Budi", "Operations", "IDR 6,500,000", "Manager Approval", "3h ago"],
]

export default function ApprovalsPage() {
  return (
    <AppShell title="Pending Approvals" description="Current user's actionable queue, scoped by role and organization.">
      <Card title="Approval queue">
        <Table headers={["Request ID", "Title", "Requester", "Department", "Amount", "Current Step", "Waiting Since"]} rows={approvals.map((row) => row.map((cell) => cell))} />
      </Card>
    </AppShell>
  )
}
