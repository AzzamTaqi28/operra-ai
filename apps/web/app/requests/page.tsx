import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Card, Chip, Table } from "@/components/ui"

const rows = [
  ["PR-1024", "Laptop refresh", "Siti", "Finance", "IDR 45,000,000", "in_review", "Director Approval", "2026-05-08"],
  ["PR-1023", "Office chairs", "Ari", "People Ops", "IDR 12,500,000", "draft", "Draft", "2026-05-08"],
  ["PR-1022", "Cloud storage", "Dina", "Engineering", "IDR 21,000,000", "revision_requested", "Manager Approval", "2026-05-07"],
]

export default function RequestsPage() {
  return (
    <AppShell title="Purchase Requests" description="Drafts, submissions, approval timelines, comments, and attachments.">
      <div className="toolbar">
        <div className="toolbar-filters">
          <Chip>Status</Chip>
          <Chip>Department</Chip>
          <Chip>Date range</Chip>
          <Chip>Search</Chip>
        </div>
        <div className="toolbar-actions">
          <Link href="/requests/new" className="button button-solid">
            Create request
          </Link>
          <a href="/api/v1/exports/purchase-requests.csv" className="button button-outline">
            Export CSV
          </a>
        </div>
      </div>

      <Card title="Request list" description="Request ID, title, requester, department, amount, status, and current step.">
        <Table
          headers={["Request ID", "Title", "Requester", "Department", "Amount", "Status", "Current Step", "Created At"]}
          rows={rows.map((row) => [
            <Link key={row[0]} href={`/requests/${row[0].toLowerCase()}`}>{row[0]}</Link>,
            row[1],
            row[2],
            row[3],
            row[4],
            row[5],
            row[6],
            row[7],
          ])}
        />
      </Card>
    </AppShell>
  )
}
