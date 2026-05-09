import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Card } from "@/components/ui"

export default function ExportsPage() {
  return (
    <AppShell title="Exports" description="CSV exports for purchase requests, approval history, and audit logs.">
      <div className="stack">
        <Card title="Purchase requests CSV">
          <Link href="/api/v1/exports/purchase-requests.csv" className="button button-solid">Download purchase requests CSV</Link>
        </Card>
        <Card title="Approval history CSV">
          <Link href="/api/v1/exports/approval-history.csv" className="button button-outline">Download approval history CSV</Link>
        </Card>
        <Card title="Audit logs CSV">
          <Link href="/api/v1/exports/audit-logs.csv" className="button button-outline">Download audit logs CSV</Link>
        </Card>
      </div>
    </AppShell>
  )
}
