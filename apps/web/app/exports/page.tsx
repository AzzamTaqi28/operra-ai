import Link from "next/link"

import { AppShell } from "@/components/app-shell"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function ExportsPage() {
  return (
    <AppShell title="Exports" description="CSV exports for purchase requests, approval history, and audit logs.">
      <div className="stack">
        <Card>
          <CardHeader>
            <CardTitle>Purchase requests CSV</CardTitle>
            <CardDescription>Snapshot of request metadata and status history.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild>
              <Link href="/api/v1/exports/purchase-requests.csv">Download purchase requests CSV</Link>
            </Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Approval history CSV</CardTitle>
            <CardDescription>Full approval chain for audit and analysis.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="outline">
              <Link href="/api/v1/exports/approval-history.csv">Download approval history CSV</Link>
            </Button>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Audit logs CSV</CardTitle>
            <CardDescription>Organization-scoped event trail export.</CardDescription>
          </CardHeader>
          <CardContent>
            <Button asChild variant="outline">
              <Link href="/api/v1/exports/audit-logs.csv">Download audit logs CSV</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    </AppShell>
  )
}
