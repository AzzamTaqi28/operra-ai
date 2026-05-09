import { AppShell } from "@/components/app-shell"
import { Card, MetricCard, Table } from "@/components/ui"

const metrics = [
  { label: "Total purchase requests", value: "128", hint: "All organizations scoped to the current tenant" },
  { label: "Pending approvals", value: "8", hint: "Awaiting action on the active workflow step" },
  { label: "Approved this month", value: "74", hint: "Includes completed and processing requests" },
  { label: "Rejected this month", value: "6", hint: "Logged with a full audit trail" },
  { label: "Average approval time", value: "14h 32m", hint: "Median for the last 30 days" },
]

const recentRequests = [
  ["PR-1024", "Laptop refresh", "Siti", "Finance", "IDR 45,000,000", "in_review", "Director Approval", "2026-05-08"],
  ["PR-1023", "Office chairs", "Ari", "People Ops", "IDR 12,500,000", "draft", "Draft", "2026-05-08"],
  ["PR-1022", "Cloud storage", "Dina", "Engineering", "IDR 21,000,000", "revision_requested", "Manager Approval", "2026-05-07"],
]

const pendingApprovals = [
  ["PR-1024", "Laptop refresh", "Siti", "Finance", "IDR 45,000,000", "Director Approval", "1h ago"],
  ["PR-1019", "Security training", "Budi", "Operations", "IDR 6,500,000", "Manager Approval", "3h ago"],
]

export default function DashboardPage() {
  return (
    <AppShell title="Dashboard" description="Operational overview for purchase approvals, workflows, and audit trails.">
      <section className="card-grid">
        {metrics.map((item) => (
          <MetricCard key={item.label} label={item.label} value={item.value} hint={item.hint} />
        ))}
      </section>

      <div className="stack">
        <Card title="Recent requests" description="Status should be obvious at a glance.">
          <Table
            headers={["Request ID", "Title", "Requester", "Department", "Amount", "Status", "Current Step", "Created At"]}
            rows={recentRequests.map((row) => row.map((cell) => cell))}
          />
        </Card>

        <Card title="My pending approvals" description="Fast access to the current actionable step.">
          <Table
            headers={["Request ID", "Title", "Requester", "Department", "Amount", "Current Step", "Waiting Since"]}
            rows={pendingApprovals.map((row) => row.map((cell) => cell))}
          />
        </Card>
      </div>
    </AppShell>
  )
}
